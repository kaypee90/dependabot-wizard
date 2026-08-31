package main

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/kaypee90/dependabot-wizard/github"
)

// HealthReport summarises how a repository is handling its dependency pull
// requests. It answers the question the config leads to: is dependabot
// actually working here, or is nobody merging what it opens?
type HealthReport struct {
	Owner  string
	Repo   string
	Label  string
	Total  int
	Open   int
	Merged int
	Closed int

	// MergeRate is merged as a percentage of resolved (merged plus closed)
	// pull requests. Open ones are excluded because they are undecided.
	MergeRate float64

	// OldestOpen is the longest-waiting open pull request, if there is one.
	OldestOpen    *github.PullRequest
	OldestOpenAge time.Duration
}

// BuildHealthReport derives the summary from raw pull requests. It is a pure
// function so the reporting logic is testable without touching the network.
func BuildHealthReport(owner string, repo string, label string, pullRequests []github.PullRequest, now time.Time) HealthReport {
	report := HealthReport{Owner: owner, Repo: repo, Label: label, Total: len(pullRequests)}

	for index := range pullRequests {
		pullRequest := pullRequests[index]

		switch strings.ToUpper(pullRequest.State) {
		case "OPEN":
			report.Open++
			created, err := time.Parse(time.RFC3339, pullRequest.CreatedAt)
			if err != nil {
				continue
			}
			age := now.Sub(created)
			if report.OldestOpen == nil || age > report.OldestOpenAge {
				report.OldestOpen = &pullRequests[index]
				report.OldestOpenAge = age
			}
		case "MERGED":
			report.Merged++
		case "CLOSED":
			report.Closed++
		}
	}

	if resolved := report.Merged + report.Closed; resolved > 0 {
		report.MergeRate = float64(report.Merged) / float64(resolved) * 100
	}

	return report
}

// Markdown renders the report for a GitHub Actions step summary or a comment.
func (report HealthReport) Markdown() string {
	var builder strings.Builder

	fmt.Fprintf(&builder, "## Dependabot health: %s/%s\n\n", report.Owner, report.Repo)

	if report.Total == 0 {
		fmt.Fprintf(&builder, "No pull requests labelled `%s` were found.\n\n", report.Label)
		fmt.Fprintf(&builder, "That usually means dependabot is not configured, or it is configured "+
			"without that label. Run `depbot` in this repository to generate a configuration.\n")
		return builder.String()
	}

	fmt.Fprintf(&builder, "| Metric | Value |\n| --- | ---: |\n")
	fmt.Fprintf(&builder, "| Total pull requests | %d |\n", report.Total)
	fmt.Fprintf(&builder, "| Open | %d |\n", report.Open)
	fmt.Fprintf(&builder, "| Merged | %d |\n", report.Merged)
	fmt.Fprintf(&builder, "| Closed without merging | %d |\n", report.Closed)

	if report.Merged+report.Closed > 0 {
		fmt.Fprintf(&builder, "| Merge rate | %.0f%% |\n", report.MergeRate)
	}

	if report.OldestOpen != nil {
		fmt.Fprintf(&builder, "| Oldest open | [#%d](%s), %d days |\n",
			report.OldestOpen.Number, report.OldestOpen.Url, int(report.OldestOpenAge.Hours()/24))
	}

	builder.WriteString("\n")
	builder.WriteString(report.verdict())

	return builder.String()
}

// verdict turns the numbers into the one line a reader actually acts on.
func (report HealthReport) verdict() string {
	switch {
	case report.Open == 0:
		return "Nothing is waiting. Dependabot updates are being kept up with.\n"
	case report.Merged+report.Closed == 0:
		return fmt.Sprintf("%d open, none resolved yet. Nobody has reviewed a dependabot pull request in this repository.\n", report.Open)
	case report.MergeRate < 50:
		return fmt.Sprintf("%d open, and only %.0f%% of resolved pull requests were merged. Most updates are being closed rather than taken — consider narrowing the config with `ignore` rules so the ones that remain are worth reviewing.\n", report.Open, report.MergeRate)
	case report.Open > 10:
		return fmt.Sprintf("%d open pull requests is a large queue. Grouping updates with `groups` in dependabot.yml would reduce the count without reducing coverage.\n", report.Open)
	default:
		return fmt.Sprintf("%d open, %.0f%% of resolved pull requests merged. This repository is keeping up.\n", report.Open, report.MergeRate)
	}
}

// pullRequestsFromResponse flattens the GraphQL edge/node nesting.
func pullRequestsFromResponse(response github.PullRequestResponse) []github.PullRequest {
	edges := response.Data.Repository.PullRequests.Edges
	pullRequests := make([]github.PullRequest, 0, len(edges))

	for _, edge := range edges {
		pullRequests = append(pullRequests, edge.Node)
	}

	// Newest first, so any listing reads in a predictable order.
	sort.Slice(pullRequests, func(i, j int) bool {
		return pullRequests[i].CreatedAt > pullRequests[j].CreatedAt
	})

	return pullRequests
}

// appendToFileFromEnv writes to the file named by an environment variable,
// which is how GitHub Actions receives step summaries and outputs. A missing
// variable simply means we are not running inside Actions.
func appendToFileFromEnv(variable string, content string) error {
	path := os.Getenv(variable)
	if path == "" {
		return nil
	}

	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("opening %s: %w", variable, err)
	}
	defer file.Close()

	if _, err := file.WriteString(content); err != nil {
		return fmt.Errorf("writing %s: %w", variable, err)
	}

	return nil
}

// runReport fetches the pull requests and emits the summary. It prints to
// stdout always, and additionally feeds the GitHub Actions step summary and
// outputs when those environment variables are present.
func runReport(owner string, repo string, label string) error {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return fmt.Errorf("GITHUB_TOKEN is not set; the report needs a token with repo scope")
	}

	if owner == "" || repo == "" {
		// Actions provides this as "owner/repo".
		if repository := os.Getenv("GITHUB_REPOSITORY"); repository != "" {
			if parts := strings.SplitN(repository, "/", 2); len(parts) == 2 {
				owner, repo = parts[0], parts[1]
			}
		}
	}

	if owner == "" || repo == "" {
		return fmt.Errorf("both --owner and --repo are required outside GitHub Actions")
	}

	if label == "" {
		label = github.DefaultLabel
	}

	response, err := github.GetPullRequests(context.Background(), owner, repo, label, token)
	if err != nil {
		return err
	}

	report := BuildHealthReport(owner, repo, label, pullRequestsFromResponse(response), time.Now())
	markdown := report.Markdown()

	fmt.Print(markdown)

	if err := appendToFileFromEnv("GITHUB_STEP_SUMMARY", markdown); err != nil {
		return err
	}

	outputs := fmt.Sprintf("total=%d\nopen=%d\nmerged=%d\nclosed=%d\nmerge-rate=%.0f\n",
		report.Total, report.Open, report.Merged, report.Closed, report.MergeRate)

	return appendToFileFromEnv("GITHUB_OUTPUT", outputs)
}
