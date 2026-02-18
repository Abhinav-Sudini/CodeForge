package utils

import "syscall"


func GetTimeInMillSec(timeval syscall.Timeval) int64 {
	return timeval.Nano()/1000
}


