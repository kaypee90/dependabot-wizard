package main

import (
	"flag"
	"fmt"
	"os"
)

const version = "0.3.1"

func displayAppVersion() {
	fmt.Printf("Depbot %s\n", version)
	os.Exit(0)
}

func launchApplicaton() {
	// display introduction text
	printIntroductoryText()

	var updates []Update
	existingConfigurationChecked := false

	// get dependabot configuration details
	for {
		if !existingConfigurationChecked && dependabotFileExists() {
			if getConfigurationOverrideConfirmation() == NO {
				os.Exit(0)
			}

			existingConfigurationChecked = true
		}

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

		if addAdditionalPackageManager() == NO {
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
