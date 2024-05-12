package main

import "fmt"

const fileName = "dependabot.yml"

func printIntroText() {
	green := "\033[32m"
	reset := "\033[0m"
	fmt.Printf("%sDepbot wizard will help you configure dependabot in your project%s\n", green, reset)
}

func main() {
	printIntroText()

	var updates []Update

	for {
		packageEcosystem := getPackageEcosystem()
		directory := getDirectory()
		interval := getInterval()
		reviewer := getReviewer()
		openPullRequestsLimit := getOpenPullRequestLimit()

		updates = append(updates, Update{
			PackageEcosystem: packageEcosystem,
			Directory:        directory,
			Schedule: Schedule{
				Interval: interval,
			},
			Reviewers:             []string{reviewer},
			OpenPullRequestsLimit: openPullRequestsLimit,
		})

		if NO == addAdditionalPackageManager() {
			break
		}
	}

	config := Config{
		Version: 2,
		Updates: updates,
	}

	yamlData := config.ConvertToYaml()
	writeDataFile(fileName, yamlData)
}
