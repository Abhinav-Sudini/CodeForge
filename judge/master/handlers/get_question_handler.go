package handlers

import (
	"encoding/json"
	"master/types"
	"net/http"
	"strconv"
)

func (s *Server)GetQuestionHandler(w http.ResponseWriter,r *http.Request){

	question_id,err := strconv.Atoi(r.PathValue("q_id"))
	if err != nil {
		http.Error(w,"no question pathvalue",http.StatusBadRequest)
		return
	}
	question,err :=  s.queries.GetQuestion(r.Context(),int32(question_id))
	if err != nil {
		http.Error(w,"failed to get questions",http.StatusInternalServerError)
		return
	}
	var return_struct = types.Get_question_resp_json{
		QuestionId: question.QuestionID,
		QuestionCategory: question.QuestionCategory.String,
		QuestionName: question.QuestionName.String,
		QuestionDescription: question.QuestionDescription.String,
		InputDescription: question.InputDescription.String,
		OutputDescription: question.OutputDescription.String,
		ConstraintsDescription: question.ConstraintsDescription.String,
		TimeConstraint: question.TimeConstraint,
		MemConstraint: question.MemConstraint,
		ExampleInputs: question.ExampleInputs,
		ExampleOutputs: question.ExampleOutputs,
	}
	json_resp,err := json.Marshal(return_struct)
	if err != nil {
		http.Error(w,"failed to get response",http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(json_resp)
}
