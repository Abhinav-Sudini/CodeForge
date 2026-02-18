package MyLog

import (
	"fmt"
	"worker/runtime"
)


const is_production = runtime.Is_production

var init_done = false


func Printdev(exec_loc string,a ...any) {
	if init_done == false {
		init_logger()
	}
	if is_production {
		return
	}
	fmt.Print("[",exec_loc,"] : ")
	fmt.Println(a...)
}

func Print(exec_loc string,a ...any) {
	if init_done == false {
		init_logger()
	}
	fmt.Print("[",exec_loc,"] : ")
	fmt.Println(a...)
}

func init_logger(){

	init_done = true
}


