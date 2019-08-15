package main

import "testing"

func TestRouter(t *testing.T) {
	router := getRouter()
	if router == nil {
		t.Errorf("No router created")
	}
}

func TestInitRouterHandlers(t *testing.T) {
	router := getRouter()
	initRouteHandlers(router)

	route := router.Get("pop")
	if route == nil {
		t.Errorf("No route created for pop")
	}

	route = router.Get("push")
	if route == nil {
		t.Errorf("No route created for push")
	}

}

func TestPushMessage(t *testing.T) {
	message := "test"
	queue := getQueue()
	pushMessageToQueue(queue, message)
	if queue.Len() == 0 {
		t.Errorf("No message pushed in queue")
	}
}

func TestPopMessage(t *testing.T) {
	queue := getQueue()
	pushMessageToQueue(queue, "test")
	m := popNextMessageFromQueue(queue)
	if m != "test" {
		t.Errorf("No message popped from queue")
	}

}
