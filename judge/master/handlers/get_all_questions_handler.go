package handlers

import (
	"encoding/json"
	db "master/db/postgres_db"
	"master/types"
	"net/http"
)

func (s *Server) GetAllQuestionsDetailsHandler(w http.ResponseWriter, r *http.Request) {
	all_questions, err := s.queries.GetAllQuestionsMinimalDetails(r.Context())
	if err != nil {
		http.Error(w, "unalble to get questions", http.StatusInternalServerError)
		return
	}
	var temp_slice = make_question_struct(all_questions)

	resp_json, err := json.Marshal(temp_slice)
	if err != nil {
		http.Error(w, "unalble to get questions", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resp_json)
}

func make_question_struct(all_questions []db.GetAllQuestionsMinimalDetailsRow) types.Get_all_questions_resp_json {
	var temp_slice []types.Question_record_minimal
	for i := range all_questions {
		temp_slice = append(temp_slice, types.Question_record_minimal{
			QuestionId:       int(all_questions[i].QuestionID),
			QuestionName:     all_questions[i].QuestionName.String,
			QuestionCategory: all_questions[i].QuestionCategory.String,
		})
	}
	return types.Get_all_questions_resp_json{
		Questions: temp_slice,
	}
}
