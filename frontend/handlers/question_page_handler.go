package handlers

import (
	"fmt"
	"frontend/components/questionpage"
	"frontend/config"
	"net/http"
	"strconv"
)

func QuestionPageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("req recived")
	question_id, err := strconv.Atoi(r.PathValue("q_id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("question not valid"))
	}
	questionpage.Questionpage(question_id,
		config.Question_details_api_url, config.Get_question_submissions_api, config.Code_submission_api, config.Submission_verdict_api,
	).Render(r.Context(), w)
}
