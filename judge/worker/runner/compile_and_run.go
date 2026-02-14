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

func CompileIfCompilable(cur_runtime string, codeDir string) error {
	runtime_conf, ok := runtime.GetRuntime(cur_runtime)
	if ok == false {
		return errors.New("[exec error] runtime does not exist")
	}

	if len(runtime_conf.CompileComand) == 0 { //no compilation needed
		return nil
	}

	max_compilation_timeout, _ := strconv.Atoi(os.Getenv("WORKER_MAX_COMPILATION_TIME"))
	// max_compilation_timeout := 10000
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(max_compilation_timeout))
	defer cancel()

	comipler_path := runtime_conf.CompileComand[0]
	args := runtime_conf.CompileComand[1:]
	cmd := exec.CommandContext(ctx, comipler_path, args...)

	var err error
	cmd.Stdin, err = os.Open(filepath.Join(codeDir, runtime_conf.CodeFileName))
	if err != nil {
		return errors.New("[exec error] failed to open code file")
	}

	var buf_err, buf_out bytes.Buffer
	cmd.Stdout = &buf_out
	cmd.Stderr = &buf_err

	//start compilation
	if err := cmd.Run(); err != nil {
		saveStderrAndStdoutToFile(codeDir,&buf_out,&buf_err)
		return errors.New("[exec error] compilation failed with error : " + err.Error())
	}

	return nil

}

func CompileRunAndTests(runner_parms types.RunnerParamsJson) (SubmitionResult, error) {

	var FinalResult SubmitionResult
	FinalResult.Result = int(VerdictAccepted)

	if err := validateParams(runner_parms); err != nil {
		FinalResult.Result = int(VerdictBadRequest)
		return FinalResult, err
	}

	err := CompileIfCompilable(runner_parms.Runtime, runner_parms.CodeDir)
	if err != nil {
		FinalResult.Result = int(VerdictCompilationError)
		return FinalResult, err
	}

	//each test case will be grouped into a directory with
	//the test case number as the name of the dir
	for test_case_no := 1; ; test_case_no++ {
		str_test_no := strconv.Itoa(test_case_no)
		test_case_path := filepath.Join(runner_parms.TestCasesDir, str_test_no)
		if exist, err := utils.DirExistsAndValidPerms(test_case_path, "r"); err != nil || exist == false {
			break
		}
		test_inp_file := filepath.Join(test_case_path, "inp.txt")
		if exist, err := utils.FileExistsAndValidPerms(test_inp_file, "r"); err != nil || exist == false {
			fmt.Println("[exec error] file inp not found with error ", err)
			continue
		}
		test_exp_out_file := filepath.Join(test_case_path, "out.txt")
		if exist, err := utils.FileExistsAndValidPerms(test_exp_out_file, "r"); err != nil || exist == false {
			fmt.Println("[exec error] file out not found with error", err)
			continue
		}

		// run the actual test case
		verdict, resourses_used := runForSingleTestCase(runner_parms, test_inp_file, test_exp_out_file)
		FinalResult.Time_ms = max(FinalResult.Time_ms, resourses_used.Time_ms)
		FinalResult.Mem_usage = max(FinalResult.Mem_usage, resourses_used.Mem_kb)
		fmt.Println("[Executioner] test case", test_case_no, "done bro", GenerateResultMSG(FinalResult))

		if verdict != VerdictAccepted {
			FinalResult.Result = int(verdict)
			FinalResult.WA_Test_case = test_case_no
			break
		}

	}

	return FinalResult, nil

}

func runForSingleTestCase(runner_parms types.RunnerParamsJson, test_inp_file string, test_out_file string) (Verdict, ResourcesUsed) {

	var resourses_used ResourcesUsed

	runtime_conf, ok := runtime.GetRuntime(runner_parms.Runtime)
	if ok == false {
		return VerdictInternalError, resourses_used
	}
	exe_or_interpreter := runtime_conf.CompileComand[0]
	args := runtime_conf.CompileComand[1:]

	in, err1 := os.Open(test_inp_file)
	if err1 != nil {
		return VerdictInternalError, resourses_used
	}
	defer in.Close()

	var out_buf bytes.Buffer
	var err_buf bytes.Buffer
	out := &out_buf
	errfile := &err_buf

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(runner_parms.TimeConstrain))
	defer cancel()

	cmd := exec.CommandContext(ctx, exe_or_interpreter, args...)
	cmd.Stdin = in
	cmd.Stdout = out
	cmd.Stderr = errfile

	var ru_befor syscall.Rusage
	syscall.Getrusage(syscall.RUSAGE_CHILDREN, &ru_befor)
	st_time := time.Now()

	if err := cmd.Run(); err != nil {
		fmt.Println("exec error ", err)
		saveStderrAndStdoutToFile(runner_parms.CodeDir,&out_buf,&err_buf)
		return VerdictWrongAns, resourses_used
	}

	var ru_after syscall.Rusage
	syscall.Getrusage(syscall.RUSAGE_CHILDREN, &ru_after)
	end_time := time.Now()


	mem_usage := ru_after.Maxrss - ru_befor.Maxrss
	time_usage := end_time.Sub(st_time).Milliseconds()
	resourses_used.Mem_kb = int(mem_usage)
	resourses_used.Time_ms = int(time_usage)

	if mem_usage > int64(runner_parms.MemConstrain) {
		return VerdictMLE, resourses_used
	}
	if ctx.Err() == context.Canceled {
		fmt.Println("process terminated due to TLE")
		return VerdictTLE, resourses_used
	}

	exp_output, err := os.Open(test_out_file)
	if err != nil {
		return VerdictInternalError, ResourcesUsed{}
	}
	output_same := OutputJudge(exp_output, out)
	if output_same == false {
		saveStderrAndStdoutToFile(runner_parms.CodeDir,&out_buf,&err_buf)
		return VerdictWrongAns, resourses_used
	}

	return VerdictAccepted, resourses_used
}

func saveStderrAndStdoutToFile(codeDir string,stdout io.Reader,stderr io.Reader){
		err_file_path := filepath.Join(codeDir,runtime.StdErrorFileName)
		out_file_path := filepath.Join(codeDir,runtime.StdOutFileName)
		utils.SaveFileFromBuf(err_file_path,stderr)
		utils.SaveFileFromBuf(out_file_path,stdout)
}

// func verdictStr(ver Verdict) string {
// 	switch ver {
// 	case VerdictAccepted:
// 		return "Accepted"
// 	case VerdictWrongAns:
// 		return "Wrong Answer"
// 	case VerdictMLE:
// 		return "Memory Limit Exceeded"
// 	case VerdictTLE:
// 		return "Time Limit Exceeded"
// 	case VerdictCompilationError:
// 		return "Compilation Error"
// 	case VerdictCodeToBig:
// 		return "Code Size To Big"
// 	case VerdictInternalError:
// 		return "Internal Server Error"
// 	}
// 	return "Internal error bro"
// }

