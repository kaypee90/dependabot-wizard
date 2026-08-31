package main

import (
	"os"
	"testing"
	"time"

	"github.com/kaypee90/dependabot-wizard/github"
)

// TestPrintSampleReport is a documentation helper, not an assertion: run it
// with -run TestPrintSampleReport -v to regenerate the example in the README.
func TestPrintSampleReport(t *testing.T) {
	if os.Getenv("PRINT_SAMPLE") == "" {
		t.Skip("set PRINT_SAMPLE=1 to print the sample report")
	}

	now := time.Date(2026, 8, 31, 12, 0, 0, 0, time.UTC)
	day := func(n int) string { return now.AddDate(0, 0, -n).Format(time.RFC3339) }

	pullRequests := []github.PullRequest{
		{Number: 61, State: "OPEN", Url: "https://github.com/kaypee90/zebrok/pull/61", CreatedAt: day(34)},
		{Number: 63, State: "OPEN", Url: "https://github.com/kaypee90/zebrok/pull/63", CreatedAt: day(12)},
		{Number: 64, State: "OPEN", Url: "https://github.com/kaypee90/zebrok/pull/64", CreatedAt: day(3)},
		{Number: 55, State: "MERGED", CreatedAt: day(40)},
		{Number: 56, State: "MERGED", CreatedAt: day(52)},
		{Number: 57, State: "MERGED", CreatedAt: day(60)},
		{Number: 58, State: "MERGED", CreatedAt: day(71)},
		{Number: 59, State: "CLOSED", CreatedAt: day(80)},
	}

	os.Stdout.WriteString("\n" + BuildHealthReport("kaypee90", "zebrok", github.DependabotFilter(), pullRequests, now).Markdown() + "\n")
}
