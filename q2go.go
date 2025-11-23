package main

import (
	"container/list"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

type Queue struct {
	mu sync.Mutex
	l  *list.List
}

func NewQueue() *Queue {
	return &Queue{l: list.New()}
}

func (q *Queue) Push(msg string) {
	q.mu.Lock()
	q.l.PushBack(msg)
	q.mu.Unlock()
}

func (q *Queue) Pop() (string, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.l.Len() == 0 {
		return "", false
	}
	e := q.l.Front()
	msg := e.Value.(string)
	q.l.Remove(e)
	return msg, true
}

func (q *Queue) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.l.Len()
}

var queueMap map[string]*Queue
var router *mux.Router
var queueMu sync.RWMutex

func main() {
	initialize()
	log.Println("q2go started on :8077")
	log.Fatal(http.ListenAndServe(":8077", router))
}

func initialize() {
	router = mux.NewRouter()
	queueMap = make(map[string]*Queue)
	router.HandleFunc("/queue", queuePostHandler).Methods("POST").Name("queuePost")
	router.HandleFunc("/queue/{qid}", queueDeleteHandler).Methods("DELETE").Name("queueDelete")
	router.HandleFunc("/queue/{qid}/message", messagePostHandler).Methods("POST").Name("messagePost")
	router.HandleFunc("/queue/{qid}/message", messageGetHandler).Methods("GET").Name("messageGet")
	router.HandleFunc("/status", statusHandler).Methods("GET").Name("status")

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
		q.Push(m)
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

func createQueue(qm map[string]*Queue, qname string) *Queue {
	queueMu.Lock()
	defer queueMu.Unlock()
	if existing := qm[qname]; existing != nil {
		return existing
	}
	q := NewQueue()
	qm[qname] = q
	return q
}

// push/pop are methods on Queue now

func popMessage(qname string) (string, error) {
	q := getQueue(queueMap, qname)
	if q == nil {
		return "", fmt.Errorf("queue with the given name %s not found ", qname)
	}
	if msg, ok := q.Pop(); ok {
		return msg, nil
	}
	return "", nil
}

func getQueue(qm map[string]*Queue, qname string) *Queue {
	queueMu.RLock()
	q := qm[qname]
	queueMu.RUnlock()
	if q == nil {
		fmt.Printf("queue with the name %s does not exist \n", qname)
	}
	return q
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	type statusResp struct {
		Queues   int `json:"queues"`
		Messages int `json:"messages"`
	}

	queueMu.RLock()
	queues := len(queueMap)
	// make a copy of map values to avoid holding RLock while calling Len()
	qs := make([]*Queue, 0, queues)
	for _, q := range queueMap {
		qs = append(qs, q)
	}
	queueMu.RUnlock()

	total := 0
	for _, q := range qs {
		total += q.Len()
	}

	resp := statusResp{Queues: queues, Messages: total}
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(resp); err != nil {
		http.Error(w, "failed to encode status", http.StatusInternalServerError)
		return
	}
}
