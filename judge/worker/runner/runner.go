package runner

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"worker/runtime"
	"worker/types"
	"worker/utils"
)

type Runner interface {
	RunJobAndGetResult(req *types.JudgeCodeRequest) (SubmitionResult, error)
}

type runner struct {
	runtime string
}

func NewRunner(requested_runtime string) (Runner, error) {
	worker_runtime := os.Getenv("WORKER_RUNTIME")
	if worker_runtime == "" || worker_runtime != requested_runtime {
		return nil, errors.New("[runner]: need a Runtime or not a valid runtime")
	}
	run := &runner{
		runtime: requested_runtime,
	}
	return Runner(run), nil
}

func (run *runner) RunJobAndGetResult(req *types.JudgeCodeRequest) (SubmitionResult, error) {
	//TODO

	codeFileDir := os.Getenv("CODE_FILE_DIR")
	allTasksDir := os.Getenv("QUESTIONS_DIR")
	codeFileName := runtime.GetCodeFileName(req.Runtime)
	codeFilePath := filepath.Join(codeFileDir, codeFileName)

	err := utils.RemoveAllFilesInDir(codeFileDir)
	if err != nil {
		return SubmitionResult{},err
	}

	err = os.WriteFile(codeFilePath, []byte(req.Code), 0666)
	if err != nil {
		return SubmitionResult{}, err
	}

	runnerParams := types.RunnerParamsJson{
		CodeDir:       codeFileDir,
		TestCasesDir:  filepath.Join(allTasksDir, strconv.Itoa(req.QuestionId)),
		Runtime:       req.Runtime,
		TimeConstrain: req.TimeConstrain,
		MemConstrain:  req.MemConstrain,
	}
	err = runExecWorkerProcessWithParams(runnerParams)
	if err != nil {
		fmt.Println("[cmd exec] : ", err)
		return SubmitionResult{}, err
	}

	result_file := filepath.Join(codeFileDir, runtime.VerdictFileName)
	Result, err := getSubmitionResultFromFile(result_file)
	if err != nil {
		return Result, err
	}

	return Result, nil
}

func runExecWorkerProcessWithParams(runnerParams types.RunnerParamsJson) error {
	ctx := context.Background()

	cmd := exec.CommandContext(ctx, runtime.ExecWorkerStartComand)
	var buf_in, buf_out, buf_err bytes.Buffer
	cmd.Stdin = &buf_in
	cmd.Stdout = &buf_out
	cmd.Stderr = &buf_err

	params_json_bytes, _ := json.Marshal(runnerParams)
	buf_in.Write(params_json_bytes)

	if err := cmd.Run(); err != nil {
		t, _ := io.ReadAll(&buf_out)
		fmt.Println("[ecec work cout]", string(t))
		t, _ = io.ReadAll(&buf_err)
		fmt.Println("[ecec work cerr]", string(t))
		return err
	}
	if runtime.Is_production == false {
		t, _ := io.ReadAll(&buf_out)
		fmt.Println("[ecec work cout]", string(t))
		t, _ = io.ReadAll(&buf_err)
		fmt.Println("[ecec work cerr]", string(t))
	}
	return nil
}

func getSubmitionResultFromFile(path string) (SubmitionResult, error) {
	f, err := os.Open(path)
	if err != nil {
		return SubmitionResult{}, err
	}
	json_str, err := io.ReadAll(f)
	var result SubmitionResult
	err = json.Unmarshal(json_str, &result)
	if err != nil {
		return SubmitionResult{}, err
	}
	return result, nil
}
