package main

import (
	"fmt"
	"os"
	"os/exec"
	"worker/runtime"
)

func main(){
	run,_ := runtime.GetRuntime("c++17")
	cpl := run.CompileComand[0]
	arg := run.CompileComand[1:]
	fmt.Println(cpl,arg)
	cmd := exec.Command(cpl,arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err:=cmd.Run();err!=nil {
		panic(err)
	}
	fmt.Println("done compiling")
}
