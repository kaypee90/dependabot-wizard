package main

import (
	"strconv"

	"github.com/manifoldco/promptui"
)

const (
	YES = "yes"
	NO  = "no"
)

func getConfigurationOverrideConfirmation() (string, error) {
	prompt := promptui.Select{
		Label: "Dependabot is currently configured for this project. Would you like to override the current configuration?",
		Items: []string{YES, NO},
	}

	_, result, err := prompt.Run()

	return result, err
}

func getPackageEcosystem() (string, error) {
	prompt := promptui.Select{
		Label: "Select Package Ecosystem",
		Items: []string{
			"cargo", "docker", "gomod", "gradle", "maven",
			"npm", "nuget", "pip", "poetry", "swift", "terraform",
			"bundler", "composer", "devcontainers", "elm",
			"gitsubmodule", "github-actions", "mix", "pub",
		},
	}

	_, result, err := prompt.Run()

	return result, err
}

func getDirectory() (string, error) {
	prompt := promptui.Prompt{
		Label:   "Provide Directory (/)",
		Default: "/",
	}

	result, err := prompt.Run()

	return result, err
}

func getInterval() (string, error) {
	prompt := promptui.Select{
		Label: "Select Interval",
		Items: []string{"daily", "weekly", "monthly"},
	}

	_, result, err := prompt.Run()

	return result, err
}

func getReviewer() (string, error) {
	prompt := promptui.Prompt{
		Label: "Provide Reviewer (optional)",
	}

	result, err := prompt.Run()

	return result, err
}

func getOpenPullRequestLimit() (uint64, error) {
	prompt := promptui.Prompt{
		Label: "Provide Open Pull Request Limit (optional)",
	}

	result, err := prompt.Run()
	if err != nil {
		return 0, err
	}

	limit, err := strconv.ParseUint(result, 10, 32)
	if err != nil {
		return 5, nil
	}

	return limit, nil
}

func getLabel() (string, error) {
	prompt := promptui.Prompt{
		Label: "Provide Label (optional)",
	}

	result, err := prompt.Run()

	return result, err
}

func addAdditionalPackageManager() (string, error) {
	prompt := promptui.Select{
		Label: "Do you want to add another package manager?",
		Items: []string{YES, NO},
	}

	_, result, err := prompt.Run()

	return result, err
}
