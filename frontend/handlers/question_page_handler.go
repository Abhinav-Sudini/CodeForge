package handlers

import (
	"fmt"
	"frontend/components/questionpage"
	"net/http"
	"strconv"
)

func QuestionPageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("req recived")
	question_id,err := strconv.Atoi(r.PathValue("q_id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("question not valid"))
	}
	questionpage.Questionpage(question_id).Render(r.Context(),w)
}
