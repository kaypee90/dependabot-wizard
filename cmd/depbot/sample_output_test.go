package main

import (
	"os"
	"testing"
	"time"

	"github.com/kaypee90/dependabot-wizard/github"
)

// TestPrintSampleReport is a documentation helper, not an assertion: run it
// with PRINT_SAMPLE=1 -run TestPrintSampleReport -v to regenerate the example
// in the README, so the documented output cannot drift from the real one.
func TestPrintSampleReport(t *testing.T) {
	if os.Getenv("PRINT_SAMPLE") == "" {
		t.Skip("set PRINT_SAMPLE=1 to print the sample report")
	}

	now := time.Date(2026, 8, 31, 12, 0, 0, 0, time.UTC)
	day := func(n int) string { return now.AddDate(0, 0, -n).Format(time.RFC3339) }

	open := func(number, age int, ecosystem, title string) github.PullRequest {
		return github.PullRequest{
			Number: number, State: "OPEN", Title: title,
			Url:         "https://github.com/kaypee90/zebrok/pull/" + itoa(number),
			CreatedAt:   day(age),
			HeadRefName: "dependabot/" + ecosystem + "/dep",
		}
	}
	merged := func(number, age, took int, ecosystem string) github.PullRequest {
		return github.PullRequest{
			Number: number, State: "MERGED", Title: "Bump a dependency",
			CreatedAt: day(age), MergedAt: day(age - took),
			HeadRefName: "dependabot/" + ecosystem + "/dep",
		}
	}

	pullRequests := []github.PullRequest{
		open(61, 34, "npm_and_yarn", "Bump webpack from 5.88.0 to 5.94.0"),
		open(63, 12, "npm_and_yarn", "Bump eslint from 8.50.0 to 9.9.0"),
		open(64, 3, "go_modules", "Bump golang.org/x/net from 0.17.0 to 0.28.0"),
		merged(55, 20, 2, "npm_and_yarn"),
		merged(56, 26, 1, "npm_and_yarn"),
		merged(57, 40, 3, "go_modules"),
		merged(58, 55, 1, "github_actions"),
		{Number: 59, State: "CLOSED", Title: "Bump a major", CreatedAt: day(60),
			ClosedAt: day(59), HeadRefName: "dependabot/npm_and_yarn/dep"},
	}

	os.Stdout.WriteString("\n" + BuildHealthReport("kaypee90", "zebrok",
		github.DependabotFilter(), pullRequests, now).Markdown() + "\n")
}

func itoa(value int) string {
	if value == 0 {
		return "0"
	}
	digits := ""
	for value > 0 {
		digits = string(rune('0'+value%10)) + digits
		value /= 10
	}
	return digits
}
