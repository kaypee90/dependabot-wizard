package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/kaypee90/dependabot-wizard/github"
	"gotest.tools/v3/assert"
)

var reportNow = time.Date(2026, 8, 31, 12, 0, 0, 0, time.UTC)

func at(daysAgo int) string {
	return reportNow.AddDate(0, 0, -daysAgo).Format(time.RFC3339)
}

// openPR is a pull request still waiting, opened daysAgo.
func openPR(number int, daysAgo int, ecosystem string) github.PullRequest {
	return github.PullRequest{
		Number:      number,
		State:       "OPEN",
		Title:       fmt.Sprintf("Bump dependency %d", number),
		Url:         fmt.Sprintf("https://example.test/%d", number),
		CreatedAt:   at(daysAgo),
		HeadRefName: fmt.Sprintf("dependabot/%s/dep-%d", ecosystem, number),
	}
}

// mergedPR was opened daysAgo and merged afterMergeDays later.
func mergedPR(number int, daysAgo int, afterMergeDays int, ecosystem string) github.PullRequest {
	return github.PullRequest{
		Number:      number,
		State:       "MERGED",
		Title:       fmt.Sprintf("Bump dependency %d", number),
		CreatedAt:   at(daysAgo),
		MergedAt:    at(daysAgo - afterMergeDays),
		HeadRefName: fmt.Sprintf("dependabot/%s/dep-%d", ecosystem, number),
	}
}

func closedPR(number int, daysAgo int, ecosystem string) github.PullRequest {
	return github.PullRequest{
		Number:      number,
		State:       "CLOSED",
		Title:       fmt.Sprintf("Bump dependency %d", number),
		CreatedAt:   at(daysAgo),
		ClosedAt:    at(daysAgo - 1),
		HeadRefName: fmt.Sprintf("dependabot/%s/dep-%d", ecosystem, number),
	}
}

func TestEcosystemForBranchReadsDependabotBranchNames(t *testing.T) {
	assert.Equal(t, EcosystemForBranch("dependabot/npm_and_yarn/lodash-4.17.21"), "npm")
	assert.Equal(t, EcosystemForBranch("dependabot/go_modules/golang.org/x/net-0.17.0"), "gomod")
	assert.Equal(t, EcosystemForBranch("dependabot/github_actions/actions/checkout-4"), "github-actions")
	// A directory in the middle must not confuse the ecosystem segment.
	assert.Equal(t, EcosystemForBranch("dependabot/npm_and_yarn/ui/web/react-18.3.1"), "npm")
	// An unrecognised slug is reported as-is rather than dropped.
	assert.Equal(t, EcosystemForBranch("dependabot/brand_new/thing-1.0"), "brand_new")
	// Not a dependabot branch.
	assert.Equal(t, EcosystemForBranch("feature/my-change"), "")
	assert.Equal(t, EcosystemForBranch(""), "")
}

func TestBuildHealthReportCountsEachState(t *testing.T) {
	report := BuildHealthReport("kaypee90", "zebrok", github.DependabotFilter(), []github.PullRequest{
		openPR(1, 3, "npm_and_yarn"),
		mergedPR(2, 10, 2, "npm_and_yarn"),
		mergedPR(3, 12, 1, "go_modules"),
		closedPR(4, 20, "go_modules"),
	}, reportNow)

	assert.Equal(t, report.Total, 4)
	assert.Equal(t, report.Open, 1)
	assert.Equal(t, report.Merged, 2)
	assert.Equal(t, report.Closed, 1)
}

func TestBuildHealthReportExcludesOpenFromMergeRate(t *testing.T) {
	// 2 merged of 4 resolved is 50%, regardless of how many are still open.
	report := BuildHealthReport("o", "r", github.DependabotFilter(), []github.PullRequest{
		openPR(1, 1, "pip"), openPR(2, 2, "pip"),
		mergedPR(3, 5, 1, "pip"), mergedPR(4, 6, 1, "pip"),
		closedPR(5, 7, "pip"), closedPR(6, 8, "pip"),
	}, reportNow)

	assert.Equal(t, report.MergeRate, 50.0)
}

func TestBuildHealthReportCountsStaleOpenPullRequests(t *testing.T) {
	report := BuildHealthReport("o", "r", github.DependabotFilter(), []github.PullRequest{
		openPR(1, 45, "npm_and_yarn"),
		openPR(2, 31, "npm_and_yarn"),
		openPR(3, 29, "npm_and_yarn"),
		openPR(4, 2, "npm_and_yarn"),
	}, reportNow)

	assert.Equal(t, report.Open, 4)
	assert.Equal(t, report.StaleOpen, 2)
}

func TestBuildHealthReportTakesTheMedianTimeToMerge(t *testing.T) {
	// Merge durations of 1, 3 and 40 days: the median is 3, and the outlier
	// must not drag the headline number the way a mean would.
	report := BuildHealthReport("o", "r", github.DependabotFilter(), []github.PullRequest{
		mergedPR(1, 50, 1, "pip"),
		mergedPR(2, 40, 3, "pip"),
		mergedPR(3, 60, 40, "pip"),
	}, reportNow)

	assert.Assert(t, report.HasMergeTimes)
	assert.Equal(t, days(report.MedianTimeToMerge), 3)
}

func TestBuildHealthReportCountsRecentActivity(t *testing.T) {
	report := BuildHealthReport("o", "r", github.DependabotFilter(), []github.PullRequest{
		openPR(1, 5, "pip"),        // opened recently
		openPR(2, 200, "pip"),      // old
		mergedPR(3, 12, 2, "pip"),  // merged 10 days ago
		mergedPR(4, 400, 1, "pip"), // merged long ago
	}, reportNow)

	assert.Equal(t, report.OpenedRecently, 1)
	assert.Equal(t, report.MergedRecently, 1)
}

func TestBuildHealthReportBreaksDownByEcosystem(t *testing.T) {
	report := BuildHealthReport("o", "r", github.DependabotFilter(), []github.PullRequest{
		openPR(1, 3, "npm_and_yarn"),
		mergedPR(2, 10, 1, "npm_and_yarn"),
		mergedPR(3, 11, 1, "npm_and_yarn"),
		mergedPR(4, 12, 1, "go_modules"),
	}, reportNow)

	assert.Equal(t, len(report.Ecosystems), 2)
	// Busiest first.
	assert.Equal(t, report.Ecosystems[0].Ecosystem, "npm")
	assert.Equal(t, report.Ecosystems[0].Total, 3)
	assert.Equal(t, report.Ecosystems[0].MergeRate, 100.0)
	assert.Equal(t, report.Ecosystems[1].Ecosystem, "gomod")
}

func TestBuildHealthReportGroupsUnattributablePullRequests(t *testing.T) {
	report := BuildHealthReport("o", "r", github.DependabotFilter(), []github.PullRequest{
		{Number: 1, State: "OPEN", CreatedAt: at(2), HeadRefName: "some/other/branch"},
	}, reportNow)

	assert.Equal(t, len(report.Ecosystems), 1)
	assert.Equal(t, report.Ecosystems[0].Ecosystem, "unknown")
}

func TestBuildHealthReportRanksTheOldestOpenPullRequests(t *testing.T) {
	report := BuildHealthReport("o", "r", github.DependabotFilter(), []github.PullRequest{
		openPR(1, 3, "pip"), openPR(2, 45, "pip"), openPR(3, 9, "pip"),
		openPR(4, 60, "pip"), openPR(5, 1, "pip"), openPR(6, 30, "pip"), openPR(7, 12, "pip"),
	}, reportNow)

	assert.Assert(t, report.OldestOpen != nil)
	assert.Equal(t, report.OldestOpen.Number, 4)
	assert.Equal(t, days(report.OldestOpenAge), 60)
	// Capped at five, oldest first.
	assert.Equal(t, len(report.OldestOpenList), 5)
	assert.Equal(t, report.OldestOpenList[0].PullRequest.Number, 4)
	assert.Equal(t, report.OldestOpenList[1].PullRequest.Number, 2)
}

func TestBuildHealthReportHandlesNoPullRequests(t *testing.T) {
	report := BuildHealthReport("o", "r", github.DependabotFilter(), nil, reportNow)

	assert.Equal(t, report.Total, 0)
	assert.Equal(t, report.MergeRate, 0.0)
	assert.Assert(t, report.OldestOpen == nil)
	assert.Assert(t, strings.Contains(report.Markdown(), "No pull requests opened by app/dependabot"))
}

func TestBuildHealthReportIgnoresUnparseableTimestamps(t *testing.T) {
	broken := github.PullRequest{Number: 9, State: "OPEN", CreatedAt: "not-a-date"}

	report := BuildHealthReport("o", "r", github.DependabotFilter(), []github.PullRequest{broken}, reportNow)

	assert.Equal(t, report.Open, 1)
	assert.Assert(t, report.OldestOpen == nil, "an unparseable date must not become the oldest")
}

func TestMarkdownReportsCountsAndAVerdict(t *testing.T) {
	markdown := BuildHealthReport("kaypee90", "zebrok", github.DependabotFilter(), []github.PullRequest{
		openPR(1, 3, "npm_and_yarn"),
		mergedPR(2, 10, 1, "npm_and_yarn"),
		mergedPR(3, 12, 1, "npm_and_yarn"),
		closedPR(4, 20, "npm_and_yarn"),
	}, reportNow).Markdown()

	assert.Assert(t, strings.Contains(markdown, "## Dependabot health: kaypee90/zebrok"))
	assert.Assert(t, strings.Contains(markdown, "| Total pull requests | 4 |"))
	assert.Assert(t, strings.Contains(markdown, "| Merge rate | 67% |"))
	assert.Assert(t, strings.Contains(markdown, "| Median time to merge | 1 day |"))
	assert.Assert(t, strings.Contains(markdown, "keeping up"))
}

// TestMarkdownEmitsGroupingYamlForALargeQueue is the loop the whole feature
// exists to close: the finding names the configuration that fixes it, and
// that yaml is rendered from the same model the generator writes.
func TestMarkdownEmitsGroupingYamlForALargeQueue(t *testing.T) {
	var pullRequests []github.PullRequest
	for i := 1; i <= 15; i++ {
		pullRequests = append(pullRequests, openPR(i, i, "npm_and_yarn"))
	}
	pullRequests = append(pullRequests, mergedPR(90, 30, 1, "npm_and_yarn"))

	markdown := BuildHealthReport("o", "r", github.DependabotFilter(), pullRequests, reportNow).Markdown()

	assert.Assert(t, strings.Contains(markdown, "is a queue, not a backlog"))
	assert.Assert(t, strings.Contains(markdown, "```yaml"))
	assert.Assert(t, strings.Contains(markdown, "groups:"))
	assert.Assert(t, strings.Contains(markdown, "all-minor-and-patch:"))
	assert.Assert(t, strings.Contains(markdown, "- minor"))
}

func TestMarkdownEmitsIgnoreYamlForALowMergeRate(t *testing.T) {
	markdown := BuildHealthReport("o", "r", github.DependabotFilter(), []github.PullRequest{
		openPR(1, 1, "pip"),
		closedPR(2, 5, "pip"), closedPR(3, 6, "pip"), closedPR(4, 7, "pip"),
		mergedPR(5, 8, 1, "pip"),
	}, reportNow).Markdown()

	assert.Assert(t, strings.Contains(markdown, "ignore:"))
	assert.Assert(t, strings.Contains(markdown, "version-update:semver-major"))
}

// A weak merge rate over three resolved pull requests is not evidence.
func TestMarkdownDoesNotRecommendOnThinEvidence(t *testing.T) {
	markdown := BuildHealthReport("o", "r", github.DependabotFilter(), []github.PullRequest{
		openPR(1, 1, "pip"), closedPR(2, 5, "pip"), closedPR(3, 6, "pip"), mergedPR(4, 7, 1, "pip"),
	}, reportNow).Markdown()

	assert.Assert(t, !strings.Contains(markdown, "ignore:"),
		"three resolved pull requests is too small a sample to recommend on:\n%s", markdown)
}

func TestMarkdownFlagsStalePullRequests(t *testing.T) {
	markdown := BuildHealthReport("o", "r", github.DependabotFilter(), []github.PullRequest{
		openPR(1, 90, "pip"), mergedPR(2, 10, 1, "pip"),
	}, reportNow).Markdown()

	assert.Assert(t, strings.Contains(markdown, "open longer than 30 days"))
	assert.Assert(t, strings.Contains(markdown, "open-pull-requests-limit"))
}

func TestMarkdownFlagsSlowReview(t *testing.T) {
	markdown := BuildHealthReport("o", "r", github.DependabotFilter(), []github.PullRequest{
		mergedPR(1, 60, 20, "pip"), mergedPR(2, 50, 25, "pip"), mergedPR(3, 40, 30, "pip"),
	}, reportNow).Markdown()

	assert.Assert(t, strings.Contains(markdown, "Median time to merge is"))
}

func TestMarkdownIncludesTheEcosystemBreakdownOnlyWhenItSaysSomething(t *testing.T) {
	oneEcosystem := BuildHealthReport("o", "r", github.DependabotFilter(), []github.PullRequest{
		openPR(1, 1, "pip"), mergedPR(2, 5, 1, "pip"),
	}, reportNow).Markdown()
	assert.Assert(t, !strings.Contains(oneEcosystem, "### By ecosystem"))

	twoEcosystems := BuildHealthReport("o", "r", github.DependabotFilter(), []github.PullRequest{
		openPR(1, 1, "pip"), mergedPR(2, 5, 1, "go_modules"),
	}, reportNow).Markdown()
	assert.Assert(t, strings.Contains(twoEcosystems, "### By ecosystem"))
	assert.Assert(t, strings.Contains(twoEcosystems, "`gomod`"))
}

func TestMarkdownListsWhatIsWaitingLongest(t *testing.T) {
	markdown := BuildHealthReport("o", "r", github.DependabotFilter(), []github.PullRequest{
		openPR(41, 34, "npm_and_yarn"), openPR(42, 2, "go_modules"),
	}, reportNow).Markdown()

	assert.Assert(t, strings.Contains(markdown, "### Waiting longest"))
	assert.Assert(t, strings.Contains(markdown, "[#41](https://example.test/41)"))
	assert.Assert(t, strings.Contains(markdown, "34 days"))
}

func TestHumanDurationUsesTheCoarsestUsefulUnit(t *testing.T) {
	assert.Equal(t, humanDuration(72*time.Hour), "3 days")
	assert.Equal(t, humanDuration(24*time.Hour), "1 day")
	assert.Equal(t, humanDuration(5*time.Hour), "5 hours")
	assert.Equal(t, humanDuration(time.Hour), "1 hour")
	assert.Equal(t, humanDuration(30*time.Minute), "30 minutes")
}

func TestCountOfPluralisesCorrectly(t *testing.T) {
	assert.Equal(t, countOf(1, "pull request"), "1 pull request")
	assert.Equal(t, countOf(3, "pull request"), "3 pull requests")
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

	err := runReport("kaypee90", "zebrok", github.DependabotFilter())

	assert.ErrorContains(t, err, "GITHUB_TOKEN is not set")
}

func TestRunReportRequiresARepository(t *testing.T) {
	t.Setenv("GITHUB_TOKEN", "token")
	t.Setenv("GITHUB_REPOSITORY", "")

	err := runReport("", "", github.DependabotFilter())

	assert.ErrorContains(t, err, "--owner and --repo are required")
}
