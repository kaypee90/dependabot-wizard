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
	OpenPullRequestsLimit int      `yaml:"open-pull-requests-limit,omitempty"`
}

type Config struct {
	Version int      `yaml:"version"`
	Updates []Update `yaml:"updates"`
}

func (config *Config) ConvertToYaml() []byte {
	data, err := yaml.Marshal(config)
	if err != nil {
		panic(err)
	}

	return data
}
