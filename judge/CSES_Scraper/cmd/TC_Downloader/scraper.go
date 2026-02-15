package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	scraper "scraper/scraper"
	"strconv"
)

const download_dir = "./scraped_testcases"
const input_put_json_file = "CSES_Scraped_Task_Details.json"

func main(){
	// scraper.Test_html_get_home()
	

	f,err := os.Open(input_put_json_file)
	if err != nil {
		panic(err)
	}
	json_tasks,err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}

	var all_tasks []scraper.Task
	err = json.Unmarshal(json_tasks,&all_tasks)
	if err != nil{
		panic(err)
	}

	for _,task := range(all_tasks) {
		id_str:= strconv.Itoa(task.Task_id)
		pbl_dir := filepath.Join(download_dir,id_str)
		err := os.MkdirAll(pbl_dir,0777)
		if err != nil {
			fmt.Println("failet to ccreate file - ",pbl_dir,"with err : ",err)
		}
		err = scraper.Get_and_save_file(pbl_dir,task)
		if err != nil {
			panic(err)
		}
		fmt.Println("[done] for ",task.TaskName)
	} 
	
	fmt.Println("done bro")
}


