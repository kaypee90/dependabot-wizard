package main

import "gopkg.in/yaml.v2"

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

func (config *Config) ConvertToYaml() []byte {
	data, err := yaml.Marshal(config)
	if err != nil {
		panic(err)
	}

	return data
}
