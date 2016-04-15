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
	payload := decodePayload(req)

	fmt.Printf("request: %v", payload.Repo)
}

func decodePayload(req *http.Request) github.WebHookPayload {
	decoder := json.NewDecoder(req.Body)

	var payload github.WebHookPayload

	err := decoder.Decode(&payload)

	if err != nil {
		panic(err)
	}

	return payload
}
