package main

type CompanyYear struct {
	CompanyID int
	Year      int
}

type MetricValues map[string]float64
type Dataset map[CompanyYear]MetricValues
type ScoringContext map[string]Dataset

type MetricOperation struct {
	Type       string        `yaml:"type"`
	Parameters []MetricParam `yaml:"parameters"`
}

type MetricParam struct {
	Source string `yaml:"source"`
	Param  string `yaml:"param,omitempty"`
}

type MetricConfig struct {
	Name      string          `yaml:"name"`
	Operation MetricOperation `yaml:"operation"`
}

type ScoreConfig struct {
	Name    string         `yaml:"name"`
	Metrics []MetricConfig `yaml:"metrics"`
}
