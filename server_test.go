package main

import (
	"net/http"
	"strings"
	"testing"

	"github.com/google/go-github/github"
	"github.com/stretchr/testify/assert"
)

func TestGetEventReturnsDeploymentEvent(t *testing.T) {
	assert := assert.New(t)

	request_body := strings.NewReader("Fake request body")
	req, _ := http.NewRequest("POST", "/webhook", request_body)

	event, _ := GetEvent(req, "deployment")

	var deploy_event github.DeploymentEvent
	assert.IsType(deploy_event, event)
}

func TestGetEventReturnsDeploymentStatusEvent(t *testing.T) {
	assert := assert.New(t)

	request_body := strings.NewReader("Fake request body")
	req, _ := http.NewRequest("POST", "/webhook", request_body)

	event, _ := GetEvent(req, "deployment_status")

	var deploy_event github.DeploymentStatusEvent
	assert.IsType(deploy_event, event)
}
