package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const (
	dependabotFileName = "dependabot.yml"
	githubDirectory    = ".github"
)

func printIntroductoryText() {
	green := "\033[32m"
	reset := "\033[0m"
	fmt.Printf("%sDepbot wizard will assist you to configure dependabot in your project%s\n", green, reset)
}

func writeBytesToFile(fileName string, data []byte) {
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

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func createDirectoryIfItDoesNotExist(dir string) (hasCreatedNewDir bool, err error) {
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

func createConfigurationFile(fullFilePath string, destinationDir string, data []byte, skipCreatingDir bool) {
	if !skipCreatingDir {
		createDirectoryIfItDoesNotExist(destinationDir)
	}

	writeBytesToFile(fullFilePath, data)
}

func getWorkingDirectory() string {
	currentDirectory, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return filepath.Base(currentDirectory)
}

func checkIfCurrentDirectoryIsGithub() bool {
	return getWorkingDirectory() == githubDirectory
}

func generateDependabotYamlFilePath(fileName string, destinationDir string) string {
	fullFilePath := fileName

	if !checkIfCurrentDirectoryIsGithub() {
		fullFilePath = destinationDir + "/" + fileName
	}

	return fullFilePath
}

func createDependabotYamlFile(data []byte) {
	skipCreatingDir := checkIfCurrentDirectoryIsGithub()
	fullFilePath := generateDependabotYamlFilePath(dependabotFileName, githubDirectory)
	createConfigurationFile(fullFilePath, githubDirectory, data, skipCreatingDir)
}

func checkIfDependabotFileExists() bool {
	dependabotYmlFilePath := generateDependabotYamlFilePath(dependabotFileName, githubDirectory)
	return fileExists(dependabotYmlFilePath)
}
