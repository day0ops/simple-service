package test

import (
	"net/http"
	"testing"
)

func TestHealthzHandler(t *testing.T) {
	// Set up the Test Client using the reusable function
	client := NewTestClient()
	// Perform a request to the route
	response := client.PerformRequest("GET", "/healthz", nil, nil)
	// Assert the response
	AssertResponse(t, response, http.StatusOK, "")
}

func TestNonExistentHandler(t *testing.T) {
	// Set up the Test Client using the reusable function
	client := NewTestClient()
	// Perform a request to the route
	response := client.PerformRequest("GET", "/non-existence", nil, nil)
	// Assert the response
	AssertResponseStatus(t, response, http.StatusNotFound)
}

func TestRootHandler(t *testing.T) {
	t.Run("handle when no podname or hostname defined", WrappedTestCase(func(wt *testing.T, tc *WrappedTestClient) {
		// Setup test data
		SetupPodName(wt, "")
		SetupHostname(wt, "")
		// Perform a request to the route
		response := tc.PerformRequest("GET", "/", nil, nil)
		// Assert the response
		AssertResponseStatus(t, response, http.StatusInternalServerError)
	}))

	t.Run("handle when only podname is defined", WrappedTestCase(func(wt *testing.T, tc *WrappedTestClient) {
		expected := "fake-podname"
		// Setup test data
		SetupPodName(wt, expected)
		// Perform a request to the route
		response := tc.PerformRequest("GET", "/", nil, nil)
		// Assert the response
		AssertResponse(wt, response, http.StatusOK, buildExpectedBody(expected))
	}))

	t.Run("handle when only hostname is defined", WrappedTestCase(func(wt *testing.T, tc *WrappedTestClient) {
		expected := "fake-hostmame"
		// Setup test data
		SetupHostname(wt, expected)
		// Perform a request to the route
		response := tc.PerformRequest("GET", "/", nil, nil)
		// Assert the response
		AssertResponse(wt, response, http.StatusOK, buildExpectedBody(expected))
	}))
}

func buildExpectedBody(expected string) string {
	return `{"host": "` + expected + `"}`
}
