package main

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEventReturnsDeploymentEvent(t *testing.T) {
	assert := assert.New(t)

	request_body := strings.NewReader("Fake request body")
	req, _ := http.NewRequest("POST", "/webhook", request_body)

	event, _ := NewGithubEvent(req, "deployment")

	assert.Equal(event.EventType, "Deployment")
}

func TestGetEventReturnsDeploymentStatusEvent(t *testing.T) {
	assert := assert.New(t)

	request_body := strings.NewReader("Fake request body")
	req, _ := http.NewRequest("POST", "/webhook", request_body)

	event, _ := NewGithubEvent(req, "deployment_status")

	assert.Equal(event.EventType, "DeploymentStatus")
}
