package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/manifoldco/promptui"
)

const (
	YES = "yes"
	NO  = "no"

	defaultOpenPullRequestsLimit = 5
)

func getConfigurationOverrideConfirmation() (string, error) {
	prompt := promptui.Select{
		Label: "Dependabot is currently configured for this project. Would you like to override the current configuration?",
		Items: []string{YES, NO},
	}

	_, result, err := prompt.Run()

	return result, err
}

// useDetectedConfiguration asks whether to accept the scan results instead of
// answering the questions one ecosystem at a time.
func useDetectedConfiguration() (string, error) {
	prompt := promptui.Select{
		Label: "Use this detected configuration?",
		Items: []string{YES, NO},
	}

	_, result, err := prompt.Run()

	return result, err
}

func getPackageEcosystem() (string, error) {
	prompt := promptui.Select{
		Label: "Select Package Ecosystem",
		Items: supportedEcosystems,
	}

	_, result, err := prompt.Run()

	return result, err
}

func getDirectory() (string, error) {
	prompt := promptui.Prompt{
		Label:   "Provide Directory (/)",
		Default: "/",
		Validate: func(input string) error {
			if !strings.HasPrefix(strings.TrimSpace(input), "/") {
				return fmt.Errorf("directory must start with /")
			}
			return nil
		},
	}

	result, err := prompt.Run()

	return strings.TrimSpace(result), err
}

func getInterval() (string, error) {
	prompt := promptui.Select{
		Label: "Select Interval",
		Items: supportedIntervals,
	}

	_, result, err := prompt.Run()

	return result, err
}

func getReviewer() (string, error) {
	prompt := promptui.Prompt{
		Label: "Provide Reviewer (optional)",
	}

	result, err := prompt.Run()

	return strings.TrimSpace(result), err
}

// getOpenPullRequestLimit accepts a blank answer as "use the default" but
// rejects a non-numeric one, rather than silently substituting a value the
// user did not ask for.
func getOpenPullRequestLimit() (uint64, error) {
	prompt := promptui.Prompt{
		Label:   "Provide Open Pull Request Limit",
		Default: strconv.Itoa(defaultOpenPullRequestsLimit),
		Validate: func(input string) error {
			trimmed := strings.TrimSpace(input)
			if trimmed == "" {
				return nil
			}
			if _, err := strconv.ParseUint(trimmed, 10, 32); err != nil {
				return fmt.Errorf("enter a whole number, for example 5")
			}
			return nil
		},
	}

	result, err := prompt.Run()
	if err != nil {
		return 0, err
	}

	trimmed := strings.TrimSpace(result)
	if trimmed == "" {
		return defaultOpenPullRequestsLimit, nil
	}

	limit, err := strconv.ParseUint(trimmed, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("reading open pull request limit: %w", err)
	}

	return limit, nil
}

func getLabel() (string, error) {
	prompt := promptui.Prompt{
		Label: "Provide Label (optional)",
	}

	result, err := prompt.Run()

	return strings.TrimSpace(result), err
}

// getGroupingStrategy asks how to batch updates. Grouping is the single
// biggest reducer of pull request volume, so it is worth asking about rather
// than leaving people to discover it after the queue is unmanageable.
func getGroupingStrategy() (GroupingStrategy, error) {
	items := make([]string, 0, len(supportedGroupingStrategies))
	for _, strategy := range supportedGroupingStrategies {
		items = append(items, DescribeGrouping(strategy))
	}

	prompt := promptui.Select{
		Label: "How should updates be batched?",
		Items: items,
		Size:  len(items),
	}

	index, _, err := prompt.Run()
	if err != nil {
		return GroupingNone, err
	}

	return supportedGroupingStrategies[index], nil
}

// getIgnoreMajorUpdates asks whether to stop opening major version bumps.
func getIgnoreMajorUpdates() (bool, error) {
	prompt := promptui.Select{
		Label: "Open pull requests for major version updates?",
		Items: []string{YES, NO},
	}

	_, result, err := prompt.Run()
	if err != nil {
		return false, err
	}

	return result == NO, nil
}

func addAdditionalPackageManager() (string, error) {
	prompt := promptui.Select{
		Label: "Do you want to add another package manager?",
		Items: []string{YES, NO},
	}

	_, result, err := prompt.Run()

	return result, err
}
