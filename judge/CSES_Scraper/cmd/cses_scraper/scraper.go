package main

import (
	"fmt"
	scraper "scraper/scraper"
)

const out_put_json_file = "CSES_Scraped_Task_Details.json"

func main(){
	// scraper.Test_html_get_home()
	scraped_tasks := scraper.Scrape_all_task()
	err := scraper.Save_tasksJson_to_file(out_put_json_file,scraped_tasks)
	if err != nil{
		panic(err)
	}
	
	fmt.Println(scraped_tasks)
}
