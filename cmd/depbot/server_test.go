package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kaypee90/dependabot-wizard/github"
	"gotest.tools/v3/assert"
)

func newTestServer(t *testing.T) http.Handler {
	t.Helper()

	handler, err := newRouter()
	assert.NilError(t, err)

	return handler
}

func do(t *testing.T, handler http.Handler, method string, target string, body string) *httptest.ResponseRecorder {
	t.Helper()

	var reader *strings.Reader
	if body == "" {
		reader = strings.NewReader("")
	} else {
		reader = strings.NewReader(body)
	}

	request := httptest.NewRequest(method, target, reader)
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)

	return recorder
}

func TestHealthCheckReturnsOK(t *testing.T) {
	response := do(t, newTestServer(t), http.MethodGet, "/healthz", "")

	assert.Equal(t, response.Code, http.StatusOK)
	assert.Assert(t, strings.Contains(response.Body.String(), "Healthy"))
}

func TestSaveConfigurationWritesTheFileAndReportsThePath(t *testing.T) {
	root := t.TempDir()
	chdir(t, root)

	payload := `{"configuration":"version: 2\nupdates:\n- package-ecosystem: npm\n  directory: /\n  schedule:\n    interval: weekly\n"}`
	response := do(t, newTestServer(t), http.MethodPost, "/api/configurations", payload)

	assert.Equal(t, response.Code, http.StatusCreated)

	var decoded ConfigurationResponse
	assert.NilError(t, json.Unmarshal(response.Body.Bytes(), &decoded))
	assert.Equal(t, decoded.Path, filepath.Join(githubDirectory, dependabotFileName))

	written, err := os.ReadFile(filepath.Join(root, githubDirectory, dependabotFileName))
	assert.NilError(t, err)
	assert.Assert(t, strings.HasPrefix(string(written), "version: 2"))
}

// TestSaveConfigurationRejectsUnparseableYaml is the server-side half of the
// guard against the bug that shipped: invalid config never reaches disk.
func TestSaveConfigurationRejectsUnparseableYaml(t *testing.T) {
	root := t.TempDir()
	chdir(t, root)

	payload := `{"configuration":"version: 2\n  updates:\n  - package-ecosystem: npm\n"}`
	response := do(t, newTestServer(t), http.MethodPost, "/api/configurations", payload)

	assert.Equal(t, response.Code, http.StatusUnprocessableEntity)

	_, err := os.Stat(filepath.Join(root, githubDirectory, dependabotFileName))
	assert.Assert(t, os.IsNotExist(err), "nothing should have been written")
}

func TestSaveConfigurationRejectsInvalidEcosystem(t *testing.T) {
	chdir(t, t.TempDir())

	payload := `{"configuration":"version: 2\nupdates:\n- package-ecosystem: cocoapods\n  directory: /\n  schedule:\n    interval: weekly\n"}`
	response := do(t, newTestServer(t), http.MethodPost, "/api/configurations", payload)

	assert.Equal(t, response.Code, http.StatusUnprocessableEntity)
	assert.Assert(t, strings.Contains(response.Body.String(), "unsupported package-ecosystem"))
}

func TestSaveConfigurationRejectsNonPostRequests(t *testing.T) {
	root := t.TempDir()
	chdir(t, root)

	response := do(t, newTestServer(t), http.MethodGet, "/api/configurations", "")

	assert.Equal(t, response.Code, http.StatusMethodNotAllowed)

	_, err := os.Stat(filepath.Join(root, githubDirectory, dependabotFileName))
	assert.Assert(t, os.IsNotExist(err), "a GET must never write a file")
}

func TestSaveConfigurationRejectsAnEmptyBody(t *testing.T) {
	chdir(t, t.TempDir())

	response := do(t, newTestServer(t), http.MethodPost, "/api/configurations", "")

	assert.Equal(t, response.Code, http.StatusBadRequest)
}

func TestPullRequestsRequiresCredentials(t *testing.T) {
	response := do(t, newTestServer(t), http.MethodPost, "/api/pull-requests", `{"owner":"","repo":"","token":""}`)

	assert.Equal(t, response.Code, http.StatusBadRequest)
}

// TestPullRequestsFilterDefaultsToDependabotRegardlessOfLabel pins the fix for
// the reported bug: no label is required for dependabot's pull requests to be
// found, because the filter is on the author.
func TestPullRequestsFilterDefaultsToDependabot(t *testing.T) {
	request := GetPullRequestsRequest{Owner: "o", Repo: "r", Token: "t"}

	filter := request.filter()

	assert.Equal(t, filter.Author, github.DefaultAuthor)
	assert.Equal(t, filter.Label, "")
	assert.Equal(t, github.BuildSearchQuery("o", "r", filter),
		"repo:o/r is:pr sort:created-desc author:app/dependabot")
}

func TestPullRequestsFilterAcceptsAnOptionalLabel(t *testing.T) {
	request := GetPullRequestsRequest{Owner: "o", Repo: "r", Token: "t", Label: "go mod"}

	assert.Equal(t, request.filter().Label, "go mod")
}

func TestPullRequestsFilterAllowsClearingTheAuthor(t *testing.T) {
	empty := ""
	request := GetPullRequestsRequest{Owner: "o", Repo: "r", Token: "t", Author: &empty}

	assert.Equal(t, request.filter().Author, "")
}

func TestPullRequestsRejectsNonPostRequests(t *testing.T) {
	response := do(t, newTestServer(t), http.MethodGet, "/api/pull-requests", "")

	assert.Equal(t, response.Code, http.StatusMethodNotAllowed)
}

func TestDetectReportsEcosystemsForTheWorkingDirectory(t *testing.T) {
	root := t.TempDir()
	assert.NilError(t, os.WriteFile(filepath.Join(root, "go.mod"), []byte("module x\n"), 0o644))
	chdir(t, root)

	response := do(t, newTestServer(t), http.MethodGet, "/api/detect", "")
	assert.Equal(t, response.Code, http.StatusOK)

	var decoded DetectResponse
	assert.NilError(t, json.Unmarshal(response.Body.Bytes(), &decoded))
	assert.Equal(t, len(decoded.Ecosystems), 1)
	assert.Equal(t, decoded.Ecosystems[0].Ecosystem, "gomod")
}

func TestSinglePageApplicationIsServedAtRoot(t *testing.T) {
	response := do(t, newTestServer(t), http.MethodGet, "/", "")

	assert.Equal(t, response.Code, http.StatusOK)
	assert.Assert(t, strings.Contains(strings.ToLower(response.Body.String()), "<!doctype html"))
}

// TestNoWildcardCORSHeader guards the hole where any site the user visited
// could POST to the local server and overwrite a file on disk.
func TestNoWildcardCORSHeader(t *testing.T) {
	response := do(t, newTestServer(t), http.MethodGet, "/healthz", "")

	assert.Equal(t, response.Header().Get("Access-Control-Allow-Origin"), "")
}
