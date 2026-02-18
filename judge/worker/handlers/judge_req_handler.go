package handlers

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"

	MyLog "worker/logger"
	"worker/runner"
	"worker/types"
	"worker/utils"
)

var WorkerMutex sync.Mutex


func validateJudgeReq(r types.JudgeCodeRequest) error {
    if r.Runtime == "" || r.Code == "" {
        return errors.New("[json vaidate] code/runtime is required")
    }
    if r.TimeConstrain <= 0 || r.MemConstrain <= 0 || r.JobId <= 0 {
        return errors.New("[json validate] invalid constraints")
    }
    return nil
}

func create_new_job(stream io.ReadCloser) error {
	var judgeReq types.JudgeCodeRequest
	err := utils.CustomJsonUnMarshal(stream,&judgeReq)
	if err != nil {
		return err
	}
	err = validateJudgeReq(judgeReq)
	if err != nil {
		return err
	}
	fmt.Println(judgeReq)

	if WorkerMutex.TryLock() == true {

		//main entry point of our code runner it compiles and 
		//runs the code async and post the verdict to the Master 
		//to the url of the master with http
		go func() {
			defer WorkerMutex.Unlock()
			//runing the job and posting the verdict
			run,err := runner.NewRunner(judgeReq.Runtime)
			if err != nil {
				return
			}
			Result,err := run.RunJobAndGetResult(&judgeReq)
			if err!=nil {
				MyLog.Print("exec parent","faile with",err)
				panic(err)
			}
			err = PostResponseToMaster(Result)
			if err != nil {
				panic(err.Error())
			}
			MyLog.Printdev("server liste for another req")
		}()

	}else{
		return errors.New("[worker] worker busy ")
	}
	 

	return  nil
}

func Compile_and_judge_handler(w http.ResponseWriter, r *http.Request){
	err := create_new_job(r.Body)
	if err != nil {
		MyLog.Print("worker","worker faied with error : ",err)
		http.Error(w,err.Error(),http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("job queued"))
}
