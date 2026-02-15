package scraper

import (
	"io"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

const cses_link = "https://cses.fi"
const cses_task_url = "https://cses.fi/problemset/list/"
const php_session_key = "51aee256bd3e1128bef3c78e8aa0cd01ad3289af" // i am to lazy to create a .env

type Task struct {
	TaskName string `json:"TaskName"`
	TaskUrl string	`json:"TaskUrl"`
	Task_id int `jsno:"Task_id"`
	Task_Group string `json:"Task_Group"`
}


var is_li_node_check = func(n *html.Node) bool {
	if n.Type != html.ElementNode {
		return false
	}
	if n.DataAtom != atom.Li {
		return false
	}
	if len(n.Attr) != 1 {
		return false
	}
	if n.Attr[0].Key != "class" || n.Attr[0].Val != "task" {
		return false
	} 
	return true
}

var is_task_href_check = func (n *html.Node) bool {
	if n.Type != html.ElementNode {
		return false
	}	
	if n.DataAtom != atom.A {
		return false
	}
	return is_li_node_check(n.Parent)
}

func get_cses_task_html() io.ReadCloser {
	client := http.Client{}

	var body io.Reader
	req, err := http.NewRequest("GET", cses_task_url,body)
	if err != nil {
		panic(err)
	}

	// 3. Set Headers
	Util_set_header(req,cses_task_url)

	resp,err := client.Do(req)
	if err != nil {
		panic("get not get the task url")
	}
	return resp.Body
}

func get_all_task_urls(dom *html.Node) []Task {
	Output := []Task{}
	var dfs func(n *html.Node)
	dfs = func(n *html.Node){
		if is_task_href_check(n) {
			Output = append(Output,util_const_task_struct(n))
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
        dfs(c)
    }
	}
	dfs(dom)
	return Output
}

func util_const_task_struct(n *html.Node) Task {
	url := n.Attr[0].Val
	field := strings.FieldsFunc(url,func(r rune) bool {
		if r == '/' {
			return true
		}
		return false
	})
	id,_ := strconv.Atoi(field[len(field)-1])
	group := n.Parent.Parent.PrevSibling.FirstChild.Data

	return  Task{
		TaskName: n.FirstChild.Data,
		TaskUrl: url,
		Task_id: id,
		Task_Group: group,
	}
}


func Scrape_all_task() []Task {

	html_readcloser := get_cses_task_html()
	defer html_readcloser.Close()

	dom_tree,err := html.Parse(html_readcloser)
	if err != nil {
		panic("failed to parse dom")
	}

	tasks := get_all_task_urls(dom_tree)
	return tasks
}

func Util_set_header(req *http.Request,cur_url string){
	req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("accept-language", "en-US,en;q=0.9")
	req.Header.Set("cache-control", "max-age=0")
	req.Header.Set("connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("cookie", "PHPSESSID="+php_session_key)
	req.Header.Set("referer", cur_url)
	req.Header.Set("sec-ch-ua", `"Not(A:Brand";v="8", "Chromium";v="144", "Google Chrome";v="144"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"Windows"`)
	req.Header.Set("sec-fetch-dest", "document")
	req.Header.Set("sec-fetch-mode", "navigate")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("sec-fetch-user", "?1")
	req.Header.Set("upgrade-insecure-requests", "1")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36")
}
