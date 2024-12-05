package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/kaypee90/dependabot-wizard/github"
)

//go:embed static/*
var staticFiles embed.FS

type ConfigurationRequest struct {
	Configuration string `json:"configuration"`
}

type GetPullRequestsRequest struct {
	Owner string `json:"owner"`
	Repo  string `json:"repo"`
	Token string `json:"token"`
}

func getPullRequestsHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var request GetPullRequestsRequest
	if err := decoder.Decode(&request); err != nil {
		if err == io.EOF {
			http.Error(w, "Empty request body", http.StatusBadRequest)
			return
		}
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	pullRequests := github.GetPullRequests(request.Owner, request.Repo, request.Token)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pullRequests)
}

func saveConfigurationHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var config ConfigurationRequest
	if err := decoder.Decode(&config); err != nil {
		if err == io.EOF {
			http.Error(w, "Empty request body", http.StatusBadRequest)
			return
		}
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	createDependabotYamlFile([]byte(config.Configuration))
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{\"message\": \"Healthy\"}"))
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set the necessary CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight request
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func startWebApplication() {
	mux := http.NewServeMux()

	fs := http.FileServer(http.FS(staticFiles))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("/", healthCheckHandler)
	mux.HandleFunc("/api/configurations", saveConfigurationHandler)
	mux.HandleFunc("/api/pull-requests", getPullRequestsHandler)

	port := ":3001"
	fmt.Println("Serving on http://localhost" + port)
	http.ListenAndServe(port, enableCORS(mux))
}
