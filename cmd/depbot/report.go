package main

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/kaypee90/dependabot-wizard/github"
	"gopkg.in/yaml.v2"
)

const (
	// staleAfter is when an open dependency update stops being a queue and
	// starts being a liability: the branch has drifted, and the update it
	// carries is usually superseded.
	staleAfter = 30 * 24 * time.Hour

	// recentWindow is the period the throughput figures cover.
	recentWindow = 30 * 24 * time.Hour

	// oldestOpenListSize caps the "waiting longest" table.
	oldestOpenListSize = 5
)

// ecosystemFromBranch maps dependabot's branch slug to the name used in
// dependabot.yml, so the report and the config speak the same language.
var ecosystemFromBranch = map[string]string{
	"npm_and_yarn":   "npm",
	"go_modules":     "gomod",
	"github_actions": "github-actions",
	"submodules":     "gitsubmodule",
	"pip":            "pip",
	"uv":             "pip",
	"bundler":        "bundler",
	"docker":         "docker",
	"composer":       "composer",
	"cargo":          "cargo",
	"maven":          "maven",
	"gradle":         "gradle",
	"nuget":          "nuget",
	"terraform":      "terraform",
	"hex":            "mix",
	"pub":            "pub",
	"swift":          "swift",
	"elm":            "elm",
	"devcontainers":  "devcontainers",
}

// EcosystemStat is one row of the per-ecosystem breakdown.
type EcosystemStat struct {
	Ecosystem string
	Total     int
	Open      int
	Merged    int
	Closed    int
	MergeRate float64
}

// AgedPullRequest pairs a pull request with how long it has been waiting.
type AgedPullRequest struct {
	PullRequest github.PullRequest
	Age         time.Duration
}

// HealthReport summarises how a repository is handling its dependency pull
// requests. It answers the question the config leads to: is dependabot
// actually working here, or is nobody merging what it opens?
type HealthReport struct {
	Owner  string
	Repo   string
	Filter github.Filter

	Total  int
	Open   int
	Merged int
	Closed int

	// MergeRate is merged as a percentage of resolved (merged plus closed)
	// pull requests. Open ones are excluded because they are undecided.
	MergeRate float64

	// StaleOpen counts open pull requests older than staleAfter.
	StaleOpen int

	// MedianTimeToMerge is the median age at merge, over merged pull requests
	// with a usable timestamp. Median rather than mean, because one pull
	// request left open for a year should not move the headline number.
	MedianTimeToMerge time.Duration
	HasMergeTimes     bool

	OpenedRecently int
	MergedRecently int

	OldestOpen     *github.PullRequest
	OldestOpenAge  time.Duration
	OldestOpenList []AgedPullRequest

	Ecosystems []EcosystemStat
}

func parseTime(value string) (time.Time, bool) {
	if value == "" {
		return time.Time{}, false
	}

	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return time.Time{}, false
	}

	return parsed, true
}

// EcosystemForBranch extracts the package ecosystem from a dependabot branch
// name such as "dependabot/npm_and_yarn/ui/web/lodash-4.17.21".
func EcosystemForBranch(branch string) string {
	segments := strings.Split(branch, "/")
	if len(segments) < 2 || segments[0] != "dependabot" {
		return ""
	}

	slug := segments[1]
	if mapped, ok := ecosystemFromBranch[slug]; ok {
		return mapped
	}

	return slug
}

func median(durations []time.Duration) time.Duration {
	if len(durations) == 0 {
		return 0
	}

	sort.Slice(durations, func(i, j int) bool { return durations[i] < durations[j] })

	middle := len(durations) / 2
	if len(durations)%2 == 1 {
		return durations[middle]
	}

	return (durations[middle-1] + durations[middle]) / 2
}

// BuildHealthReport derives the summary from raw pull requests. It is a pure
// function so the reporting logic is testable without touching the network.
func BuildHealthReport(owner string, repo string, filter github.Filter, pullRequests []github.PullRequest, now time.Time) HealthReport {
	report := HealthReport{Owner: owner, Repo: repo, Filter: filter, Total: len(pullRequests)}

	perEcosystem := map[string]*EcosystemStat{}
	var mergeDurations []time.Duration
	var openWithAge []AgedPullRequest

	for index := range pullRequests {
		pullRequest := pullRequests[index]
		state := strings.ToUpper(pullRequest.State)

		ecosystem := EcosystemForBranch(pullRequest.HeadRefName)
		if ecosystem == "" {
			ecosystem = "unknown"
		}
		stat, seen := perEcosystem[ecosystem]
		if !seen {
			stat = &EcosystemStat{Ecosystem: ecosystem}
			perEcosystem[ecosystem] = stat
		}
		stat.Total++

		created, hasCreated := parseTime(pullRequest.CreatedAt)

		switch state {
		case "OPEN":
			report.Open++
			stat.Open++

			if hasCreated {
				age := now.Sub(created)
				openWithAge = append(openWithAge, AgedPullRequest{PullRequest: pullRequest, Age: age})

				if age > staleAfter {
					report.StaleOpen++
				}
				if now.Sub(created) <= recentWindow {
					report.OpenedRecently++
				}
			}

		case "MERGED":
			report.Merged++
			stat.Merged++

			if merged, ok := parseTime(pullRequest.MergedAt); ok {
				if hasCreated && merged.After(created) {
					mergeDurations = append(mergeDurations, merged.Sub(created))
				}
				if now.Sub(merged) <= recentWindow {
					report.MergedRecently++
				}
			}

		case "CLOSED":
			report.Closed++
			stat.Closed++
		}
	}

	if resolved := report.Merged + report.Closed; resolved > 0 {
		report.MergeRate = float64(report.Merged) / float64(resolved) * 100
	}

	if len(mergeDurations) > 0 {
		report.MedianTimeToMerge = median(mergeDurations)
		report.HasMergeTimes = true
	}

	sort.Slice(openWithAge, func(i, j int) bool { return openWithAge[i].Age > openWithAge[j].Age })
	if len(openWithAge) > 0 {
		oldest := openWithAge[0]
		report.OldestOpen = &oldest.PullRequest
		report.OldestOpenAge = oldest.Age

		limit := len(openWithAge)
		if limit > oldestOpenListSize {
			limit = oldestOpenListSize
		}
		report.OldestOpenList = openWithAge[:limit]
	}

	for _, stat := range perEcosystem {
		if resolved := stat.Merged + stat.Closed; resolved > 0 {
			stat.MergeRate = float64(stat.Merged) / float64(resolved) * 100
		}
		report.Ecosystems = append(report.Ecosystems, *stat)
	}

	// Busiest first, then alphabetical, so the table is stable.
	sort.Slice(report.Ecosystems, func(i, j int) bool {
		if report.Ecosystems[i].Total != report.Ecosystems[j].Total {
			return report.Ecosystems[i].Total > report.Ecosystems[j].Total
		}
		return report.Ecosystems[i].Ecosystem < report.Ecosystems[j].Ecosystem
	})

	return report
}

func days(duration time.Duration) int {
	return int(duration.Hours() / 24)
}

// countOf renders a count with a correctly pluralised noun.
func countOf(count int, singular string) string {
	if count == 1 {
		return fmt.Sprintf("%d %s", count, singular)
	}

	return fmt.Sprintf("%d %ss", count, singular)
}

// humanDuration renders a duration at the coarsest useful precision.
func humanDuration(duration time.Duration) string {
	if duration >= 24*time.Hour {
		return countOf(days(duration), "day")
	}
	if duration >= time.Hour {
		return countOf(int(duration.Hours()), "hour")
	}
	return countOf(int(duration.Minutes()), "minute")
}

// yamlSnippet renders a fragment of dependabot.yml from the real model, so
// the recommendation cannot drift from what the generator would produce.
func yamlSnippet(value any) string {
	rendered, err := yaml.Marshal(value)
	if err != nil {
		return ""
	}

	var indented strings.Builder
	for _, line := range strings.Split(strings.TrimRight(string(rendered), "\n"), "\n") {
		indented.WriteString("  " + line + "\n")
	}

	return indented.String()
}

// Markdown renders the report for a GitHub Actions step summary or a comment.
func (report HealthReport) Markdown() string {
	var builder strings.Builder

	fmt.Fprintf(&builder, "## Dependabot health: %s/%s\n\n", report.Owner, report.Repo)

	if report.Total == 0 {
		fmt.Fprintf(&builder, "No pull requests %s were found.\n\n", report.Filter.Describe())
		fmt.Fprintf(&builder, "That usually means dependabot is not configured for this repository yet. "+
			"Run `depbot` here to generate a configuration.\n")
		return builder.String()
	}

	report.writeSummary(&builder)
	report.writeEcosystems(&builder)
	report.writeOldestOpen(&builder)
	report.writeRecommendations(&builder)

	return builder.String()
}

func (report HealthReport) writeSummary(builder *strings.Builder) {
	fmt.Fprintf(builder, "| Metric | Value |\n| --- | ---: |\n")
	fmt.Fprintf(builder, "| Total pull requests | %d |\n", report.Total)
	fmt.Fprintf(builder, "| Open | %d |\n", report.Open)
	fmt.Fprintf(builder, "| Merged | %d |\n", report.Merged)
	fmt.Fprintf(builder, "| Closed without merging | %d |\n", report.Closed)

	if report.Merged+report.Closed > 0 {
		fmt.Fprintf(builder, "| Merge rate | %.0f%% |\n", report.MergeRate)
	}

	if report.HasMergeTimes {
		fmt.Fprintf(builder, "| Median time to merge | %s |\n", humanDuration(report.MedianTimeToMerge))
	}

	if report.Open > 0 {
		fmt.Fprintf(builder, "| Open longer than %s | %d |\n", humanDuration(staleAfter), report.StaleOpen)
	}

	fmt.Fprintf(builder, "| Opened in the last %s | %d |\n", humanDuration(recentWindow), report.OpenedRecently)
	fmt.Fprintf(builder, "| Merged in the last %s | %d |\n", humanDuration(recentWindow), report.MergedRecently)

	builder.WriteString("\n")
}

func (report HealthReport) writeEcosystems(builder *strings.Builder) {
	// One ecosystem tells you nothing the summary did not already say.
	if len(report.Ecosystems) < 2 {
		return
	}

	builder.WriteString("### By ecosystem\n\n")
	builder.WriteString("| Ecosystem | Total | Open | Merged | Closed | Merge rate |\n")
	builder.WriteString("| --- | ---: | ---: | ---: | ---: | ---: |\n")

	for _, stat := range report.Ecosystems {
		mergeRate := "—"
		if stat.Merged+stat.Closed > 0 {
			mergeRate = fmt.Sprintf("%.0f%%", stat.MergeRate)
		}
		fmt.Fprintf(builder, "| `%s` | %d | %d | %d | %d | %s |\n",
			stat.Ecosystem, stat.Total, stat.Open, stat.Merged, stat.Closed, mergeRate)
	}

	builder.WriteString("\n")
}

func (report HealthReport) writeOldestOpen(builder *strings.Builder) {
	if len(report.OldestOpenList) == 0 {
		return
	}

	builder.WriteString("### Waiting longest\n\n")
	builder.WriteString("| Pull request | Ecosystem | Age |\n| --- | --- | ---: |\n")

	for _, aged := range report.OldestOpenList {
		ecosystem := EcosystemForBranch(aged.PullRequest.HeadRefName)
		if ecosystem == "" {
			ecosystem = "—"
		}

		title := aged.PullRequest.Title
		if len(title) > 60 {
			title = title[:57] + "…"
		}

		fmt.Fprintf(builder, "| [#%d](%s) %s | `%s` | %s |\n",
			aged.PullRequest.Number, aged.PullRequest.Url, title, ecosystem, humanDuration(aged.Age))
	}

	builder.WriteString("\n")
}

// writeRecommendations closes the loop: every finding names the configuration
// that addresses it, rendered from the same model the generator uses.
func (report HealthReport) writeRecommendations(builder *strings.Builder) {
	builder.WriteString("### What to do\n\n")

	wrote := false

	if report.Open > 10 {
		wrote = true
		fmt.Fprintf(builder, "**%d open pull requests is a queue, not a backlog.** Group minor and "+
			"patch updates so they arrive together and majors stay reviewable on their own:\n\n",
			report.Open)
		builder.WriteString("```yaml\n")
		builder.WriteString(yamlSnippet(map[string]any{"groups": GroupsFor(GroupingAll)}))
		builder.WriteString("```\n\n")
		builder.WriteString("Run `depbot` and choose *" + DescribeGrouping(GroupingAll) + "* to generate this.\n\n")
	}

	if report.Merged+report.Closed >= 4 && report.MergeRate < 50 {
		wrote = true
		fmt.Fprintf(builder, "**Only %.0f%% of resolved pull requests were merged.** Most updates are "+
			"being closed rather than taken, which is usually major versions arriving faster than "+
			"anyone can absorb them. Stop opening those and take them deliberately:\n\n", report.MergeRate)
		builder.WriteString("```yaml\n")
		builder.WriteString(yamlSnippet(map[string]any{"ignore": IgnoreMajorUpdates()}))
		builder.WriteString("```\n\n")
	}

	if report.StaleOpen > 0 {
		wrote = true
		verb := "have"
		if report.StaleOpen == 1 {
			verb = "has"
		}
		fmt.Fprintf(builder, "**%s %s been open longer than %s.** Those branches have drifted and the "+
			"updates they carry are probably superseded; closing them costs nothing because "+
			"dependabot reopens what is still needed. Lowering `open-pull-requests-limit` keeps the "+
			"queue at a size that gets reviewed.\n\n",
			countOf(report.StaleOpen, "pull request"), verb, humanDuration(staleAfter))
	}

	if report.HasMergeTimes && report.MedianTimeToMerge > 14*24*time.Hour {
		wrote = true
		fmt.Fprintf(builder, "**Median time to merge is %s.** Updates are being reviewed, but slowly "+
			"enough that each one is stale on arrival. A `weekly` schedule with grouping usually "+
			"beats `daily` without it.\n\n", humanDuration(report.MedianTimeToMerge))
	}

	if report.Open == 0 && report.Total > 0 {
		wrote = true
		builder.WriteString("**Nothing is waiting.** Dependabot updates are being kept up with here.\n\n")
	}

	if !wrote {
		fmt.Fprintf(builder, "**%d open, %.0f%% of resolved pull requests merged.** This repository is "+
			"keeping up; no configuration change is called for.\n\n", report.Open, report.MergeRate)
	}
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
func runReport(owner string, repo string, filter github.Filter) error {
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

	pullRequests, err := github.GetPullRequests(context.Background(), owner, repo, filter, token)
	if err != nil {
		return err
	}

	report := BuildHealthReport(owner, repo, filter, pullRequests, time.Now())
	markdown := report.Markdown()

	fmt.Print(markdown)

	if err := appendToFileFromEnv("GITHUB_STEP_SUMMARY", markdown); err != nil {
		return err
	}

	outputs := fmt.Sprintf("total=%d\nopen=%d\nmerged=%d\nclosed=%d\nmerge-rate=%.0f\nstale=%d\nmedian-merge-days=%d\n",
		report.Total, report.Open, report.Merged, report.Closed, report.MergeRate,
		report.StaleOpen, days(report.MedianTimeToMerge))

	return appendToFileFromEnv("GITHUB_OUTPUT", outputs)
}
