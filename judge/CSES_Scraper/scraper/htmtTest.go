package scraper

import (
	"fmt"
	"io"
)

func Test_html_get_home(){
	html_readcloser := get_cses_task_html()
	defer html_readcloser.Close()
	html,err := io.ReadAll(html_readcloser)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(html))
}
