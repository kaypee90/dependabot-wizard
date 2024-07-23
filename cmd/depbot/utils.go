package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const (
	dependabotFileName = "dependabot.yml"
	githubDir          = ".github"
)

func displayIntroductoryText() {
	green := "\033[32m"
	reset := "\033[0m"
	fmt.Printf("%sDepbot wizard will assist you to configure dependabot in your project%s\n", green, reset)
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

func createConfigurationFile(fileName string, destinationDir string, data []byte, skipCreatingDir bool) {
	fullFilePath := fileName

	if !skipCreatingDir {
		log.Printf("Creating %s directory..", githubDir)
		createDirIfItDoesNotExit(destinationDir)
		fullFilePath = destinationDir + "/" + fileName
	} else {
		log.Printf("Skipping directory creation since %s already exists", githubDir)
	}

	writeDataToFile(fullFilePath, data)
}

func getCurrentDir() string {
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return filepath.Base(currentDir)
}

func createDependabotYmlFile(data []byte) {
	skipCreatingDir := getCurrentDir() == githubDir
	createConfigurationFile(dependabotFileName, githubDir, data, skipCreatingDir)
}
