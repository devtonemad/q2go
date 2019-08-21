package main

import (
	"container/list"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

var queue *list.List

func main() {
	router := getRouter()
	queue = getQueue()
	initRouteHandlers(router)
	http.ListenAndServe(":8080", router)

}

func initRouteHandlers(r *mux.Router) {
	r.HandleFunc("/pop", popNextMessage).Methods("GET").Name("pop")
	r.HandleFunc("/push", pushMessage).Methods("POST").Name("push")
}

func getRouter() *mux.Router {
	return mux.NewRouter()
}

func getQueue() *list.List {
	return list.New()
}

func popNextMessage(writer http.ResponseWriter, request *http.Request) {
	message := popNextMessageFromQueue(queue)
	writer.Write([]byte(message))
}

func popNextMessageFromQueue(q *list.List) string {
	var message string
	if q.Len() > 0 {
		e := q.Front()
		message = e.Value.(string)
		fmt.Printf("message removed from queue: %s \n", message)
		q.Remove(e)
	} else {
		fmt.Printf("no more messages in queue \n")
	}
	return message
}

func pushMessage(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	m := request.FormValue("message")
	go pushMessageToQueue(queue, m)
}

func pushMessageToQueue(q *list.List, msg string) {
	q.PushBack(msg)
	//some time consuming process
	rtd := time.Duration(rand.Intn(500))
	time.Sleep(time.Millisecond * rtd)
	fmt.Printf("message pushed to queue   : %s \n", msg)
}
