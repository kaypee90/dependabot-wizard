package main

import (
	"flag"
	"fmt"
	"os"
)

const version = "0.4.0"

func displayAppVersion() {
	fmt.Printf("Depbot %s\n", version)
	os.Exit(0)
}

func handlePromptError(err error) {
	if err != nil {
		os.Exit(0)
	}
}

func launchApplicaton() {
	// display introduction text
	printIntroductoryText()

	var updates []Update
	existingConfigurationChecked := false

	// get dependabot configuration details
	for {
		if !existingConfigurationChecked && dependabotFileExists() {
			if result, err := getConfigurationOverrideConfirmation(); err != nil || result == NO {
				os.Exit(0)
			}

			existingConfigurationChecked = true
		}

		packageEcosystem, err := getPackageEcosystem()
		handlePromptError(err)

		directory, err := getDirectory()
		handlePromptError(err)

		interval, err := getInterval()
		handlePromptError(err)

		reviewer, err := getReviewer()
		handlePromptError(err)

		openPullRequestsLimit, err := getOpenPullRequestLimit()
		handlePromptError(err)

		updates = append(updates, Update{
			PackageEcosystem: packageEcosystem,
			Directory:        directory,
			Schedule: Schedule{
				Interval: interval,
			},
			Reviewers:             []string{reviewer},
			OpenPullRequestsLimit: openPullRequestsLimit,
		})

		if result, err := addAdditionalPackageManager(); err != nil || result == NO {
			handlePromptError(err)
			break
		}
	}

	config := Config{
		Version: 2,
		Updates: updates,
	}

	// write configration data to yaml file
	yamlData := config.ConvertToYaml()
	createDependabotYamlFile(yamlData)
	os.Exit(0)
}

func main() {
	showVersion := flag.Bool("version", false, "Display app version")
	flag.Parse()
	if *showVersion {
		displayAppVersion()
	} else {
		launchApplicaton()
	}
}
