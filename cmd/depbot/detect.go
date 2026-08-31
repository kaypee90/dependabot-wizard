package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// DetectedEcosystem is one package manager found at one location in the repo.
type DetectedEcosystem struct {
	Ecosystem string `json:"ecosystem"`
	Directory string `json:"directory"`
	// Evidence is the manifest file that identified the ecosystem, so the
	// user can tell why we proposed it.
	Evidence string `json:"evidence"`
}

// skippedDirectories are never worth walking: they hold dependencies rather
// than manifests, and descending into them produces false positives.
var skippedDirectories = map[string]bool{
	".git":         true,
	"node_modules": true,
	"vendor":       true,
	"dist":         true,
	"build":        true,
	"target":       true,
	".venv":        true,
	"venv":         true,
	"__pycache__":  true,
	".next":        true,
	".terraform":   true,
}

// exactManifests maps a file name to the ecosystem it proves.
var exactManifests = map[string]string{
	"go.mod":            "gomod",
	"package.json":      "npm",
	"requirements.txt":  "pip",
	"Pipfile":           "pip",
	"setup.py":          "pip",
	"pyproject.toml":    "pip",
	"poetry.lock":       "poetry",
	"Gemfile":           "bundler",
	"composer.json":     "composer",
	"Cargo.toml":        "cargo",
	"pom.xml":           "maven",
	"build.gradle":      "gradle",
	"build.gradle.kts":  "gradle",
	"Dockerfile":        "docker",
	"mix.exs":           "mix",
	"pubspec.yaml":      "pub",
	"elm.json":          "elm",
	"Package.swift":     "swift",
	".gitmodules":       "gitsubmodule",
	"devcontainer.json": "devcontainers",
}

// suffixManifests maps a file extension to the ecosystem it proves.
var suffixManifests = map[string]string{
	".csproj": "nuget",
	".fsproj": "nuget",
	".vbproj": "nuget",
	".sln":    "nuget",
	".tf":     "terraform",
}

// ecosystemForFile returns the ecosystem a file name implies, if any.
func ecosystemForFile(name string) (string, bool) {
	if ecosystem, ok := exactManifests[name]; ok {
		return ecosystem, true
	}

	for suffix, ecosystem := range suffixManifests {
		if strings.HasSuffix(name, suffix) {
			return ecosystem, true
		}
	}

	return "", false
}

// toDependabotDirectory converts a filesystem path relative to the repo root
// into the "/"-rooted, forward-slashed form dependabot.yml expects.
func toDependabotDirectory(relativeDir string) string {
	if relativeDir == "." || relativeDir == "" {
		return "/"
	}

	return "/" + filepath.ToSlash(relativeDir)
}

// DetectEcosystems walks root and reports every package ecosystem it can
// identify, deduplicated by ecosystem and directory and sorted for stable
// output. Unreadable subtrees are skipped rather than failing the whole walk.
func DetectEcosystems(root string) ([]DetectedEcosystem, error) {
	seen := make(map[string]DetectedEcosystem)

	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			// A directory we cannot read tells us nothing; keep walking.
			if entry != nil && entry.IsDir() {
				return fs.SkipDir
			}
			return nil
		}

		if entry.IsDir() {
			if path != root && skippedDirectories[entry.Name()] {
				return fs.SkipDir
			}
			return nil
		}

		relativePath, err := filepath.Rel(root, path)
		if err != nil {
			return nil
		}

		ecosystem, ok := detectFor(relativePath, entry.Name())
		if !ok {
			return nil
		}

		directory := toDependabotDirectory(filepath.Dir(relativePath))
		// GitHub Actions workflows are always declared against the repo root,
		// never against .github/workflows itself.
		if ecosystem == "github-actions" {
			directory = "/"
		}

		key := ecosystem + "\x00" + directory
		if _, exists := seen[key]; !exists {
			seen[key] = DetectedEcosystem{
				Ecosystem: ecosystem,
				Directory: directory,
				Evidence:  filepath.ToSlash(relativePath),
			}
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("scanning %s: %w", root, err)
	}

	detected := make([]DetectedEcosystem, 0, len(seen))
	for _, item := range seen {
		detected = append(detected, item)
	}

	sort.Slice(detected, func(i, j int) bool {
		if detected[i].Ecosystem != detected[j].Ecosystem {
			return detected[i].Ecosystem < detected[j].Ecosystem
		}
		return detected[i].Directory < detected[j].Directory
	})

	return detected, nil
}

// detectFor resolves the ecosystem for one file, handling the cases that
// depend on where the file sits rather than only what it is called.
func detectFor(relativePath string, name string) (string, bool) {
	slashPath := filepath.ToSlash(relativePath)

	// Workflow files only count as github-actions inside .github/workflows.
	if strings.HasPrefix(slashPath, ".github/workflows/") {
		if strings.HasSuffix(name, ".yml") || strings.HasSuffix(name, ".yaml") {
			return "github-actions", true
		}
		return "", false
	}

	// A Dockerfile may carry a suffix, e.g. Dockerfile.ci.
	if strings.HasPrefix(name, "Dockerfile") {
		return "docker", true
	}

	return ecosystemForFile(name)
}

// DetectEcosystemsInWorkingDirectory scans the directory depbot was run from.
// When invoked inside .github it scans the repository root instead, so the
// paths it reports match the config it is about to write.
func DetectEcosystemsInWorkingDirectory() ([]DetectedEcosystem, error) {
	root := "."

	inGithubDir, err := workingDirectoryIsGithub()
	if err != nil {
		return nil, err
	}
	if inGithubDir {
		root = ".."
	}

	absoluteRoot, err := filepath.Abs(root)
	if err != nil {
		return nil, fmt.Errorf("resolving scan root: %w", err)
	}

	if _, err := os.Stat(absoluteRoot); err != nil {
		return nil, fmt.Errorf("scan root %s: %w", absoluteRoot, err)
	}

	return DetectEcosystems(absoluteRoot)
}

// UpdatesFromDetection turns detection results into dependabot update blocks.
func UpdatesFromDetection(detected []DetectedEcosystem, interval string, openPullRequestsLimit uint64) []Update {
	updates := make([]Update, 0, len(detected))

	for _, item := range detected {
		updates = append(updates, Update{
			PackageEcosystem: item.Ecosystem,
			Directory:        item.Directory,
			Schedule: Schedule{
				Interval: interval,
			},
			OpenPullRequestsLimit: openPullRequestsLimit,
		})
	}

	return updates
}
