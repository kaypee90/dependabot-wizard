package main

import (
	"os"
	"testing"

	"gotest.tools/v3/assert"
)

const testDir = "./test"

func TestCreateDirectoryIfItDoesNotExistIfDirectoryAlreadyDoesntExist(t *testing.T) {
	defer os.RemoveAll(testDir)

	hasCreatedNewDir, err := createDirectoryIfItDoesNotExist(testDir)

	assert.NilError(t, err)
	assert.Equal(t, hasCreatedNewDir, true)
}

func TestCreateDirectoryIfItDoesNotExistIfDirectoryAlreadyExist(t *testing.T) {
	defer os.RemoveAll(testDir)

	_ = os.Mkdir(testDir, 0755)

	hasCreatedNewDir, err := createDirectoryIfItDoesNotExist(testDir)

	assert.NilError(t, err)
	assert.Equal(t, hasCreatedNewDir, false)
}

func TestCreateConfigurationFile(t *testing.T) {
	const fileName = "dependabot.yml"
	const destinationDir = ".github"

	defer os.RemoveAll(destinationDir)

	_ = os.Mkdir(destinationDir, 0755)

	data := []byte("version: 2\nupdates:\n  - package-ecosystem: \"npm\"")
	filePath := getDependabotYamlFilePath(fileName, destinationDir)

	createConfigurationFile(filePath, destinationDir, data, false)

	_, err := os.Stat(destinationDir + "/" + fileName)

	// Exists
	assert.NilError(t, err)
}

func TestGetWorkingDirectory(t *testing.T) {
	dirName := getWorkingDirectory()

	assert.Equal(t, dirName, "depbot")
}

func TestFileExistsWhenFileIsNotCreated(t *testing.T) {
	defer os.RemoveAll(testDir)

	hasFile := fileExists(testDir + "/dependabot.yml")

	assert.Equal(t, hasFile, false)
}

func TestFileExistsWhenFileAlrieadyExists(t *testing.T) {
	filePath := testDir + "/dependabot.yml"

	// Create the file
	file, err := os.Create(filePath)
	if err != nil {
		return
	}
	defer file.Close()
	defer os.RemoveAll(testDir)

	hasFile := fileExists(filePath)

	assert.Equal(t, hasFile, true)
}

func TestGetCliHelpText(t *testing.T) {
	text := getCliHelpText()

	assert.Equal(t, text, cliHelpText)
}
