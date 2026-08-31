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

// supportedUpdateTypes are the values a group's update-types accepts.
var supportedUpdateTypes = []string{"major", "minor", "patch"}

// supportedDependencyTypes are the values a group's dependency-type accepts.
var supportedDependencyTypes = []string{"production", "development"}

// supportedAppliesTo are the values a group's applies-to accepts.
var supportedAppliesTo = []string{"version-updates", "security-updates"}

// supportedIgnoreUpdateTypes are the values an ignore rule's update-types
// accepts. Note these are prefixed, unlike a group's.
var supportedIgnoreUpdateTypes = []string{
	"version-update:semver-major",
	"version-update:semver-minor",
	"version-update:semver-patch",
}

// GroupingStrategy is a named, ready-made grouping configuration. These exist
// so the recommendation the health report makes can actually be applied:
// the report says "group these updates", and this is what produces the yaml.
type GroupingStrategy string

const (
	// GroupingNone leaves every dependency in its own pull request.
	GroupingNone GroupingStrategy = "none"

	// GroupingAll collapses minor and patch updates into one pull request,
	// leaving majors on their own because they need individual review. This
	// is the strategy that removes most of the noise.
	GroupingAll GroupingStrategy = "all"

	// GroupingByDependencyType separates production from development
	// dependencies, so a risky production bump is never buried in a batch of
	// tooling updates.
	GroupingByDependencyType GroupingStrategy = "by-dependency-type"
)

var supportedGroupingStrategies = []GroupingStrategy{
	GroupingNone, GroupingAll, GroupingByDependencyType,
}

// DescribeGrouping explains a strategy in one line, for prompts and help.
func DescribeGrouping(strategy GroupingStrategy) string {
	switch strategy {
	case GroupingAll:
		return "one pull request for all minor and patch updates, majors separate"
	case GroupingByDependencyType:
		return "separate batches for production and development dependencies"
	default:
		return "no grouping, one pull request per dependency"
	}
}

// GroupsFor renders a strategy into dependabot group definitions.
func GroupsFor(strategy GroupingStrategy) map[string]Group {
	switch strategy {
	case GroupingAll:
		return map[string]Group{
			"all-minor-and-patch": {
				Patterns:    []string{"*"},
				UpdateTypes: []string{"minor", "patch"},
			},
		}
	case GroupingByDependencyType:
		return map[string]Group{
			"production-dependencies": {
				DependencyType: "production",
				Patterns:       []string{"*"},
				UpdateTypes:    []string{"minor", "patch"},
			},
			"development-dependencies": {
				DependencyType: "development",
				Patterns:       []string{"*"},
				UpdateTypes:    []string{"minor", "patch"},
			},
		}
	default:
		return nil
	}
}

// IgnoreMajorUpdates is the rule to apply when major bumps are consistently
// closed rather than merged: stop opening them, and take them deliberately.
func IgnoreMajorUpdates() []Ignore {
	return []Ignore{{
		DependencyName: "*",
		UpdateTypes:    []string{"version-update:semver-major"},
	}}
}

type Schedule struct {
	Interval string `yaml:"interval"`
}

// Group is one entry under an update's `groups` key. Grouping is the single
// most effective reducer of dependabot pull request volume: without it every
// dependency gets its own pull request, which is what buries review queues.
type Group struct {
	AppliesTo       string   `yaml:"applies-to,omitempty"`
	DependencyType  string   `yaml:"dependency-type,omitempty"`
	Patterns        []string `yaml:"patterns,omitempty"`
	ExcludePatterns []string `yaml:"exclude-patterns,omitempty"`
	UpdateTypes     []string `yaml:"update-types,omitempty"`
}

// Ignore is one entry under an update's `ignore` key.
type Ignore struct {
	DependencyName string   `yaml:"dependency-name"`
	Versions       []string `yaml:"versions,omitempty"`
	UpdateTypes    []string `yaml:"update-types,omitempty"`
}

type Update struct {
	PackageEcosystem      string           `yaml:"package-ecosystem"`
	Directory             string           `yaml:"directory"`
	Schedule              Schedule         `yaml:"schedule"`
	Reviewers             []string         `yaml:"reviewers,omitempty"`
	OpenPullRequestsLimit uint64           `yaml:"open-pull-requests-limit,omitempty"`
	Labels                []string         `yaml:"labels,omitempty"`
	Groups                map[string]Group `yaml:"groups,omitempty"`
	Ignore                []Ignore         `yaml:"ignore,omitempty"`
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

		if err := validateGroups(position, update.Groups); err != nil {
			return err
		}

		if err := validateIgnore(position, update.Ignore); err != nil {
			return err
		}
	}

	return nil
}

func validateGroups(position string, groups map[string]Group) error {
	for name, group := range groups {
		where := fmt.Sprintf("%s.groups[%q]", position, name)

		if strings.TrimSpace(name) == "" {
			return fmt.Errorf("%s: group names must not be empty", position)
		}

		if group.DependencyType != "" && !contains(supportedDependencyTypes, group.DependencyType) {
			return fmt.Errorf("%s: unsupported dependency-type %q, expected one of %s",
				where, group.DependencyType, strings.Join(supportedDependencyTypes, ", "))
		}

		if group.AppliesTo != "" && !contains(supportedAppliesTo, group.AppliesTo) {
			return fmt.Errorf("%s: unsupported applies-to %q, expected one of %s",
				where, group.AppliesTo, strings.Join(supportedAppliesTo, ", "))
		}

		for _, updateType := range group.UpdateTypes {
			if !contains(supportedUpdateTypes, updateType) {
				return fmt.Errorf("%s: unsupported update-type %q, expected one of %s",
					where, updateType, strings.Join(supportedUpdateTypes, ", "))
			}
		}

		// A group with no selector matches nothing and silently does nothing.
		if len(group.Patterns) == 0 && group.DependencyType == "" {
			return fmt.Errorf("%s: a group needs patterns or a dependency-type, otherwise it matches nothing", where)
		}
	}

	return nil
}

func validateIgnore(position string, rules []Ignore) error {
	for index, rule := range rules {
		where := fmt.Sprintf("%s.ignore[%d]", position, index)

		if strings.TrimSpace(rule.DependencyName) == "" {
			return fmt.Errorf("%s: dependency-name is required", where)
		}

		for _, updateType := range rule.UpdateTypes {
			if !contains(supportedIgnoreUpdateTypes, updateType) {
				return fmt.Errorf("%s: unsupported update-type %q, expected one of %s",
					where, updateType, strings.Join(supportedIgnoreUpdateTypes, ", "))
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
