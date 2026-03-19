package handlers

import "net/http"


func ServeStaticPage(w http.ResponseWriter,r *http.Request){
	http.ServeFile(w,r,"./static/html/")
}
