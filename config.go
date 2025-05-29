package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

func LoadScoreConfig(path string) (*ScoreConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var config ScoreConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
