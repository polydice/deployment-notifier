package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/google/go-github/github"
	"github.com/julienschmidt/httprouter"
)

type Event interface{}

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
	GetEvent(req, event_type)
}

func GetEvent(req *http.Request, event_type string) (Event, error) {
	switch event_type {
	case "deployment":
		return decodeDeploymentEvent(req), nil
	case "deployment_status":
		return decodeDeploymentStatusEvent(req), nil
	}

	return nil, errors.New("Error: no matched event type")
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
