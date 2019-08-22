package main

import (
	"container/list"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

var queueMap map[string]*list.List
var router *mux.Router

func main() {
	initialize()
	http.ListenAndServe(":8080", router)
}

func initialize() {
	router = mux.NewRouter()
	queueMap = make(map[string]*list.List)
	router.HandleFunc("/queue", createQueueHandler).Methods("POST").Name("createQueue")
	router.HandleFunc("/queue/{qid}/push", pushMessageHandler).Methods("POST").Name("push")
	router.HandleFunc("/queue/{qid}/pop", popMessageHandler).Methods("GET").Name("pop")
}

func createQueueHandler(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	qname := request.FormValue("qname")
	createQueue(queueMap, qname)
	writer.Write([]byte(qname))
}

func pushMessageHandler(writer http.ResponseWriter, request *http.Request) {
	v := mux.Vars(request)
	qname := v["qid"]
	request.ParseForm()
	m := request.FormValue("message")
	q := getQueue(queueMap, qname)
	go pushMessage(q, m)
}

func popMessageHandler(writer http.ResponseWriter, request *http.Request) {
	v := mux.Vars(request)
	qname := v["qid"]
	message := popMessage(qname)
	writer.Write([]byte(message))
}

func createQueue(qm map[string]*list.List, qname string) *list.List {
	//check if name already exists
	q := list.New()
	qm[qname] = q
	return q
}

func pushMessage(q *list.List, msg string) {
	q.PushBack(msg)
	//some time consuming process
	rtd := time.Duration(rand.Intn(500))
	time.Sleep(time.Millisecond * rtd)
	fmt.Printf("message pushed to queue   : %s \n", msg)
}

func popMessage(qname string) string {

	var message string
	q := getQueue(queueMap, qname)

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

func getQueue(qm map[string]*list.List, qname string) *list.List {
	q := qm[qname]
	if q == nil {
		fmt.Printf("queue with the name %s does not exist \n", qname)
	}
	return q
}
