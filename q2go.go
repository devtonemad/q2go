package main

import (
	"container/list"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

var queueMap map[string]*list.List
var router *mux.Router
var queueMu sync.RWMutex

func main() {
	initialize()
	log.Println("q2go started on :8077")
	log.Fatal(http.ListenAndServe(":8077", router))
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
	if qname == "" {
		http.Error(writer, "qname required", http.StatusBadRequest)
		return
	}
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
	q := getQueue(queueMap, qname)
	if q == nil {
		writer.WriteHeader(404)
	} else {
		pushMessage(q, m)
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
		return
	}
	queueMu.Lock()
	delete(queueMap, qname)
	queueMu.Unlock()
}

func createQueue(qm map[string]*list.List, qname string) *list.List {
	queueMu.Lock()
	defer queueMu.Unlock()
	if existing := qm[qname]; existing != nil {
		return existing
	}
	q := list.New()
	qm[qname] = q
	return q
}

func pushMessage(q *list.List, msg string) {
	queueMu.Lock()
	q.PushBack(msg)
	queueMu.Unlock()
}

func popMessage(qname string) (string, error) {

	var message string
	queueMu.Lock()
	defer queueMu.Unlock()
	q := getQueue(queueMap, qname)
	if q == nil {
		return "", fmt.Errorf("queue with the given name %s not found ", qname)
	}

	if q.Len() > 0 {
		e := q.Front()
		message = e.Value.(string)
		q.Remove(e)
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
