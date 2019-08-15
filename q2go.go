package main

import (
	"container/list"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)
/* test */
var queue = list.New()

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/pop", popnextmessage).Methods("GET")
	router.HandleFunc("/push", pushmessage).Methods("POST")
	http.ListenAndServe(":8080", router)
}

func popnextmessage(writer http.ResponseWriter, request *http.Request) {
	if queue.Len() > 0 {
		e := queue.Front()
		m := e.Value.(string)
		fmt.Printf("message removed from queue: %s \n", m)
		queue.Remove(e)
		writer.Write([]byte(m))
	} else {
		fmt.Printf("no more messages in queue \n")
	}
}

func pushmessage(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	m := request.FormValue("message")
	fmt.Printf("message pushed to queue   : %s \n", m)
	queue.PushBack(m)
}
