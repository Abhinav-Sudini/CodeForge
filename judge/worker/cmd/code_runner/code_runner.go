package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os/user"
	"path/filepath"
	"strconv"
	"syscall"
	"worker/runner"
	"worker/runtime"
	"worker/utils"
)

const nobodyUID = 65534
const nobodyGID = 65534

func checkIfRoot() (bool,error) {
		currentUser, err := user.Current()
    if err != nil {
			return false,errors.New("[exec error] failed to read UID and Gid")
    }
    // Print UID and GID
		uid,err := strconv.Atoi(currentUser.Uid)
		if err != nil {
			return false,errors.New("[exec error] failed to convert uit to int")
		}
		if uid==0 {
			return true,nil
		}else{
			return false,nil
		}
}

// the program must run as root as setGid and SetUid require privleges
func main(){
	isRoot,err := checkIfRoot()
	if err!=nil {
		panic(err)
	}
	if isRoot == false {
		fmt.Println("[exec error] you need to be root in order to run the application")
	}

	// drop privleges immediatly from root
	// err = syscall.Unshare(syscall.CLONE_NEWNET)
	// if err != nil {
	// 		panic(err)
	// }

	//Drop all groups
	if err := syscall.Setgroups([]int{}); err != nil {
			panic(err)
	}

	// Drop GID first
	if err := syscall.Setgid(nobodyGID); err != nil {
			panic(err)
	}

	// Drop UID
	if err := syscall.Setuid(nobodyUID); err != nil {
			panic(err)
	}
	fmt.Println("[worker] all privleges droped ")
	
	runner_params,err := runner.GetRunnerParams() //params are passed in as json through stdin
	if err != nil {
		fmt.Println("[exec error] plz pass in the correct params")
	}

	fmt.Println("[New Job] running the compilation proc")
	fmt.Println("[New Job] running the compilation proc")
	fmt.Println("[New Job] running the compilation proc")
	fmt.Println(runner_params)

	//call to compile and test
	verdict,err := runner.CompileRunAndTests(runner_params)
	if err != nil {
		panic(err)
	}

	//generating the output file
	verdict.MSG = runner.GenerateResultMSG(verdict)
	json_res,err := json.Marshal(verdict)
	if err != nil {
		fmt.Println("[worker] unable to marshal")
		panic(err)
	}
	io_reader := bytes.NewBuffer(json_res)
	verdict_file_path := filepath.Join(runner_params.CodeDir,runtime.VerdictFileName)
	utils.SaveFileFromBuf(verdict_file_path,io_reader)

	fmt.Println("[Verdict] : ",string(json_res))

}
