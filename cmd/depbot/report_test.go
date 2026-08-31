package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/kaypee90/dependabot-wizard/github"
	"gotest.tools/v3/assert"
)

var reportNow = time.Date(2026, 8, 31, 12, 0, 0, 0, time.UTC)

func pr(number int, state string, daysAgo int) github.PullRequest {
	return github.PullRequest{
		Number:    number,
		State:     state,
		Url:       "https://example.test/" + string(rune('0'+number%10)),
		Title:     "Bump something",
		CreatedAt: reportNow.AddDate(0, 0, -daysAgo).Format(time.RFC3339),
	}
}

func TestBuildHealthReportCountsEachState(t *testing.T) {
	report := BuildHealthReport("kaypee90", "zebrok", "dependencies", []github.PullRequest{
		pr(1, "OPEN", 3), pr(2, "MERGED", 10), pr(3, "MERGED", 12), pr(4, "CLOSED", 20),
	}, reportNow)

	assert.Equal(t, report.Total, 4)
	assert.Equal(t, report.Open, 1)
	assert.Equal(t, report.Merged, 2)
	assert.Equal(t, report.Closed, 1)
}

func TestBuildHealthReportExcludesOpenFromMergeRate(t *testing.T) {
	// 2 merged of 4 resolved is 50%, regardless of how many are still open.
	report := BuildHealthReport("o", "r", "dependencies", []github.PullRequest{
		pr(1, "OPEN", 1), pr(2, "OPEN", 2),
		pr(3, "MERGED", 5), pr(4, "MERGED", 6), pr(5, "CLOSED", 7), pr(6, "CLOSED", 8),
	}, reportNow)

	assert.Equal(t, report.MergeRate, 50.0)
}

func TestBuildHealthReportFindsTheOldestOpenPullRequest(t *testing.T) {
	report := BuildHealthReport("o", "r", "dependencies", []github.PullRequest{
		pr(1, "OPEN", 3), pr(2, "OPEN", 45), pr(3, "OPEN", 9), pr(4, "MERGED", 100),
	}, reportNow)

	assert.Assert(t, report.OldestOpen != nil)
	assert.Equal(t, report.OldestOpen.Number, 2)
	assert.Equal(t, int(report.OldestOpenAge.Hours()/24), 45)
}

func TestBuildHealthReportHandlesNoPullRequests(t *testing.T) {
	report := BuildHealthReport("o", "r", "dependencies", nil, reportNow)

	assert.Equal(t, report.Total, 0)
	assert.Equal(t, report.MergeRate, 0.0)
	assert.Assert(t, report.OldestOpen == nil)
	assert.Assert(t, strings.Contains(report.Markdown(), "No pull requests labelled"))
}

func TestBuildHealthReportIgnoresUnparseableTimestamps(t *testing.T) {
	broken := github.PullRequest{Number: 9, State: "OPEN", CreatedAt: "not-a-date"}

	report := BuildHealthReport("o", "r", "dependencies", []github.PullRequest{broken}, reportNow)

	assert.Equal(t, report.Open, 1)
	assert.Assert(t, report.OldestOpen == nil, "an unparseable date must not become the oldest")
}

func TestMarkdownReportsCountsAndAVerdict(t *testing.T) {
	report := BuildHealthReport("kaypee90", "zebrok", "dependencies", []github.PullRequest{
		pr(1, "OPEN", 3), pr(2, "MERGED", 10), pr(3, "MERGED", 12), pr(4, "CLOSED", 20),
	}, reportNow)

	markdown := report.Markdown()

	assert.Assert(t, strings.Contains(markdown, "## Dependabot health: kaypee90/zebrok"))
	assert.Assert(t, strings.Contains(markdown, "| Total pull requests | 4 |"))
	assert.Assert(t, strings.Contains(markdown, "| Merge rate | 67% |"))
	assert.Assert(t, strings.Contains(markdown, "keeping up"))
}

func TestMarkdownRecommendsGroupingForALargeQueue(t *testing.T) {
	var pullRequests []github.PullRequest
	for i := 1; i <= 15; i++ {
		pullRequests = append(pullRequests, pr(i, "OPEN", i))
	}
	pullRequests = append(pullRequests, pr(90, "MERGED", 30))

	markdown := BuildHealthReport("o", "r", "dependencies", pullRequests, reportNow).Markdown()

	assert.Assert(t, strings.Contains(markdown, "groups"))
}

func TestMarkdownFlagsALowMergeRate(t *testing.T) {
	markdown := BuildHealthReport("o", "r", "dependencies", []github.PullRequest{
		pr(1, "OPEN", 1), pr(2, "CLOSED", 5), pr(3, "CLOSED", 6), pr(4, "MERGED", 7),
	}, reportNow).Markdown()

	assert.Assert(t, strings.Contains(markdown, "ignore"))
}

func TestPullRequestsFromResponseFlattensAndSortsNewestFirst(t *testing.T) {
	response := github.PullRequestResponse{}
	response.Data.Repository.PullRequests.Edges = []github.PullRequestNode{
		{Node: pr(1, "OPEN", 30)},
		{Node: pr(2, "OPEN", 1)},
		{Node: pr(3, "OPEN", 10)},
	}

	flattened := pullRequestsFromResponse(response)

	assert.Equal(t, len(flattened), 3)
	assert.Equal(t, flattened[0].Number, 2)
	assert.Equal(t, flattened[2].Number, 1)
}

// TestAppendToFileFromEnv covers how the Action receives the summary.
func TestAppendToFileFromEnv(t *testing.T) {
	path := filepath.Join(t.TempDir(), "summary.md")
	t.Setenv("GITHUB_STEP_SUMMARY", path)

	assert.NilError(t, appendToFileFromEnv("GITHUB_STEP_SUMMARY", "first\n"))
	assert.NilError(t, appendToFileFromEnv("GITHUB_STEP_SUMMARY", "second\n"))

	written, err := os.ReadFile(path)
	assert.NilError(t, err)
	assert.Equal(t, string(written), "first\nsecond\n")
}

func TestAppendToFileFromEnvIsANoOpOutsideActions(t *testing.T) {
	t.Setenv("GITHUB_STEP_SUMMARY", "")

	assert.NilError(t, appendToFileFromEnv("GITHUB_STEP_SUMMARY", "ignored"))
}

func TestRunReportRequiresAToken(t *testing.T) {
	t.Setenv("GITHUB_TOKEN", "")

	err := runReport("kaypee90", "zebrok", "dependencies")

	assert.ErrorContains(t, err, "GITHUB_TOKEN is not set")
}

func TestRunReportRequiresARepository(t *testing.T) {
	t.Setenv("GITHUB_TOKEN", "token")
	t.Setenv("GITHUB_REPOSITORY", "")

	err := runReport("", "", "dependencies")

	assert.ErrorContains(t, err, "--owner and --repo are required")
}
