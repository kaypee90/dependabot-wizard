package main

import (
	"strconv"

	"github.com/manifoldco/promptui"
)

const YES = "yes"
const NO = "no"

func getPackageEcosystem() string {
	prompt := promptui.Select{
		Label: "Select Package Ecosystem",
		Items: []string{"cargo", "docker", "gomod", "gradle", "maven",
			"npm", "nuget", "pip", "poetry", "swift", "terraform",
			"bundler", "composer", "devcontainers", "elm",
			"gitsubmodule", "github-actions", "mix", "pub"},
	}

	_, result, err := prompt.Run()

	if err != nil {
		panic(err)
	}

	return result
}

// TOD0: Add validataion
func getDirectory() string {
	prompt := promptui.Prompt{
		Label:   "Provide Directory (/)",
		Default: "/",
	}

	result, err := prompt.Run()

	if err != nil {
		panic(err)
	}

	return result
}

func getInterval() string {
	prompt := promptui.Select{
		Label: "Select Interval",
		Items: []string{"daily", "weekly", "monthly"},
	}

	_, result, err := prompt.Run()

	if err != nil {
		panic(err)
	}

	return result
}

// TOD0: Add validataion
func getReviewer() string {
	prompt := promptui.Prompt{
		Label: "Provide Reviewer (optional)",
	}

	result, err := prompt.Run()

	if err != nil {
		panic(err)
	}

	return result
}

func getOpenPullRequestLimit() int {
	prompt := promptui.Prompt{
		Label: "Provide Open Pull Request Limit (optional)",
	}

	result, err := prompt.Run()

	if err != nil {
		panic(err)
	}

	limit, err := strconv.Atoi(result)
	if err != nil {
		return 5
	}

	return limit
}

func addAdditionalPackageManager() string {
	prompt := promptui.Select{
		Label: "Do you want to add another package manager?",
		Items: []string{YES, NO},
	}

	_, result, err := prompt.Run()

	if err != nil {
		panic(err)
	}

	return result
}
