package main

import (
	"fmt"
	"log"
	"os"
)

const (
	dependabotFileName = "dependabot.yml"
	githubDir          = ".github"
)

func displayIntroductoryText() {
	green := "\033[32m"
	reset := "\033[0m"
	fmt.Printf("%sDepbot wizard will help you configure dependabot in your project%s\n", green, reset)
}

func createDirIfItDoesNotExit(dir string) (hasCreatedNewDir bool, err error) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.Mkdir(dir, 0755)

		if err != nil {
			log.Fatalf("Error creating directory %s: %v", dir, err)
			return false, err
		}

		return true, nil
	}

	return false, nil
}

func writeDataToFile(fileName string, data []byte) {
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("Error creating file %s: %v", fileName, err)
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		log.Fatalf("Error writing to file to %s: %v", fileName, err)
	}

	log.Printf("Dependabot config created successfully!")
}

func createConfigurationFile(fileName string, destinationDir string, data []byte) {
	createDirIfItDoesNotExit(destinationDir)
	fullFilePath := destinationDir + "/" + fileName
	writeDataToFile(fullFilePath, data)
}

func createDependabotYmlFile(data []byte) {
	createConfigurationFile(dependabotFileName, githubDir, data)
}
