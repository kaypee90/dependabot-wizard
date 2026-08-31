package main

import (
	"strings"
	"testing"

	"gotest.tools/v3/assert"
)

func validConfig() Config {
	return Config{
		Version: 2,
		Updates: []Update{{
			PackageEcosystem:      "npm",
			Directory:             "/",
			Schedule:              Schedule{Interval: "weekly"},
			OpenPullRequestsLimit: 5,
		}},
	}
}

// TestConvertToYamlRoundTrips is the regression test for the bug where the
// generated file did not parse. Anything we render must read back.
func TestConvertToYamlRoundTrips(t *testing.T) {
	config := validConfig()
	config.Updates[0].Reviewers = []string{"kaypee90"}
	config.Updates[0].Labels = []string{"dependencies"}

	data, err := config.ConvertToYaml()
	assert.NilError(t, err)

	parsed, err := ParseConfig(data)
	assert.NilError(t, err)

	assert.DeepEqual(t, *parsed, config)
}

func TestConvertToYamlRendersTopLevelKeysUnindented(t *testing.T) {
	config := validConfig()
	data, err := config.ConvertToYaml()
	assert.NilError(t, err)

	lines := strings.Split(string(data), "\n")
	assert.Equal(t, lines[0], "version: 2")
	assert.Equal(t, lines[1], "updates:")
}

func TestConvertToYamlOmitsEmptyOptionalFields(t *testing.T) {
	config := validConfig()
	data, err := config.ConvertToYaml()
	assert.NilError(t, err)

	rendered := string(data)
	assert.Assert(t, !strings.Contains(rendered, "reviewers"), "empty reviewers must be omitted, got:\n%s", rendered)
	assert.Assert(t, !strings.Contains(rendered, "labels"), "empty labels must be omitted, got:\n%s", rendered)
}

func TestValidateRejectsBlankReviewer(t *testing.T) {
	config := validConfig()
	config.Updates[0].Reviewers = []string{""}

	_, err := config.ConvertToYaml()

	assert.ErrorContains(t, err, "reviewers must not contain empty entries")
}

func TestValidateRejectsUnknownEcosystem(t *testing.T) {
	config := validConfig()
	config.Updates[0].PackageEcosystem = "cocoapods"

	assert.ErrorContains(t, config.Validate(), "unsupported package-ecosystem")
}

func TestValidateRejectsUnknownInterval(t *testing.T) {
	config := validConfig()
	config.Updates[0].Schedule.Interval = "fortnightly"

	assert.ErrorContains(t, config.Validate(), "unsupported schedule interval")
}

func TestValidateRejectsRelativeDirectory(t *testing.T) {
	config := validConfig()
	config.Updates[0].Directory = "ui"

	assert.ErrorContains(t, config.Validate(), "directory must start with")
}

func TestValidateRejectsEmptyUpdates(t *testing.T) {
	config := Config{Version: 2}

	assert.ErrorContains(t, config.Validate(), "at least one update block is required")
}

// TestParseConfigRejectsTheHistoricallyBrokenOutput pins the exact shape the
// web UI used to write, so it can never be accepted again.
func TestParseConfigRejectsTheHistoricallyBrokenOutput(t *testing.T) {
	broken := "version: 2\n  updates:\n  - package-ecosystem: npm\n"

	_, err := ParseConfig([]byte(broken))

	assert.ErrorContains(t, err, "parsing yaml")
}
