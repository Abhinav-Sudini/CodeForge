package scraper

import (
	"io"
	"net/http"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

const cses_link = "https://cses.fi"
const cses_task_url = "https://cses.fi/problemset/list/"
const phh_session_key = "51aee256bd3e1128bef3c78e8aa0cd01ad3289af"  // i am to lazy to create a .env

type Task struct {
	TaskName string
	TaskUrl string
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
	util_set_header(req,cses_task_url)

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
	return  Task{
		TaskName: n.FirstChild.Data,
		TaskUrl: n.Attr[0].Val,
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

func util_set_header(req *http.Request,cur_url string){
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36")
	req.Header.Set("Cookie", "PHPSESSID="+phh_session_key)
	req.Header.Set("Origin", "https://cses.fi")
	req.Header.Set("Referer", cur_url)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8")
}
