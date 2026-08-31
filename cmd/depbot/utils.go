package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

const (
	dependabotFileName = "dependabot.yml"
	githubDirectory    = ".github"
	cliHelpText        = `Depbot - A CLI tool to help you configure dependabot in your project

  Usage:
  depbot [flags]

  Flags:
  --help     -h  Show help
  --version  -v  Display app version
  --web      -w  Start web application
  --port         Port for the web application (default 3001)
  --report       Print a dependabot health report for a repository
  --owner        Repository owner, with --report
  --repo         Repository name, with --report
  --label        Only count pull requests with this label, with --report
  --author       Pull request author to count, with --report (default "app/dependabot")

  Environment:
  DEPBOT_NO_BROWSER  Set to any value to keep --web from opening a browser
  GITHUB_TOKEN       Token with repo scope, required by --report
`
)

func getCliHelpText() string {
	return cliHelpText
}

// openBrowser asks the operating system to open url in the default browser.
// A failure here is not fatal: the caller can still print the URL.
func openBrowser(url string) error {
	switch runtime.GOOS {
	case "linux":
		return exec.Command("xdg-open", url).Start()
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		return exec.Command("open", url).Start()
	default:
		return fmt.Errorf("unsupported platform %q", runtime.GOOS)
	}
}

func printIntroductoryText() {
	green := "\033[32m"
	reset := "\033[0m"
	fmt.Printf("%sDepbot wizard will assist you to configure dependabot in your project%s\n", green, reset)
}

func writeBytesToFile(fileName string, data []byte) error {
	if err := os.WriteFile(fileName, data, 0o644); err != nil {
		return fmt.Errorf("writing %s: %w", fileName, err)
	}

	return nil
}

// createDirectoryIfItDoesNotExist reports whether it had to create dir.
func createDirectoryIfItDoesNotExist(dir string) (hasCreatedNewDir bool, err error) {
	info, err := os.Stat(dir)
	if err == nil {
		if !info.IsDir() {
			return false, fmt.Errorf("%s exists but is not a directory", dir)
		}
		return false, nil
	}

	if !errors.Is(err, fs.ErrNotExist) {
		return false, fmt.Errorf("inspecting %s: %w", dir, err)
	}

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return false, fmt.Errorf("creating %s: %w", dir, err)
	}

	return true, nil
}

// fileExists distinguishes "absent" from "could not tell", so a permissions
// error is never reported as a missing file.
func fileExists(filename string) (bool, error) {
	if _, err := os.Stat(filename); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return false, nil
		}
		return false, fmt.Errorf("inspecting %s: %w", filename, err)
	}

	return true, nil
}

func createConfigurationFile(fullFilePath string, destinationDir string, data []byte, skipCreatingDir bool) error {
	if !skipCreatingDir {
		if _, err := createDirectoryIfItDoesNotExist(destinationDir); err != nil {
			return err
		}
	}

	return writeBytesToFile(fullFilePath, data)
}

func getWorkingDirectory() (string, error) {
	currentDirectory, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("resolving working directory: %w", err)
	}

	return filepath.Base(currentDirectory), nil
}

func workingDirectoryIsGithub() (bool, error) {
	dir, err := getWorkingDirectory()
	if err != nil {
		return false, err
	}

	return dir == githubDirectory, nil
}

// getDependabotYamlFilePath writes alongside the caller when they are already
// inside .github, and into .github/ otherwise.
func getDependabotYamlFilePath(fileName string, destinationDir string) (string, error) {
	inGithubDir, err := workingDirectoryIsGithub()
	if err != nil {
		return "", err
	}

	if inGithubDir {
		return fileName, nil
	}

	return filepath.Join(destinationDir, fileName), nil
}

// createDependabotYamlFile writes data and returns the path it wrote to.
func createDependabotYamlFile(data []byte) (string, error) {
	skipCreatingDir, err := workingDirectoryIsGithub()
	if err != nil {
		return "", err
	}

	fullFilePath, err := getDependabotYamlFilePath(dependabotFileName, githubDirectory)
	if err != nil {
		return "", err
	}

	if err := createConfigurationFile(fullFilePath, githubDirectory, data, skipCreatingDir); err != nil {
		return "", err
	}

	return fullFilePath, nil
}

func dependabotFileExists() (bool, error) {
	dependabotYmlFilePath, err := getDependabotYamlFilePath(dependabotFileName, githubDirectory)
	if err != nil {
		return false, err
	}

	return fileExists(dependabotYmlFilePath)
}
