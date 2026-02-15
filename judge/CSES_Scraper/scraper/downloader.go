package scraper

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func Get_csrf_token(url string) (string,error) {
	client := http.Client{}

	var body io.Reader
	req, err := http.NewRequest("GET", url,body)
	if err != nil {
		return "",errors.New("could not make req")
	}

	Util_set_header(req,url)

	// dump, err := httputil.DumpRequest(req, true)
	// if err != nil {
	// 	panic("bye")
	// }
	// fmt.Printf("%s", dump)

	resp,err := client.Do(req)
	if err != nil {
		return "",errors.New("could not fetch data")
	}
	defer resp.Body.Close()


	dom,err := html.Parse(resp.Body)
	if err != nil {
		fmt.Println("dome parse failed")
		return "",err
	}
	token := Get_csrf_token_from_dom(dom)
	if token == "" {
		return "idk",errors.New("no token")
	}

	return token,nil
}

func Get_csrf_token_from_dom(dom *html.Node) (string) {
	token := ""
	var dfs func(n *html.Node)
	dfs = func(n *html.Node){
		if is_csrf_token(n) {
			token = n.Attr[2].Val
			return 
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
        dfs(c)
				if token != "" {
					break
				}
    }
	}
	dfs(dom)
	return token
}

func is_csrf_token(n *html.Node) bool {
	if n.Type != html.ElementNode {
		return false
	}
	if 	n.Parent==nil || n.Parent.DataAtom != atom.Form || len(n.Attr)!=3 {
		return false
	}
	return true
	
}

func Get_and_save_file(pbl_dir string , task Task) error {
	down_url := "https://cses.fi/problemset/tests/" + strconv.Itoa(task.Task_id) + "/"
	token,err := Get_csrf_token(down_url)
	if err != nil {
		return err
	}

	data := url.Values{}
	data.Set("csrf_token", token)
	data.Set("download", "true")

	client := http.Client{}
	req,err := http.NewRequest("POST",down_url,strings.NewReader(data.Encode()))
	

	if err != nil {
		return err
	}
	Util_set_header(req,down_url)

	// dump, err := httputil.DumpRequest(req, true)
	// if err != nil {
	// 	panic("bye")
	// }
	// fmt.Printf("%s", dump)
	resp,err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	zip_name := filepath.Join(pbl_dir,"prob.zip")
	zip,err := os.Create(zip_name)
	if err != nil {
		return err
	}
	defer zip.Close()

	_, err = io.Copy(zip, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
