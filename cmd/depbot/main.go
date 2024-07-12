package main

func main() {
	// display introduction text
	displayIntroductoryText()

	var updates []Update

	// get dependabot configuration details
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
	createDependabotYmlFile(yamlData)
}
