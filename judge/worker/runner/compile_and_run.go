package runner

import (
	"context"
	"errors"
	"fmt"
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
)

type testCaseResult struct{
	result Verdict
	mem_usage int64
	time_ms int
}

func validateParams(runner_parms types.RunnerParamsJson) error {
	if exist,err := utils.DirExists(runner_parms.TestCasesDir); err!=nil || exist == false {
		fmt.Println(err,exist)
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

func CompileIfCompilable(cur_runtime string,codeDir string) (error) {
	runtime_conf,ok := runtime.GetRuntime(cur_runtime)
	if ok == false {
		return errors.New("[exec error] runtime does not exist")
	}
	
	if len(runtime_conf.CompileComand) == 0 { //no compilation needed
		return nil
	}

	// max_compilation_timeout,_ := strconv.Atoi(os.Getenv("WORKER_MAX_COMPILATION_TIME"))
	max_compilation_timeout := 10000
	ctx, cancel := context.WithTimeout(context.Background(),time.Millisecond * time.Duration(max_compilation_timeout))
	defer cancel()

	comipler_path := runtime_conf.CompileComand[0]
	args := runtime_conf.CompileComand[1:]
	cmd := exec.CommandContext(ctx,comipler_path,args...)

	var err error
	cmd.Stdin,err = os.Open(filepath.Join(codeDir,runtime_conf.CodeFileName))
	if err != nil {
		return errors.New("[exec error] failed to open code file")
	}
	
	cmd.Stdout,err = os.OpenFile(runtime.StdOutFileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return errors.New("[exec error] failed to open stdout file")
	}

	cmd.Stderr,err = os.OpenFile(runtime.StdErrorFileName,os.O_RDWR|os.O_CREATE|os.O_APPEND,0644)
	if err != nil {
		return errors.New("[exec error] failed to open stderr file")
	}

	//start compilation
	if err := cmd.Run(); err != nil {
		return errors.New("[exec error] compilation failed with error : "+err.Error())
	}
	
	return nil
	
}

func CompileRunAndTests(runner_parms types.RunnerParamsJson) (Verdict,error) {

	if err:= validateParams(runner_parms); err != nil {
		return VerdictInternalError,err
	}

	err := CompileIfCompilable(runner_parms.Runtime,runner_parms.CodeDir)
	if err != nil {
		return VerdictCompilationError,err
	}

	//each test case will be grouped into a directory with 
	//the test case number as the name of the dir
	FinalVerdict := VerdictAccepted
	for test_case_no := 1;;test_case_no++ {
		str_test_no := strconv.Itoa(test_case_no)
		test_case_path := filepath.Join(runner_parms.TestCasesDir,str_test_no)
		if exist,err := utils.DirExists(test_case_path); err != nil || exist == false {
			break
		}
		test_inp_file := filepath.Join(test_case_path,"inp.txt")
		if exist,err := utils.FileExists(test_inp_file); err!=nil || exist == false {
			fmt.Println("[exec error] file inp not found with error")
			continue
		}
		test_exp_out_file := filepath.Join(test_case_path,"out.txt")
		if exist,err := utils.FileExists(test_exp_out_file); err!=nil || exist == false {
			fmt.Println("[exec error] file out not found with error")
			continue
		}

		verdict := runForSingleTestCase(runner_parms,test_inp_file,test_exp_out_file)
		fmt.Println("test case",test_case_no,"done bro",verdict)
		if verdict != VerdictAccepted {
			FinalVerdict = verdict
		}

	}

	return FinalVerdict,nil

}


func runForSingleTestCase(runner_parms types.RunnerParamsJson,test_inp_file string,test_out_file string) Verdict {
	runtime_conf,ok := runtime.GetRuntime(runner_parms.Runtime)
	if ok == false {
		return VerdictInternalError
	}
	exe_or_interpreter := runtime_conf.CompileComand[0]
	args := runtime_conf.CompileComand[1:]

	in,err1  := os.Open(test_inp_file)
	out,err2 := os.OpenFile(runtime.StdOutFileName,os.O_CREATE | os.O_TRUNC | os.O_RDWR,0644)
	errfile,err3 := os.OpenFile(runtime.StdErrorFileName,os.O_CREATE | os.O_APPEND | os.O_RDWR,0644)
	if err1!=nil || err2!=nil || err3!=nil {
		return VerdictInternalError
	}
	defer in.Close()
	defer out.Close()
	defer errfile.Close()

	ctx, cancel := context.WithTimeout(context.Background(),time.Millisecond * time.Duration(runner_parms.TimeConstrain))
	defer cancel()

	cmd := exec.CommandContext(ctx,exe_or_interpreter,args...)
	cmd.Stdin = in
	cmd.Stdout = out
	cmd.Stderr = errfile

	var ru_befor syscall.Rusage
	syscall.Getrusage(syscall.RUSAGE_CHILDREN, &ru_befor)

	if err := cmd.Run(); err != nil {
		fmt.Println("exec error ", err)
		return VerdictWrongAns
	}

	var ru_after syscall.Rusage
	syscall.Getrusage(syscall.RUSAGE_CHILDREN, &ru_after)

	data, err := os.ReadFile(runtime.StdOutFileName)
	if err != nil {
		return VerdictInternalError
	}
	fmt.Println("[runner] test case output : ",string(data))

	mem_usage := ru_after.Maxrss - ru_befor.Maxrss
	if mem_usage > int64(runner_parms.MemConstrain) {
		return VerdictMLE
	}
	if ctx.Err() == context.Canceled { 
		fmt.Println("process terminated due to TLE")
		return VerdictTLE
	}

	return VerdictAccepted
}

