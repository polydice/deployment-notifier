package main

import (
	"crypto/sha1"
	"net/http"
	"strings"
	"testing"

	"github.com/google/go-github/github"
	"github.com/stretchr/testify/assert"
)

func TestNewDatadogEventReturnsDeploymentEvent(t *testing.T) {
	assert := assert.New(t)
	test_repo_name := "polydice/test"
	test_sha := string(sha1.New().Sum(nil))
	test_event_name := test_repo_name + ":" + test_sha

	event := GithubEvent{
		Repo: &github.Repository{
			FullName: &test_repo_name,
		},
		Deployment: &github.Deployment{
			SHA: &test_sha,
		},
	}
	datadog_event := NewDatadogEvent(&event)
	assert.Contains(datadog_event.Title, test_event_name)
}

func TestNewDatadogEventReturnsDeploymentStatusEvent(t *testing.T) {
	assert := assert.New(t)
	test_repo_name := "polydice/test"
	test_sha := string(sha1.New().Sum(nil))
	test_event_name := test_repo_name + ":" + test_sha
	test_target_url := "https://example.com"
	test_state := "pending"

	event := GithubEvent{
		Repo: &github.Repository{
			FullName: &test_repo_name,
		},
		Deployment: &github.Deployment{
			SHA: &test_sha,
		},
		DeploymentStatus: &github.DeploymentStatus{
			State:     &test_state,
			TargetURL: &test_target_url,
		},
	}
	datadog_event := NewDatadogEvent(&event)
	assert.Contains(datadog_event.Title, test_event_name)
	assert.Contains(datadog_event.Text, test_state)
	assert.Contains(datadog_event.Text, test_target_url)
}

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
