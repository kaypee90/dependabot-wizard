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

func TestBuildSearchQueryFiltersByAuthorNotLabel(t *testing.T) {
	query := BuildSearchQuery("kaypee90", "zebrok", DependabotFilter())

	assert.Equal(t, query, "repo:kaypee90/zebrok is:pr sort:created-desc author:app/dependabot")
}

func TestBuildSearchQueryAddsALabelOnlyWhenAsked(t *testing.T) {
	query := BuildSearchQuery("o", "r", Filter{Author: DefaultAuthor, Label: "dependencies"})

	assert.Equal(t, query, `repo:o/r is:pr sort:created-desc author:app/dependabot label:"dependencies"`)
}

func TestBuildSearchQueryQuotesLabelsContainingSpaces(t *testing.T) {
	query := BuildSearchQuery("o", "r", Filter{Label: "go mod"})

	assert.Equal(t, query, `repo:o/r is:pr sort:created-desc label:"go mod"`)
}

func TestBuildSearchQueryOmitsAuthorWhenBlank(t *testing.T) {
	query := BuildSearchQuery("o", "r", Filter{})

	assert.Equal(t, query, "repo:o/r is:pr sort:created-desc")
}

func TestDescribeExplainsTheActiveFilter(t *testing.T) {
	assert.Equal(t, DependabotFilter().Describe(), "opened by app/dependabot")
	assert.Equal(t, Filter{Label: "deps"}.Describe(), `labelled "deps"`)
	assert.Equal(t, Filter{Author: "app/renovate", Label: "deps"}.Describe(), `opened by app/renovate and labelled "deps"`)
	assert.Equal(t, Filter{}.Describe(), "in this repository")
}

// TestUnlabelledDependabotPullRequestsAreFound is the regression test for the
// reported bug: dependabot had opened pull requests, none of them carried the
// "dependencies" label, and the dashboard showed nothing.
func TestUnlabelledDependabotPullRequestsAreFound(t *testing.T) {
	stubGitHub(t, func(w http.ResponseWriter, r *http.Request) {
		var received SearchQuery
		assert.NilError(t, json.NewDecoder(r.Body).Decode(&received))

		// The search must not constrain on a label by default.
		assert.Equal(t, received.Variables["search"],
			"repo:kaypee90/zebrok is:pr sort:created-desc author:app/dependabot")

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"data":{"search":{"issueCount":2,"edges":[
			{"node":{"title":"Bump golang.org/x/net","number":12,"url":"https://example.test/12","state":"OPEN","createdAt":"2026-08-01T10:00:00Z","author":{"login":"dependabot"}}},
			{"node":{"title":"Bump react","number":14,"url":"https://example.test/14","state":"OPEN","createdAt":"2026-08-20T10:00:00Z","author":{"login":"dependabot"}}}
		]}}}`))
	})

	pullRequests, err := GetPullRequests(context.Background(), "kaypee90", "zebrok", DependabotFilter(), "test-token")

	assert.NilError(t, err)
	assert.Equal(t, len(pullRequests), 2)
	// Newest first.
	assert.Equal(t, pullRequests[0].Number, 14)
	assert.Equal(t, pullRequests[1].Number, 12)
	assert.Equal(t, pullRequests[0].AuthorLogin(), "dependabot")
}

func TestGetPullRequestsSendsTheToken(t *testing.T) {
	stubGitHub(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Header.Get("Authorization"), "Bearer test-token")
		w.Write([]byte(`{"data":{"search":{"issueCount":0,"edges":[]}}}`))
	})

	pullRequests, err := GetPullRequests(context.Background(), "o", "r", DependabotFilter(), "test-token")

	assert.NilError(t, err)
	assert.Equal(t, len(pullRequests), 0)
}

func TestGetPullRequestsSkipsNonPullRequestResults(t *testing.T) {
	stubGitHub(t, func(w http.ResponseWriter, r *http.Request) {
		// Search over type ISSUE can return an issue, which the inline
		// fragment leaves as an empty node.
		w.Write([]byte(`{"data":{"search":{"issueCount":2,"edges":[
			{"node":{}},
			{"node":{"title":"Bump x","number":3,"state":"MERGED","createdAt":"2026-08-02T10:00:00Z"}}
		]}}}`))
	})

	pullRequests, err := GetPullRequests(context.Background(), "o", "r", DependabotFilter(), "t")

	assert.NilError(t, err)
	assert.Equal(t, len(pullRequests), 1)
	assert.Equal(t, pullRequests[0].Number, 3)
}

func TestGetPullRequestsHandlesAMissingAuthor(t *testing.T) {
	stubGitHub(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"data":{"search":{"issueCount":1,"edges":[
			{"node":{"title":"Bump x","number":5,"state":"OPEN","createdAt":"2026-08-02T10:00:00Z","author":null}}
		]}}}`))
	})

	pullRequests, err := GetPullRequests(context.Background(), "o", "r", DependabotFilter(), "t")

	assert.NilError(t, err)
	assert.Equal(t, pullRequests[0].AuthorLogin(), "")
}

func TestGetPullRequestsSurfacesGraphqlErrors(t *testing.T) {
	stubGitHub(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"data":{"search":null},"errors":[{"message":"Could not resolve to a Repository with the name 'kaypee90/missing'.","type":"NOT_FOUND"}]}`))
	})

	_, err := GetPullRequests(context.Background(), "kaypee90", "missing", DependabotFilter(), "test-token")

	assert.ErrorContains(t, err, "Could not resolve to a Repository")
}

func TestGetPullRequestsReportsAnUnauthorizedToken(t *testing.T) {
	stubGitHub(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"message":"Bad credentials"}`))
	})

	_, err := GetPullRequests(context.Background(), "o", "r", DependabotFilter(), "expired")

	assert.ErrorContains(t, err, "401")
}

func TestGetPullRequestsReportsMalformedJSON(t *testing.T) {
	stubGitHub(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"data":`))
	})

	_, err := GetPullRequests(context.Background(), "o", "r", DependabotFilter(), "t")

	assert.ErrorContains(t, err, "decoding response")
}

func TestGetPullRequestsDoesNotPanicOnTransportFailure(t *testing.T) {
	previous := endpoint
	endpoint = "http://127.0.0.1:1"
	t.Cleanup(func() { endpoint = previous })

	_, err := GetPullRequests(context.Background(), "o", "r", DependabotFilter(), "t")

	assert.ErrorContains(t, err, "calling GitHub")
}
