package utils

import "syscall"


func GetTimeInMillSec(tv syscall.Timeval) int64 {
	return int64(tv.Sec)*1e3 + int64(tv.Usec)/1000
}


