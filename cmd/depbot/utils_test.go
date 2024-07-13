package main

import (
	"os"
	"testing"

	"gotest.tools/v3/assert"
)

const testDir = "./test"

func TestCreateDirIfItDoesNotExitIfDirectoryAlreadyDoesntExist(t *testing.T) {
	defer os.RemoveAll(testDir)

	hasCreatedNewDir, err := createDirIfItDoesNotExit(testDir)

	assert.NilError(t, err)
	assert.Equal(t, hasCreatedNewDir, true)
}

func TestCreateDirIfItDoesNotExitIfDirectoryAlreadyExist(t *testing.T) {
	defer os.RemoveAll(testDir)

	_ = os.Mkdir(testDir, 0755)

	hasCreatedNewDir, err := createDirIfItDoesNotExit(testDir)

	assert.NilError(t, err)
	assert.Equal(t, hasCreatedNewDir, false)
}

func TestCreateConfigurationFile(t *testing.T) {
	const fileName = "dependabot.yml"

	defer os.RemoveAll(testDir)

	_ = os.Mkdir(testDir, 0755)

	data := []byte("version: 2\nupdates:\n  - package-ecosystem: \"npm\"")

	createConfigurationFile(fileName, testDir, data)

	_, err := os.Stat(testDir + "/" + fileName)

	// Exists
	assert.NilError(t, err)
}
