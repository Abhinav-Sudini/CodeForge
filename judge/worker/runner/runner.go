package runner

import (
	"errors"
	"os"
	"worker/types"
)

type Runner interface{
	RunJobAndGetResult(req *types.JudgeCodeRequest) types.JudgeCodeReqVerdict 
}

type runner struct{
	runtime string
}

func NewRunner(requested_runtime string) (Runner,error) {
	worker_runtime := os.Getenv("WORKER_ENV_CF")
	if worker_runtime == "" || worker_runtime != requested_runtime{
		return nil,errors.New("[runner]: need a Runtime or not a valid runtime")
	}
	run := &runner{
		runtime: requested_runtime,
	}
	return Runner(run),nil
}

func (run *runner)RunJobAndGetResult(req *types.JudgeCodeRequest) types.JudgeCodeReqVerdict {
	//TODO
	return types.JudgeCodeReqVerdict{}
}
