package main

import (
	"embed"
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/kaypee90/dependabot-wizard/github"
)

//go:embed static
var staticFiles embed.FS

type ConfigurationRequest struct {
	Configuration string `json:"configuration"`
}

type ConfigurationResponse struct {
	Path string `json:"path"`
}

type GetPullRequestsRequest struct {
	Owner string `json:"owner"`
	Repo  string `json:"repo"`
	Token string `json:"token"`
}

type DetectResponse struct {
	Ecosystems []DetectedEcosystem `json:"ecosystems"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("writing response: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, ErrorResponse{Error: message})
}

// requireMethod rejects anything but the expected verb. Without this a
// plain GET could trigger a filesystem write.
func requireMethod(w http.ResponseWriter, r *http.Request, method string) bool {
	if r.Method != method {
		w.Header().Set("Allow", method)
		writeError(w, http.StatusMethodNotAllowed, "This endpoint only accepts "+method+" requests.")
		return false
	}

	return true
}

func decodeJSONBody(w http.ResponseWriter, r *http.Request, target any) bool {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(target); err != nil {
		if errors.Is(err, io.EOF) {
			writeError(w, http.StatusBadRequest, "The request body was empty.")
			return false
		}
		writeError(w, http.StatusBadRequest, "The request body could not be read as JSON.")
		return false
	}

	return true
}

func getPullRequestsHandler(w http.ResponseWriter, r *http.Request) {
	if !requireMethod(w, r, http.MethodPost) {
		return
	}

	var request GetPullRequestsRequest
	if !decodeJSONBody(w, r, &request) {
		return
	}

	if request.Owner == "" || request.Repo == "" || request.Token == "" {
		writeError(w, http.StatusBadRequest, "Owner, repo and token are all required.")
		return
	}

	pullRequests, err := github.GetPullRequests(r.Context(), request.Owner, request.Repo, github.DefaultLabel, request.Token)
	if err != nil {
		log.Printf("fetching pull requests for %s/%s: %v", request.Owner, request.Repo, err)
		writeError(w, http.StatusBadGateway, "Could not read pull requests from GitHub: "+err.Error())
		return
	}

	writeJSON(w, http.StatusOK, pullRequests)
}

func saveConfigurationHandler(w http.ResponseWriter, r *http.Request) {
	if !requireMethod(w, r, http.MethodPost) {
		return
	}

	var request ConfigurationRequest
	if !decodeJSONBody(w, r, &request) {
		return
	}

	// Validate before writing. A config that will not parse never reaches disk.
	if _, err := ParseConfig([]byte(request.Configuration)); err != nil {
		writeError(w, http.StatusUnprocessableEntity, "That configuration is not valid: "+err.Error())
		return
	}

	path, err := createDependabotYamlFile([]byte(request.Configuration))
	if err != nil {
		log.Printf("saving configuration: %v", err)
		writeError(w, http.StatusInternalServerError, "Could not write the configuration file: "+err.Error())
		return
	}

	log.Printf("wrote %s", path)
	writeJSON(w, http.StatusCreated, ConfigurationResponse{Path: path})
}

func detectHandler(w http.ResponseWriter, r *http.Request) {
	if !requireMethod(w, r, http.MethodGet) {
		return
	}

	detected, err := DetectEcosystemsInWorkingDirectory()
	if err != nil {
		log.Printf("detecting ecosystems: %v", err)
		writeError(w, http.StatusInternalServerError, "Could not scan this project: "+err.Error())
		return
	}

	writeJSON(w, http.StatusOK, DetectResponse{Ecosystems: detected})
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"message": "Healthy"})
}

func newRouter() (http.Handler, error) {
	// The embedded tree is rooted at "static/"; serving its subtree lets the
	// app live at "/" and keeps the built asset paths correct.
	assets, err := fs.Sub(staticFiles, "static")
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.FS(assets)))
	mux.HandleFunc("/healthz", healthCheckHandler)
	mux.HandleFunc("/api/configurations", saveConfigurationHandler)
	mux.HandleFunc("/api/pull-requests", getPullRequestsHandler)
	mux.HandleFunc("/api/detect", detectHandler)

	return mux, nil
}

// serve runs the web application on an already-bound listener, so callers
// can be certain the port was acquired before advertising the URL.
func serve(listener net.Listener) error {
	handler, err := newRouter()
	if err != nil {
		return err
	}

	server := &http.Server{
		Handler:           handler,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	return server.Serve(listener)
}
