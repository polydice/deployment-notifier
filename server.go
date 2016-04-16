package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/google/go-github/github"
	"github.com/julienschmidt/httprouter"
)

func main() {
	route := httprouter.New()
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	route.POST("/webhook", DeploymentHandler)

	fmt.Println("Starting server on :" + port)

	http.ListenAndServe(":"+port, route)
}

func DeploymentHandler(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
	event_type := req.Header.Get("X-Github-Event")
	switch event_type {
	case "deployment":
		event := decodeDeploymentEvent(req)
		fmt.Printf("event: %#v", event)
	case "deployment_status":
		event := decodeDeploymentStatusEvent(req)
		fmt.Printf("event: %#v", event)
	}

}

func decodeDeploymentStatusEvent(req *http.Request) github.DeploymentStatusEvent {
	decoder := json.NewDecoder(req.Body)

	var event github.DeploymentStatusEvent

	err := decoder.Decode(&event)

	if err != nil {
		fmt.Print(err)
	}

	return event
}

func decodeDeploymentEvent(req *http.Request) github.DeploymentEvent {
	decoder := json.NewDecoder(req.Body)

	var event github.DeploymentEvent

	err := decoder.Decode(&event)

	if err != nil {
		fmt.Print(err)
	}

	return event
}
