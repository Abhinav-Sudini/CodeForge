package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"syscall"
	MyLog "worker/logger"
	"worker/runner"
	"worker/runtime"
	"worker/utils"
)

const nobodyUID = 65534
const nobodyGID = 65534

func checkIfRoot() (bool, error) {
	currentUser, err := user.Current()
	if err != nil {
		return false, errors.New("[exec error] failed to read UID and Gid")
	}
	// Print UID and GID
	uid, err := strconv.Atoi(currentUser.Uid)
	if err != nil {
		return false, errors.New("[exec error] failed to convert uit to int")
	}
	if uid == 0 {
		return true, nil
	} else {
		return false, nil
	}
}

// the program must run as root as setGid and SetUid require privleges
func main() {
	MyLog.Printdev("exce worker", "exec worker start")
	isRoot, err := checkIfRoot()
	if err != nil {
		panic(err)
	}
	if isRoot == false {
		fmt.Println("[exec error] you need to be root in order to run the application")
	}

	// Somtimes gives an error if in a docker container as there are container restrictions
	// err = syscall.Unshare(syscall.CLONE_NEWNET)
	// if err != nil {
	// 		panic(err)
	// }

	// drop privleges immediatly from root
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
	MyLog.Printdev("[worker] all privleges droped ")

	runner_params, err := runner.GetRunnerParams() //params are passed in as json through stdin to the program
	if err != nil {
		fmt.Println("[exec error] plz pass in the correct params")
		fmt.Println("[json error] :",err)
		os.Exit(1)
	}

	MyLog.Printdev("exec worker", "[New Job] running the compilation proc")
	MyLog.Printdev("exec worker", "params = ", runner_params)

	err = os.Chdir(runner_params.CodeDir)
	if err != nil {
		panic("can not change dir error : " + err.Error())
	}

	//main start of the execution runner
	//call to compile and test
	verdict, err := runner.CompileRunAndTests(runner_params)
	if err != nil {
		panic(err)
	}

	//generating the output file
	verdict.MSG = runner.GenerateResultMSG(verdict)
	Stderr_file_path := filepath.Join(runner_params.CodeDir, runtime.StdErrorFileName)
	verdict.Stderr = utils.GetSterrFromFile(Stderr_file_path)

	json_res, err := json.Marshal(verdict)

	if err != nil {
		fmt.Println("[worker] unable to marshal")
		panic(err)
	}
	io_reader := bytes.NewBuffer(json_res)
	verdict_file_path := filepath.Join(runner_params.CodeDir, runtime.VerdictFileName)
	utils.SaveFileFromBuf(verdict_file_path, io_reader)

	MyLog.Printdev("[Verdict] : ", string(json_res))

}
