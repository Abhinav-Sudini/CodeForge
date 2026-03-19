package scraper

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"scraper/db1/postgres_db"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

var (
	json_file_path = "./CSES_Scraped_Task_Details.json"
	start_ind      = 275 + 91
	end_ind        = 274
)

func Scrape_all_questions_and_add_to_db() {

	json_list, err := os.ReadFile(json_file_path)
	if err != nil {
		panic(err)
	}
	var all_tasks []Task
	err = json.Unmarshal(json_list, &all_tasks)
	if err != nil {
		panic(err)
	}
	// PG_DATABASE_URL=postgres://postgres:dev_password@database:5432/dev_db
	db_con, err := pgx.Connect(context.Background(), os.Getenv("PG_DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	quries := postgres_db.New(db_con)

	for it, task := range all_tasks[start_ind:] {
		params, err := scrape_question(task.Task_id, task)
		if err != nil {
			fmt.Println("failed to get question ", task.Task_id, "at it :", it)
			return
		}
		err = quries.CreateQuestion(context.Background(), params)
		if err != nil {
			fmt.Println("failed to get question ", task.Task_id, "at it :", it, "err: ", err)
			return
		}
		fmt.Println("added to db ", it, "qid", task.Task_id)
	}

}

func scrape_question(question_id int, task Task) (postgres_db.CreateQuestionParams, error) {
	html_readcloser := get_question_page(question_id)
	defer html_readcloser.Close()

	dom_tree, err := html.Parse(html_readcloser)
	if err != nil {
		panic("failed to parse dom")
	}

	var params postgres_db.CreateQuestionParams

	var tm_con float32
	var mem_con int

	var dfs func(n *html.Node)
	dfs = func(n *html.Node) {

		if is_time_constraint(n) {
			fmt.Sscanf(n.Data, " %f s", &tm_con)
			// fmt.Println(n.Data,"t :",t,"n: ",n,"err ",err)
			params.TimeConstraint = int32(tm_con)
		}
		if is_mem_constraint(n) {
			fmt.Sscanf(n.Data, " %d MB", &mem_con)
			// fmt.Println(n.Data,"t :",t,"n: ",n,"err ",err)
			params.TimeConstraint = int32(mem_con)
		}
		if is_md_node(n) {
			c := n.FirstChild
			var question string
			for c.DataAtom != atom.H1 {
				question += nodeText(c)
				c = c.NextSibling
			}
			var input string
			if c.FirstChild.Data == "Input" {
				c = c.NextSibling
				for c.DataAtom != atom.H1 {
					input += nodeText(c)
					c = c.NextSibling
				}
			}

			var output string
			if c.FirstChild.Data == "Input" || c.FirstChild.Data[:6] == "Output" {
				c = c.NextSibling
				for c.DataAtom != atom.H1 {
					output += nodeText(c)
					c = c.NextSibling
				}
			}

			var cons string
			if c.FirstChild.Data == "Constraints" {
				c = c.NextSibling
				for c.DataAtom != atom.H1 {
					cons += nodeText(c)
					c = c.NextSibling
				}
			}

			var exp_inp, exp_out []string
				for c != nil && c.DataAtom == atom.H1 {
					for c.DataAtom != atom.Pre {
						c = c.NextSibling
					}
					exp_inp = append(exp_inp, c.FirstChild.Data)
					c = c.NextSibling
					for c.DataAtom != atom.Pre {
						c = c.NextSibling
					}
					exp_out = append(exp_out, c.FirstChild.Data)
					c = c.NextSibling
				}
			if question == "" || input == "" || output == "" || len(exp_inp) == 0 || len(exp_out) == 0 || tm_con == 0 || mem_con == 0 {
				fmt.Println("failed to get question: ", question_id, "task id", task.Task_id)
				fmt.Println(question, input, output, cons, exp_inp, exp_out)
				panic("can not find data")
			}
			params = postgres_db.CreateQuestionParams{
				QuestionID:             int32(task.Task_id),
				QuestionCategory:       pgtype.Text{task.Task_Group, true},
				QuestionName:           pgtype.Text{task.TaskName, true},
				QuestionDescription:    pgtype.Text{question, true},
				InputDescription:       pgtype.Text{input, true},
				OutputDescription:      pgtype.Text{output, true},
				ConstraintsDescription: pgtype.Text{cons, true},
				TimeConstraint:         int32(tm_con),
				MemConstraint:          int32(mem_con),
				ExampleInputs:          exp_inp,
				ExampleOutputs:         exp_out,
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			dfs(c)
		}
	}

	dfs(dom_tree)

	fmt.Println("done for question: ", task.Task_id, "tm", params.TimeConstraint, "len: ", len(params.ExampleInputs))
	// fmt.Println(params)
	return params, nil
}

func is_md_node(n *html.Node) bool {
	if n.Type != html.ElementNode {
		return false
	}
	if n.DataAtom != atom.Div || len(n.Attr) != 1 {
		return false
	}
	if n.Attr[0].Val != "md" {
		return false
	}

	return true
}

func nodeText(n *html.Node) string {
	var buf bytes.Buffer
	if err := html.Render(&buf, n); err != nil {
		return ""
	}
	return buf.String()
	//
	// if n == nil {
	// 	return ""
	// }
	// var b strings.Builder
	//
	// var walk func(*html.Node)
	// walk = func(cur *html.Node) {
	// 	if cur.Type == html.TextNode {
	// 		b.WriteString(cur.Data)
	// 	}
	// 	for c := cur.FirstChild; c != nil; c = c.NextSibling {
	// 		walk(c)
	// 	}
	// }
	//
	// walk(n)
	// return b.String()
}

func is_mem_constraint(n *html.Node) bool {
	if n==nil || n.Type != html.TextNode {
		return false
	}
	if n.Parent==nil || n.Parent.DataAtom != atom.Li {
		return false
	}
	if n.PrevSibling==nil || n.PrevSibling.FirstChild==nil || n.PrevSibling.FirstChild.Data != "Memory limit:" {
		return false
	}
	return true
}
func is_time_constraint(n *html.Node) bool {
	if n==nil || n.Type != html.TextNode {
		return false
	}
	if n.Parent==nil || n.Parent.DataAtom != atom.Li {
		return false
	}
	if n.PrevSibling==nil || n.PrevSibling.FirstChild==nil || n.PrevSibling.FirstChild.Data != "Time limit:" {
		return false
	}
	return true
}

func get_question_page(question_id int) io.ReadCloser {
	url := fmt.Sprintf("https://cses.fi/problemset/task/%d/", question_id)
	client := http.Client{}

	var body io.Reader
	req, err := http.NewRequest("GET", url, body)
	if err != nil {
		panic(err)
	}

	// 3. Set Headers
	Util_set_header(req, url)

	resp, err := client.Do(req)
	if err != nil {
		panic("get not get the task url")
	}
	return resp.Body
}
