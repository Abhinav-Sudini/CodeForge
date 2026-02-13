package runner

import(
	"encoding/json"
	"os"

	"worker/types"
)

func GetRunnerParams() (types.RunnerParamsJson,error) {
	var runner_parm types.RunnerParamsJson
	dec := json.NewDecoder(os.Stdin)
	err := dec.Decode(&runner_parm)
	if err != nil {
		return runner_parm, err
	}
	return runner_parm,nil
}


