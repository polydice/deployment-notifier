package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/google/go-github/github"
	"github.com/julienschmidt/httprouter"
	"github.com/zorkian/go-datadog-api"
)

type GithubEvent struct {
	Repo             *github.Repository
	EventType        string
	DeploymentStatus *github.DeploymentStatus
}

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
	event, err := GetEvent(req, event_type)

	fmt.Printf("Received %v event of %v", event_type, *event.Repo.Name)

	if err != nil {
		fmt.Print(err)
	}

	datadog_client := GetDataDogClient()
	datadog_event := &datadog.Event{
		Title: "Received " + event_type + " event of " + *event.Repo.Name,
		Text:  "Test",
	}

	_, datadog_err := datadog_client.PostEvent(datadog_event)
	if err != nil {
		fmt.Print(datadog_err)
	}
}

func GetEvent(req *http.Request, event_type string) (GithubEvent, error) {
	switch event_type {
	case "deployment":
		return decodeDeploymentEvent(req), nil
	case "deployment_status":
		return decodeDeploymentStatusEvent(req), nil
	}

	return GithubEvent{}, errors.New("Error: no matched event type")
}

func GetDataDogClient() *datadog.Client {
	api_key := os.Getenv("DATADOG_API_KEY")
	app_key := os.Getenv("DATADOG_APP_KEY")

	client := datadog.NewClient(api_key, app_key)
	return client
}

func decodeDeploymentStatusEvent(req *http.Request) GithubEvent {
	decoder := json.NewDecoder(req.Body)

	var event github.DeploymentStatusEvent

	err := decoder.Decode(&event)

	if err != nil {
		fmt.Print(err)
	}

	github_event := GithubEvent{
		Repo:             event.Repo,
		EventType:        "DeploymentStatus",
		DeploymentStatus: event.DeploymentStatus,
	}

	return github_event
}

func decodeDeploymentEvent(req *http.Request) GithubEvent {
	decoder := json.NewDecoder(req.Body)

	var event github.DeploymentEvent

	err := decoder.Decode(&event)

	if err != nil {
		fmt.Print(err)
	}

	github_event := GithubEvent{
		Repo:      event.Repo,
		EventType: "DeploymentStatus",
	}

	return github_event
}
