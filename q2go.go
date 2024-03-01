package main

import (
	"container/list"
	"fmt"
	"io"
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
	router.HandleFunc("/queue", queuePostHandler).Methods("POST").Name("queuePost")
	router.HandleFunc("/queue/{qid}", queueDeleteHandler).Methods("DELETE").Name("queueDelete")
	router.HandleFunc("/queue/{qid}/message", messagePostHandler).Methods("POST").Name("messagePost")
	router.HandleFunc("/queue/{qid}/message", messageGetHandler).Methods("GET").Name("messageGet")

}

func queuePostHandler(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	qname := request.FormValue("qname")
	createQueue(queueMap, qname)
	writer.Write([]byte(qname))
}

func messagePostHandler(writer http.ResponseWriter, request *http.Request) {
	v := mux.Vars(request)
	qname := v["qid"]
	bodyBytes, err := io.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, "Error Reading Request Body", http.StatusInternalServerError)
		return
	}
	m := string(bodyBytes)
	//request.ParseForm()
	//m := request.FormValue("message")
	q := getQueue(queueMap, qname)
	if q == nil {
		writer.WriteHeader(404)
	} else {
		go pushMessage(q, m)
	}
}

func messageGetHandler(writer http.ResponseWriter, request *http.Request) {
	v := mux.Vars(request)
	qname := v["qid"]
	message, err := popMessage(qname)
	if err != nil {
		writer.WriteHeader(404)
	} else {
		writer.Write([]byte(message))
	}
}

func queueDeleteHandler(writer http.ResponseWriter, request *http.Request) {
	v := mux.Vars(request)
	qname := v["qid"]
	q := getQueue(queueMap, qname)
	if q == nil {
		writer.WriteHeader(404)
	}
	delete(queueMap, qname)
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

func popMessage(qname string) (string, error) {

	var message string
	q := getQueue(queueMap, qname)
	if q == nil {
		return "", fmt.Errorf("queue with the given name %s not found ", qname)
	}

	if q.Len() > 0 {
		e := q.Front()
		message = e.Value.(string)
		fmt.Printf("message removed from queue: %s \n", message)
		q.Remove(e)
	} else {
		fmt.Printf("no more messages in queue \n")
	}

	return message, nil

}

func getQueue(qm map[string]*list.List, qname string) *list.List {
	q := qm[qname]
	if q == nil {
		fmt.Printf("queue with the name %s does not exist \n", qname)
	}
	return q
}
