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
	Deployment       *github.Deployment
}

func main() {
	route := httprouter.New()
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	checkDatadogKeys()

	route.POST("/webhook", DeploymentHandler)

	fmt.Println("Starting server on :" + port)

	http.ListenAndServe(":"+port, route)
}

func DeploymentHandler(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
	event_type := req.Header.Get("X-Github-Event")
	event, err := NewGithubEvent(req, event_type)

	if err != nil {
		fmt.Print(err)
		return
	}

	fmt.Printf("Recieved %v event from %v\n", event_type, *event.Repo.FullName)

	datadog_client := NewDataDogClient()
	datadog_event := NewDatadogEvent(&event)

	_, datadog_err := datadog_client.PostEvent(datadog_event)
	if err != nil {
		fmt.Print(datadog_err)
	}
}

func NewGithubEvent(req *http.Request, event_type string) (GithubEvent, error) {
	switch event_type {
	case "deployment":
		return decodeDeploymentEvent(req), nil
	case "deployment_status":
		return decodeDeploymentStatusEvent(req), nil
	}

	return GithubEvent{}, errors.New("Not Deployment or DeploymentStatus event\n")
}

func NewDatadogEvent(event *GithubEvent) *datadog.Event {
	repoName := *event.Repo.FullName + ":" + *event.Deployment.SHA
	status := event.DeploymentStatus
	switch status {
	case nil:
		return &datadog.Event{
			Title: "Deployment of " + repoName + " started.",
		}
	default:
		return &datadog.Event{
			Title: "Deployment of " + repoName + " is " + *status.State,
			Text:  "%%% \n Status: " + "[" + *status.State + "](" + *status.TargetURL + ") \n %%%",
		}
	}
}

func checkDatadogKeys() {
	api_key := os.Getenv("DATADOG_API_KEY")
	app_key := os.Getenv("DATADOG_APP_KEY")

	if api_key == "" {
		panic("Api key can't be empty!")
	}

	if app_key == "" {
		panic("Application key can't be empty!")
	}
}

func NewDataDogClient() *datadog.Client {
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
		Deployment:       event.Deployment,
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
		Repo:       event.Repo,
		EventType:  "Deployment",
		Deployment: event.Deployment,
	}

	return github_event
}
