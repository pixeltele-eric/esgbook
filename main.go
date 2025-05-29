package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
)

const (
	ConfigPath        = "data/score_1.yaml"
	DisclosureCSVPath = "data/disclosure_data.csv"
	EmissionsCSVPath  = "data/emissions_data.csv"
	WasteCSVPath      = "data/waste_data.csv"
	OutputCSVPath     = "output.csv"
)

func main() {
	config := mustLoadConfig(ConfigPath)
	ctx := loadAllData()
	writeScoresCSV(OutputCSVPath, config, ctx)
}

func mustLoadConfig(path string) *ScoreConfig {
	config, err := LoadScoreConfig(path)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	return config
}

func loadAllData() ScoringContext {
	disclosure := loadDisclosureData(DisclosureCSVPath)
	emissions := loadLatestPerYear(EmissionsCSVPath, "emissions")
	waste := loadLatestPerYear(WasteCSVPath, "waste")

	ctx := ScoringContext{
		"disclosure": disclosure,
		"emissions":  emissions,
		"waste":      waste,
	}
	return ctx
}

func writeScoresCSV(filename string, config *ScoreConfig, ctx ScoringContext) {
	outFile, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Failed to create output file: %v", err)
	}
	defer outFile.Close()
	w := csv.NewWriter(outFile)
	defer w.Flush()

	// Build and write header
	header := buildHeader(config)
	if err := w.Write(header); err != nil {
		log.Fatalf("Failed to write header: %v", err)
	}

	// Write each data row
	for k := range ctx["disclosure"] {
		row := buildRow(k, config, ctx)
		if err := w.Write(row); err != nil {
			log.Printf("Warning: could not write row for company %d year %d", k.CompanyID, k.Year)
		}
	}
	fmt.Printf("Output written to %s\n", filename)
}

func buildHeader(config *ScoreConfig) []string {
	header := []string{"company_id", "year"}
	for _, m := range config.Metrics {
		header = append(header, m.Name)
	}
	return header
}

func buildRow(k CompanyYear, config *ScoreConfig, ctx ScoringContext) []string {
	companyID, year := k.CompanyID, k.Year
	computed := computeMetrics(companyID, year, config, ctx)
	row := []string{strconv.Itoa(companyID), strconv.Itoa(year)}
	for _, m := range config.Metrics {
		v, ok := computed[m.Name]
		if ok {
			row = append(row, fmt.Sprintf("%.2f", v))
		} else {
			row = append(row, "")
		}
	}
	return row
}
