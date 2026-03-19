package handlers

import (
	"fmt"
	"frontend/components/homepage"
	"net/http"
)

func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	fmt.Println("req recived")
	homepage.Homepage().Render(r.Context(), w)
}
