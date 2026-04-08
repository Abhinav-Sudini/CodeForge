package runner

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
	"time"
	"worker/config"
	MyLog "worker/logger"
	"worker/runtime"
	"worker/types"
	"worker/utils"
)

type Verdict int

const (
	VerdictAccepted Verdict = iota
	VerdictWrongAns
	VerdictCompilationError
	VerdictTLE
	VerdictMLE
	VerdictCodeToBig
	VerdictInternalError
	VerdictBadRequest
)

type SubmitionResult struct {
	Result       int    `json:"Result"`
	Mem_usage    int    `json:"Mem_usage"`
	Time_ms      int    `json:"Time_ms"`
	WA_Test_case int    `json:"WA_Test_case"`
	WA_Stdout    string `json:"WA_Stdout"`
	Stderr       string `json:"Stderr"`
	MSG          string `json:"MSG"`
}

type ResourcesUsed struct {
	Time_ms int
	Mem_kb  int
}

func validateParams(runner_parms types.RunnerParamsJson) error {
	if exist, err := utils.DirExistsAndValidPerms(runner_parms.TestCasesDir, "r"); err != nil || exist == false {
		fmt.Println(err, exist)
		return errors.New("[exec error] Test dir does not exist : ")
	}
	if runner_parms.Runtime != os.Getenv("WORKER_RUNTIME") {
		return errors.New("[exec error] runtime not available")
	}
	if runner_parms.TimeConstrain <= 0 || runner_parms.MemConstrain <= 0 {
		return errors.New("[exec error] time/mem constrains invalid")
	}
	return nil
}

func CompileIfCompilable(cur_runtime string, codeDir string) (bool, error) {
	runtime_conf, ok := runtime.GetRuntime(cur_runtime)
	if ok == false {
		return false, errors.New("[exec error] runtime does not exist")
	}

	if len(runtime_conf.CompileComand) == 0 { //no compilation needed
		return true, nil
	}

	max_compilation_timeout, _ := strconv.Atoi(os.Getenv("WORKER_MAX_COMPILATION_TIME"))
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(max_compilation_timeout*10))
	defer cancel()

	compiler_path := runtime_conf.CompileComand[0]
	args := runtime_conf.CompileComand[1:]
	cmd := exec.CommandContext(ctx, compiler_path, args...)

	var err error
	cmd.Stdin, err = os.Open(filepath.Join(codeDir, runtime_conf.CodeFileName))
	if err != nil {
		return false, errors.New("[exec error] failed to open code file")
	}

	var buf_err, buf_out bytes.Buffer
	cmd.Stdout = &buf_out
	cmd.Stderr = &buf_err

	//start compilation
	MyLog.Printdev("[compl comand]", compiler_path, args)
	if err := cmd.Run(); err != nil {
		saveStderrAndStdoutToFile(codeDir, &buf_out, &buf_err)
		switch exitErr := err.(type) {
		case *exec.ExitError:
			MyLog.Printdev("exec run", "exit code : ", exitErr.ExitCode(), "stderr : ", exitErr.Stderr, "err from cmd: ", err.Error())
			return false, nil
		default:
			MyLog.Print("[exec worker compile]", "failed with err:", err)
			return false, errors.New("[exec error] compilation failed with error : " + err.Error())
		}
	}

	return true, nil

}

func CompileRunAndTests(runner_parms types.RunnerParamsJson) (SubmitionResult, error) {

	var FinalResult SubmitionResult
	FinalResult.Result = int(VerdictAccepted)

	if err := validateParams(runner_parms); err != nil {
		FinalResult.Result = int(VerdictBadRequest)
		return FinalResult, err
	}

	compile_done, err := CompileIfCompilable(runner_parms.Runtime, runner_parms.CodeDir)
	if err != nil {
		return FinalResult, err
	}
	if compile_done == false {
		FinalResult.Result = int(VerdictCompilationError)
		return FinalResult, nil
	}
	MyLog.Printdev("exec worker", "compilation done")

	//each test case will be grouped into a directory with
	//the test case number as the name of the dir
	for test_case_no := 1; ; test_case_no++ {
		str_test_no := strconv.Itoa(test_case_no)
		test_case_dir := runner_parms.TestCasesDir
		inp_file := str_test_no + ".in"
		out_file := str_test_no + ".out"
		test_inp_file := filepath.Join(test_case_dir, inp_file)

		if exist, err := utils.FileExistsAndValidPerms(test_inp_file, "r"); err != nil || exist == false {
			if err != nil {
				fmt.Println("[exec error] can not open file inp with error ", err)
			}
			break
		}
		test_exp_out_file := filepath.Join(test_case_dir, out_file)
		if exist, err := utils.FileExistsAndValidPerms(test_exp_out_file, "r"); err != nil || exist == false {
			if err != nil {
				fmt.Println("[exec error] can not open file out with error ", err)
			}
			break
		}

		// run the actual test case
		verdict, resourses_used := runForSingleTestCase(runner_parms, test_inp_file, test_exp_out_file)
		FinalResult.Time_ms = max(FinalResult.Time_ms, resourses_used.Time_ms)
		FinalResult.Mem_usage = max(FinalResult.Mem_usage, resourses_used.Mem_kb)
		MyLog.Printdev("[Executioner]", "test case", test_case_no, "done bro", verdict, "time : ", resourses_used.Time_ms)

		if verdict != VerdictAccepted {
			FinalResult.Result = int(verdict)
			FinalResult.WA_Test_case = test_case_no
			break
		}

	}

	return FinalResult, nil

}

func runForSingleTestCase(runner_parms types.RunnerParamsJson, test_inp_file string, test_out_file string) (Verdict, ResourcesUsed) {

	//seting up return values and runtime
	var resourses_used ResourcesUsed
	runtime_conf, ok := runtime.GetRuntime(runner_parms.Runtime)
	if ok == false {
		return VerdictInternalError, resourses_used
	}

	exe_or_interpreter := runtime_conf.RunComand[0]
	args := runtime_conf.RunComand[1:]

	code_inp, err1 := os.ReadFile(test_inp_file)
	if err1 != nil {
		MyLog.Print("test case runner", "can not read file with err :", err1)
		return VerdictInternalError, resourses_used
	}

	buf_err := utils.BoundBuffer{N: config.MaxOutputBufferSize}
	buf_out := utils.BoundBuffer{N: config.MaxOutputBufferSize}

	//creating the comand to be run
	time_const_multiplier := 3 // as we can hit the dead line but the program might not run for whole time
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(runner_parms.TimeConstrain*time_const_multiplier))
	defer cancel()

	cmd := exec.CommandContext(ctx, exe_or_interpreter, args...)
	cmd.Stdin = bytes.NewBuffer(code_inp)
	cmd.Stdout = &buf_out
	cmd.Stderr = &buf_err

	//remove the parents ENV variables
	cmd.Env = utils.CopyEnvVariablesOfParent(config.WorkerRuntimeENVVariablesToInclude)

	//running the comand
	MyLog.Printdev("singele exec runner", "cmd comand being run", cmd.String())
	if err := cmd.Start(); err != nil {
		MyLog.Print("[exec runtime error] ", "failed starting cmd", err, ctx.Err())
		saveStderrAndStdoutToFile(runner_parms.CodeDir, &buf_out.Buf, &buf_err.Buf)
		return VerdictWrongAns, resourses_used
	}

	err := cmd.Wait()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return VerdictTLE, ResourcesUsed{
				Mem_kb:  0,
				Time_ms: runner_parms.TimeConstrain,
			}
		}
		MyLog.Printdev("[exec runtime error] ", "failed at wait call :", err)
		saveStderrAndStdoutToFile(runner_parms.CodeDir, &buf_out.Buf, &buf_err.Buf)
		return VerdictWrongAns, resourses_used
	}

	rusage := cmd.ProcessState.SysUsage().(*syscall.Rusage)
	mem_usage := rusage.Maxrss
	time_usage_in_kernal := utils.GetTimeInMillSec(rusage.Stime)
	time_usage_in_user := utils.GetTimeInMillSec(rusage.Utime)
	resourses_used.Mem_kb = int(mem_usage)
	resourses_used.Time_ms = int(time_usage_in_kernal + time_usage_in_user)

	if resourses_used.Time_ms > runner_parms.TimeConstrain {
		return VerdictTLE, ResourcesUsed{
			Mem_kb:  0,
			Time_ms: resourses_used.Time_ms,
		}
	}
	MyLog.Printdev("runstats kernerl time ", time_usage_in_kernal, "user_time", time_usage_in_user)

	if mem_usage > int64(runner_parms.MemConstrain) {
		return VerdictMLE, resourses_used
	}

	exp_output, err := os.Open(test_out_file)
	if err != nil {
		return VerdictInternalError, ResourcesUsed{}
	}

	// fmt.Println("out : ", buf_out.String())
	// buf_out_copy := append([]byte(nil), buf_out.Buf.Bytes()...)
	// buf_out = *bytes.NewBuffer(buf_out_copy)
	truncated_stdout := buf_out.Buf.Bytes()[:min(len(buf_out.Buf.Bytes()), 2000)]
	truncated_stderr := buf_err.Buf.Bytes()[:min(len(buf_err.Buf.Bytes()), 2000)]
	output_same := OutputJudge(exp_output, &buf_out.Buf)
	if output_same == false {
		saveStderrAndStdoutToFile(runner_parms.CodeDir, bytes.NewBuffer(truncated_stdout), bytes.NewBuffer(truncated_stderr))
		return VerdictWrongAns, resourses_used
	}

	return VerdictAccepted, resourses_used
}

func saveStderrAndStdoutToFile(codeDir string, stdout io.Reader, stderr io.Reader) {
	err_file_path := filepath.Join(codeDir, runtime.StdErrorFileName)
	out_file_path := filepath.Join(codeDir, runtime.StdOutFileName)
	utils.SaveFileFromBuf(err_file_path, stderr)
	utils.SaveFileFromBuf(out_file_path, stdout)
}

func (ver Verdict) String() string {
	switch ver {
	case VerdictAccepted:
		return "Accepted"
	case VerdictWrongAns:
		return "Wrong Answer"
	case VerdictMLE:
		return "Memory Limit Exceeded"
	case VerdictTLE:
		return "Time Limit Exceeded"
	case VerdictCompilationError:
		return "Compilation Error"
	case VerdictCodeToBig:
		return "Code Size To Big"
	case VerdictInternalError:
		return "Internal Server Error"
	}
	return "Internal error bro"
}
