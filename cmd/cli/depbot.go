package main

const fileName = "dependabot.yml"

func main() {
	packageEcosystem := getPackageEcosystem()
	directory := getDirectory()
	interval := getInterval()
	reviewer := getReviewer()
	openPullRequestsLimit := getOpenPullRequestLimit()

	config := Config{
		Version: 2,
		Updates: []Update{
			{
				PackageEcosystem: packageEcosystem,
				Directory:        directory,
				Schedule: Schedule{
					Interval: interval,
				},
				Reviewers:             []string{reviewer},
				OpenPullRequestsLimit: openPullRequestsLimit,
			},
		},
	}

	yamlData := config.ConvertToYaml()
	writeDataFile(fileName, yamlData)
}
