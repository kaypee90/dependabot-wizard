package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

const (
	dependabotFileName = "dependabot.yml"
	githubDirectory    = ".github"
	pathSeperator      = "/"
	cliHelpText        = `Depbot - A CLI tool to help you configure dependabot in your project
  Usage:
  depbot [flags]

  Flags:
  --help     -h  Show help
  --version  -v Display app version
  --web      -w    Start web application
`
)

func getCliHelpText() string {
	return cliHelpText
}

func openBrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		fmt.Printf("Failed to open browser: %v\n", err)
	}
}

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

func workingDirectoryIsGithub() bool {
	return getWorkingDirectory() == githubDirectory
}

func getDependabotYamlFilePath(fileName string, destinationDir string) string {
	fullFilePath := fileName

	if !workingDirectoryIsGithub() {
		fullFilePath = destinationDir + pathSeperator + fileName
	}

	return fullFilePath
}

func createDependabotYamlFile(data []byte) {
	skipCreatingDir := workingDirectoryIsGithub()
	fullFilePath := getDependabotYamlFilePath(dependabotFileName, githubDirectory)
	createConfigurationFile(fullFilePath, githubDirectory, data, skipCreatingDir)
}

func dependabotFileExists() bool {
	dependabotYmlFilePath := getDependabotYamlFilePath(dependabotFileName, githubDirectory)
	return fileExists(dependabotYmlFilePath)
}
