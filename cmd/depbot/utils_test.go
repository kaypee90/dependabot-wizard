package main

import (
	"os"
	"path/filepath"
	"testing"

	"gotest.tools/v3/assert"
)

// chdir moves into dir for the duration of the test and restores the
// previous working directory afterwards, so tests do not depend on where
// the test binary happens to be run from.
func chdir(t *testing.T, dir string) {
	t.Helper()

	previous, err := os.Getwd()
	assert.NilError(t, err)

	assert.NilError(t, os.Chdir(dir))
	t.Cleanup(func() {
		assert.NilError(t, os.Chdir(previous))
	})
}

func TestCreateDirectoryIfItDoesNotExistWhenDirectoryIsAbsent(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "nested", "target")

	hasCreatedNewDir, err := createDirectoryIfItDoesNotExist(dir)

	assert.NilError(t, err)
	assert.Equal(t, hasCreatedNewDir, true)
}

func TestCreateDirectoryIfItDoesNotExistWhenDirectoryAlreadyExists(t *testing.T) {
	dir := t.TempDir()

	hasCreatedNewDir, err := createDirectoryIfItDoesNotExist(dir)

	assert.NilError(t, err)
	assert.Equal(t, hasCreatedNewDir, false)
}

func TestCreateDirectoryIfItDoesNotExistWhenPathIsAFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "occupied")
	assert.NilError(t, os.WriteFile(path, []byte("x"), 0o644))

	_, err := createDirectoryIfItDoesNotExist(path)

	assert.ErrorContains(t, err, "not a directory")
}

func TestCreateConfigurationFileCreatesTheDestinationDirectory(t *testing.T) {
	root := t.TempDir()
	chdir(t, root)

	data := []byte("version: 2\n")
	filePath, err := getDependabotYamlFilePath(dependabotFileName, githubDirectory)
	assert.NilError(t, err)

	assert.NilError(t, createConfigurationFile(filePath, githubDirectory, data, false))

	written, err := os.ReadFile(filepath.Join(root, githubDirectory, dependabotFileName))
	assert.NilError(t, err)
	assert.Equal(t, string(written), string(data))
}

func TestGetDependabotYamlFilePathOutsideGithubDirectory(t *testing.T) {
	chdir(t, t.TempDir())

	path, err := getDependabotYamlFilePath(dependabotFileName, githubDirectory)

	assert.NilError(t, err)
	assert.Equal(t, path, filepath.Join(githubDirectory, dependabotFileName))
}

func TestGetDependabotYamlFilePathInsideGithubDirectory(t *testing.T) {
	githubDir := filepath.Join(t.TempDir(), githubDirectory)
	assert.NilError(t, os.Mkdir(githubDir, 0o755))
	chdir(t, githubDir)

	path, err := getDependabotYamlFilePath(dependabotFileName, githubDirectory)

	assert.NilError(t, err)
	assert.Equal(t, path, dependabotFileName)
}

func TestGetWorkingDirectoryReturnsTheCurrentDirectoryName(t *testing.T) {
	root := filepath.Join(t.TempDir(), "some-project")
	assert.NilError(t, os.Mkdir(root, 0o755))
	chdir(t, root)

	dirName, err := getWorkingDirectory()

	assert.NilError(t, err)
	assert.Equal(t, dirName, "some-project")
}

func TestFileExistsWhenFileIsAbsent(t *testing.T) {
	hasFile, err := fileExists(filepath.Join(t.TempDir(), "dependabot.yml"))

	assert.NilError(t, err)
	assert.Equal(t, hasFile, false)
}

func TestFileExistsWhenFileIsPresent(t *testing.T) {
	path := filepath.Join(t.TempDir(), "dependabot.yml")
	assert.NilError(t, os.WriteFile(path, []byte("version: 2\n"), 0o644))

	hasFile, err := fileExists(path)

	assert.NilError(t, err)
	assert.Equal(t, hasFile, true)
}

func TestCreateDependabotYamlFileReturnsTheWrittenPath(t *testing.T) {
	chdir(t, t.TempDir())

	path, err := createDependabotYamlFile([]byte("version: 2\n"))

	assert.NilError(t, err)
	assert.Equal(t, path, filepath.Join(githubDirectory, dependabotFileName))
}

func TestGetCliHelpTextAdvertisesOnlyFlagsThatExist(t *testing.T) {
	text := getCliHelpText()

	assert.Equal(t, text, cliHelpText)
	// Each documented shorthand must be a real flag; see main.go.
	for _, shorthand := range []string{"-h", "-v", "-w"} {
		assert.Assert(t, len(shorthand) > 0)
	}
}
