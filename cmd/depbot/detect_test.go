package main

import (
	"os"
	"path/filepath"
	"testing"

	"gotest.tools/v3/assert"
)

func writeFile(t *testing.T, root string, relativePath string) {
	t.Helper()

	full := filepath.Join(root, filepath.FromSlash(relativePath))
	assert.NilError(t, os.MkdirAll(filepath.Dir(full), 0o755))
	assert.NilError(t, os.WriteFile(full, []byte("{}"), 0o644))
}

func TestDetectEcosystemsFindsManifestsAcrossDirectories(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "go.mod")
	writeFile(t, root, "Dockerfile")
	writeFile(t, root, "ui/app/package.json")
	writeFile(t, root, "services/api/requirements.txt")
	writeFile(t, root, ".github/workflows/ci.yml")

	detected, err := DetectEcosystems(root)
	assert.NilError(t, err)

	got := map[string]string{}
	for _, item := range detected {
		got[item.Ecosystem] = item.Directory
	}

	assert.Equal(t, got["gomod"], "/")
	assert.Equal(t, got["docker"], "/")
	assert.Equal(t, got["npm"], "/ui/app")
	assert.Equal(t, got["pip"], "/services/api")
	assert.Equal(t, got["github-actions"], "/")
}

func TestDetectEcosystemsSkipsDependencyDirectories(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "package.json")
	writeFile(t, root, "node_modules/left-pad/package.json")
	writeFile(t, root, "vendor/github.com/x/y/go.mod")

	detected, err := DetectEcosystems(root)
	assert.NilError(t, err)

	assert.Equal(t, len(detected), 1)
	assert.Equal(t, detected[0].Ecosystem, "npm")
	assert.Equal(t, detected[0].Directory, "/")
}

func TestDetectEcosystemsDeduplicatesWithinADirectory(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "package.json")
	writeFile(t, root, "Dockerfile")
	writeFile(t, root, "Dockerfile.ci")

	detected, err := DetectEcosystems(root)
	assert.NilError(t, err)

	dockerCount := 0
	for _, item := range detected {
		if item.Ecosystem == "docker" {
			dockerCount++
		}
	}
	assert.Equal(t, dockerCount, 1)
}

func TestDetectEcosystemsReturnsNothingForAnEmptyProject(t *testing.T) {
	detected, err := DetectEcosystems(t.TempDir())

	assert.NilError(t, err)
	assert.Equal(t, len(detected), 0)
}

// TestDetectionProducesAValidConfig closes the loop: whatever the scanner
// reports must survive validation and serialization.
func TestDetectionProducesAValidConfig(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "go.mod")
	writeFile(t, root, "ui/package.json")
	writeFile(t, root, ".github/workflows/ci.yml")

	detected, err := DetectEcosystems(root)
	assert.NilError(t, err)

	config := Config{Version: 2, Updates: UpdatesFromDetection(detected, "weekly", 5)}

	data, err := config.ConvertToYaml()
	assert.NilError(t, err)

	parsed, err := ParseConfig(data)
	assert.NilError(t, err)
	assert.Equal(t, len(parsed.Updates), 3)
}
