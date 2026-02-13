package main

import (
	"errors"
	"fmt"
	"os/user"
	"strconv"
	"syscall"
	"worker/runner"
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

	//drop privleges immediatly from root
	err = syscall.Unshare(syscall.CLONE_NEWNET)
	if err != nil {
			panic(err)
	}
	fmt.Println("test 1")

	//Drop all groups
	if err := syscall.Setgroups([]int{}); err != nil {
			panic(err)
	}

	fmt.Println("test 1")
	// Drop GID first
	if err := syscall.Setgid(nobodyGID); err != nil {
			panic(err)
	}

	fmt.Println("test 3")
	// Drop UID
	if err := syscall.Setuid(nobodyUID); err != nil {
			panic(err)
	}
	
	runner_params,err := runner.GetRunnerParams() //params are passed in as json through stdin
	if err != nil {
		fmt.Println("[exec error] plz pass in the correct params")
	}

	fmt.Println("running the compilation proc")
	fmt.Println(runner_params)
	verdict,err := runner.CompileRunAndTests(runner_params)
	if err != nil {
		panic(err)
	}
	fmt.Println("[Verdict] : ",verdict)

}
