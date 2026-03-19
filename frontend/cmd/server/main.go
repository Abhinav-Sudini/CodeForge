package main

import (
	"fmt"
	"frontend/config"
	"frontend/handlers"
	"net/http"
	"strconv"
)


func main(){
	http.HandleFunc("/",handlers.HomePageHandler)
	http.HandleFunc("/problems/{q_id}/",handlers.QuestionPageHandler)
	fmt.Println("server start at port - ",config.SERVER_PORT)
	http.ListenAndServe(":"+strconv.Itoa(config.SERVER_PORT),nil)
}

