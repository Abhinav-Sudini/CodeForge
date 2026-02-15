package main

import (
	"fmt"
	scraper "worker/CSES_Scraper"
)


func main(){
	// scraper.Test_html_get_home()
	tasks := scraper.Scrape_all_task()
	fmt.Println(tasks)
}
