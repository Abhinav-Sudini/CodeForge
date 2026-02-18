package handlers

import (
	"fmt"
	"worker/runner"
)


func PostResponseToMaster(verdict runner.SubmitionResult) error {
	//TODO
	fmt.Println("verdict from server",verdict)
	return nil
}
