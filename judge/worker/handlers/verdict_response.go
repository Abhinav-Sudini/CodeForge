package handlers

import (
	"fmt"
	"worker/types"
)

func PostResponseToMaster(verdict types.JudgeCodeResponse) error {
	//TODO
	fmt.Println("verdict from server", verdict)
	return nil
}

func PostWorkerIsFreeReqToMaster() error {
	fmt.Println("posting is free ")
	return nil
}
