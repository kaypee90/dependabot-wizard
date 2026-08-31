package main

import (
	"errors"
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/kaypee90/dependabot-wizard/github"
	"github.com/manifoldco/promptui"
)

// version is injected at build time:
//
//	go build -ldflags "-X main.version=$(git describe --tags --always)"
var version = "dev"

// errCancelled means the user interrupted a prompt; it is not a failure.
var errCancelled = errors.New("cancelled")

func displayAppVersion() {
	fmt.Printf("Depbot %s\n", version)
}

func displayHelpText() {
	fmt.Print(getCliHelpText())
}

// promptError converts a prompt failure into either a cancellation or a
// real error, so Ctrl-C exits quietly and everything else is reported.
func promptError(err error) error {
	if errors.Is(err, promptui.ErrInterrupt) || errors.Is(err, promptui.ErrEOF) {
		return errCancelled
	}

	return err
}

func launchWebApplication(port int) error {
	// Bind before announcing anything, so a busy port fails loudly here
	// instead of opening a browser onto a page that will never load.
	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		return fmt.Errorf("cannot listen on port %d: %w", port, err)
	}

	url := fmt.Sprintf("http://%s/", listener.Addr().String())
	fmt.Printf("Depbot is running at %s\nPress Ctrl-C to stop.\n", url)

	served := make(chan error, 1)
	go func() {
		served <- serve(listener)
	}()

	// DEPBOT_NO_BROWSER keeps the server headless, for CI and smoke tests.
	if os.Getenv("DEPBOT_NO_BROWSER") != "" {
		fmt.Printf("Running headless. Open %s to continue.\n", url)
	} else if err := openBrowser(url); err != nil {
		fmt.Printf("Could not open a browser automatically. Open %s to continue.\n", url)
	}

	return <-served
}

// collectUpdatesInteractively runs the original question-by-question flow.
func collectUpdatesInteractively() ([]Update, error) {
	var updates []Update

	for {
		packageEcosystem, err := getPackageEcosystem()
		if err != nil {
			return nil, promptError(err)
		}

		directory, err := getDirectory()
		if err != nil {
			return nil, promptError(err)
		}

		interval, err := getInterval()
		if err != nil {
			return nil, promptError(err)
		}

		var reviewers []string
		reviewer, err := getReviewer()
		if err != nil {
			return nil, promptError(err)
		}
		if reviewer != "" {
			reviewers = append(reviewers, reviewer)
		}

		openPullRequestsLimit, err := getOpenPullRequestLimit()
		if err != nil {
			return nil, promptError(err)
		}

		var labels []string
		label, err := getLabel()
		if err != nil {
			return nil, promptError(err)
		}
		if label != "" {
			labels = append(labels, label)
		}

		updates = append(updates, Update{
			PackageEcosystem: packageEcosystem,
			Directory:        directory,
			Schedule: Schedule{
				Interval: interval,
			},
			Reviewers:             reviewers,
			OpenPullRequestsLimit: openPullRequestsLimit,
			Labels:                labels,
		})

		again, err := addAdditionalPackageManager()
		if err != nil {
			return nil, promptError(err)
		}
		if again == NO {
			return updates, nil
		}
	}
}

// collectUpdatesFromDetection offers whatever the scanner found. It returns
// nil when there is nothing to offer or the user declines, so the caller can
// fall back to the manual flow.
func collectUpdatesFromDetection() ([]Update, error) {
	detected, err := DetectEcosystemsInWorkingDirectory()
	if err != nil {
		// Detection is a convenience; never block the wizard on it.
		fmt.Printf("Could not scan this project (%v). Falling back to manual setup.\n", err)
		return nil, nil
	}

	if len(detected) == 0 {
		return nil, nil
	}

	fmt.Printf("\nFound %d package %s in this project:\n", len(detected), pluralize(len(detected), "ecosystem", "ecosystems"))
	for _, item := range detected {
		fmt.Printf("  %-15s %-12s (from %s)\n", item.Ecosystem, item.Directory, item.Evidence)
	}
	fmt.Println()

	accepted, err := useDetectedConfiguration()
	if err != nil {
		return nil, promptError(err)
	}
	if accepted == NO {
		return nil, nil
	}

	interval, err := getInterval()
	if err != nil {
		return nil, promptError(err)
	}

	openPullRequestsLimit, err := getOpenPullRequestLimit()
	if err != nil {
		return nil, promptError(err)
	}

	return UpdatesFromDetection(detected, interval, openPullRequestsLimit), nil
}

// applyNoiseControls attaches the grouping and ignore rules to every update
// block. These are repository-wide decisions rather than per-ecosystem ones,
// so they are asked once and applied across the board.
func applyNoiseControls(updates []Update, strategy GroupingStrategy, ignoreMajor bool) []Update {
	groups := GroupsFor(strategy)

	for index := range updates {
		if len(groups) > 0 {
			updates[index].Groups = groups
		}
		if ignoreMajor {
			updates[index].Ignore = IgnoreMajorUpdates()
		}
	}

	return updates
}

// askNoiseControls collects the grouping and major-update decisions.
func askNoiseControls() (GroupingStrategy, bool, error) {
	strategy, err := getGroupingStrategy()
	if err != nil {
		return GroupingNone, false, promptError(err)
	}

	ignoreMajor, err := getIgnoreMajorUpdates()
	if err != nil {
		return GroupingNone, false, promptError(err)
	}

	return strategy, ignoreMajor, nil
}

func pluralize(count int, singular string, plural string) string {
	if count == 1 {
		return singular
	}

	return plural
}

func launchCliApplication() error {
	printIntroductoryText()

	exists, err := dependabotFileExists()
	if err != nil {
		return err
	}
	if exists {
		result, err := getConfigurationOverrideConfirmation()
		if err != nil {
			return promptError(err)
		}
		if result == NO {
			return errCancelled
		}
	}

	updates, err := collectUpdatesFromDetection()
	if err != nil {
		return err
	}

	if len(updates) == 0 {
		updates, err = collectUpdatesInteractively()
		if err != nil {
			return err
		}
	}

	strategy, ignoreMajor, err := askNoiseControls()
	if err != nil {
		return err
	}
	updates = applyNoiseControls(updates, strategy, ignoreMajor)

	config := Config{
		Version: 2,
		Updates: updates,
	}

	yamlData, err := config.ConvertToYaml()
	if err != nil {
		return err
	}

	path, err := createDependabotYamlFile(yamlData)
	if err != nil {
		return err
	}

	fmt.Printf("\nWrote %s\n", path)

	return nil
}

func run() error {
	var (
		showVersion  bool
		startWebApp  bool
		displayHelp  bool
		reportHealth bool
		port         int
		owner        string
		repo         string
		label        string
		author       string
	)

	flag.BoolVar(&showVersion, "version", false, "Display app version")
	flag.BoolVar(&showVersion, "v", false, "Display app version (shorthand)")
	flag.BoolVar(&startWebApp, "web", false, "Start web application")
	flag.BoolVar(&startWebApp, "w", false, "Start web application (shorthand)")
	flag.BoolVar(&displayHelp, "help", false, "Display help text")
	flag.BoolVar(&displayHelp, "h", false, "Display help text (shorthand)")
	flag.BoolVar(&reportHealth, "report", false, "Print a dependabot health report for a repository")
	flag.IntVar(&port, "port", 3001, "Port for the web application")
	flag.StringVar(&owner, "owner", "", "Repository owner, for --report")
	flag.StringVar(&repo, "repo", "", "Repository name, for --report")
	flag.StringVar(&label, "label", "", "Only count pull requests with this label, for --report")
	flag.StringVar(&author, "author", github.DefaultAuthor, "Pull request author to count, for --report")

	flag.Usage = displayHelpText
	flag.Parse()

	switch {
	case showVersion:
		displayAppVersion()
		return nil
	case displayHelp:
		displayHelpText()
		return nil
	case reportHealth:
		return runReport(owner, repo, github.Filter{Author: author, Label: label})
	case startWebApp:
		return launchWebApplication(port)
	default:
		return launchCliApplication()
	}
}

func main() {
	if err := run(); err != nil {
		if errors.Is(err, errCancelled) {
			os.Exit(0)
		}
		fmt.Fprintf(os.Stderr, "depbot: %v\n", err)
		os.Exit(1)
	}
}
