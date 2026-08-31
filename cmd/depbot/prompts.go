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

func addAdditionalPackageManager() (string, error) {
	prompt := promptui.Select{
		Label: "Do you want to add another package manager?",
		Items: []string{YES, NO},
	}

	_, result, err := prompt.Run()

	return result, err
}
