package main

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v2"
)

// supportedEcosystems is the single source of truth for the package
// ecosystems depbot offers and accepts. The CLI prompt, the API and
// validation all read from here so they cannot drift apart.
var supportedEcosystems = []string{
	"bundler", "cargo", "composer", "devcontainers", "docker", "elm",
	"github-actions", "gitsubmodule", "gomod", "gradle", "maven", "mix",
	"npm", "nuget", "pip", "poetry", "pub", "swift", "terraform",
}

// supportedIntervals are the schedule intervals dependabot accepts.
var supportedIntervals = []string{"daily", "weekly", "monthly"}

type Schedule struct {
	Interval string `yaml:"interval"`
}

type Update struct {
	PackageEcosystem      string   `yaml:"package-ecosystem"`
	Directory             string   `yaml:"directory"`
	Schedule              Schedule `yaml:"schedule"`
	Reviewers             []string `yaml:"reviewers,omitempty"`
	OpenPullRequestsLimit uint64   `yaml:"open-pull-requests-limit,omitempty"`
	Labels                []string `yaml:"labels,omitempty"`
}

type Config struct {
	Version uint8    `yaml:"version"`
	Updates []Update `yaml:"updates"`
}

func contains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}

	return false
}

// Validate reports the first reason this config would be rejected by
// dependabot. It runs before every write so a malformed config never
// reaches disk.
func (config *Config) Validate() error {
	if config.Version != 2 {
		return fmt.Errorf("version must be 2, got %d", config.Version)
	}

	if len(config.Updates) == 0 {
		return fmt.Errorf("at least one update block is required")
	}

	for index, update := range config.Updates {
		position := fmt.Sprintf("updates[%d]", index)

		if !contains(supportedEcosystems, update.PackageEcosystem) {
			return fmt.Errorf("%s: unsupported package-ecosystem %q, expected one of %s",
				position, update.PackageEcosystem, strings.Join(supportedEcosystems, ", "))
		}

		if !strings.HasPrefix(update.Directory, "/") {
			return fmt.Errorf("%s: directory must start with %q, got %q",
				position, "/", update.Directory)
		}

		if !contains(supportedIntervals, update.Schedule.Interval) {
			return fmt.Errorf("%s: unsupported schedule interval %q, expected one of %s",
				position, update.Schedule.Interval, strings.Join(supportedIntervals, ", "))
		}

		for _, reviewer := range update.Reviewers {
			if strings.TrimSpace(reviewer) == "" {
				return fmt.Errorf("%s: reviewers must not contain empty entries", position)
			}
		}

		for _, label := range update.Labels {
			if strings.TrimSpace(label) == "" {
				return fmt.Errorf("%s: labels must not contain empty entries", position)
			}
		}
	}

	return nil
}

// ConvertToYaml validates the config and renders it, so callers cannot
// serialize something dependabot would reject.
func (config *Config) ConvertToYaml() ([]byte, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return nil, fmt.Errorf("rendering yaml: %w", err)
	}

	return data, nil
}

// ParseConfig reads a dependabot.yml back into a Config and validates it.
// This is what makes the write path verifiable: anything we generate must
// survive a round trip through here.
func ParseConfig(data []byte) (*Config, error) {
	var config Config
	if err := yaml.UnmarshalStrict(data, &config); err != nil {
		return nil, fmt.Errorf("parsing yaml: %w", err)
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return &config, nil
}
