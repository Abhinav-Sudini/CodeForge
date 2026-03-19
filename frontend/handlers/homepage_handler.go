package handlers

import (
	"fmt"
	"frontend/components/homepage"
	"frontend/config"
	"net/http"
)

func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	fmt.Println("req recived")
	homepage.Homepage(config.Question_details_api_url).Render(r.Context(), w)
}
