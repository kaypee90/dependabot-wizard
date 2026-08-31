package github

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"gotest.tools/v3/assert"
)

// stubGitHub points the client at a local server for the duration of a test.
func stubGitHub(t *testing.T, handler http.HandlerFunc) {
	t.Helper()

	server := httptest.NewServer(handler)
	previous := endpoint
	endpoint = server.URL

	t.Cleanup(func() {
		endpoint = previous
		server.Close()
	})
}

func TestGetPullRequestsReturnsPullRequests(t *testing.T) {
	stubGitHub(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Header.Get("Authorization"), "Bearer test-token")

		var received PullRequestQuery
		assert.NilError(t, json.NewDecoder(r.Body).Decode(&received))
		// Owner and repo must travel as variables, not be interpolated.
		assert.Equal(t, received.Variables["owner"], "kaypee90")
		assert.Equal(t, received.Variables["repo"], "zebrok")
		assert.Equal(t, received.Variables["label"], "dependencies")

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"data":{"repository":{"pullRequests":{"edges":[
			{"node":{"title":"Bump left-pad","number":7,"url":"https://example.test/7","state":"MERGED","createdAt":"2026-01-04T10:00:00Z","author":{"login":"dependabot"}}}
		]}}}}`))
	})

	response, err := GetPullRequests(context.Background(), "kaypee90", "zebrok", DefaultLabel, "test-token")

	assert.NilError(t, err)
	edges := response.Data.Repository.PullRequests.Edges
	assert.Equal(t, len(edges), 1)
	assert.Equal(t, edges[0].Node.Number, 7)
	assert.Equal(t, edges[0].Node.State, "MERGED")
}

// TestGetPullRequestsSurfacesGraphqlErrors covers the case that used to look
// like success: HTTP 200 with an errors block and no data.
func TestGetPullRequestsSurfacesGraphqlErrors(t *testing.T) {
	stubGitHub(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"data":{"repository":null},"errors":[{"message":"Could not resolve to a Repository with the name 'kaypee90/missing'.","type":"NOT_FOUND"}]}`))
	})

	_, err := GetPullRequests(context.Background(), "kaypee90", "missing", DefaultLabel, "test-token")

	assert.ErrorContains(t, err, "Could not resolve to a Repository")
}

func TestGetPullRequestsReportsAnUnauthorizedToken(t *testing.T) {
	stubGitHub(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"message":"Bad credentials"}`))
	})

	_, err := GetPullRequests(context.Background(), "kaypee90", "zebrok", DefaultLabel, "expired-token")

	assert.ErrorContains(t, err, "401")
}

func TestGetPullRequestsReportsMalformedJSON(t *testing.T) {
	stubGitHub(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"data":`))
	})

	_, err := GetPullRequests(context.Background(), "kaypee90", "zebrok", DefaultLabel, "test-token")

	assert.ErrorContains(t, err, "decoding response")
}

func TestGetPullRequestsDoesNotPanicOnTransportFailure(t *testing.T) {
	previous := endpoint
	endpoint = "http://127.0.0.1:1"
	t.Cleanup(func() { endpoint = previous })

	_, err := GetPullRequests(context.Background(), "kaypee90", "zebrok", DefaultLabel, "test-token")

	assert.ErrorContains(t, err, "calling GitHub")
}
