package scraper

import (
	"encoding/json"
	"os"
)

func Save_tasksJson_to_file(path string,task_list []Task) error {
	
	f,err := os.OpenFile(path,os.O_CREATE|os.O_TRUNC|os.O_WRONLY,0644)
	if err != nil {
		return err
	}
	
	task_json,err := json.Marshal(task_list)
	if err != nil {
		return err
	}

	f.Write(task_json)
	return nil
}


