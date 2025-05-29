package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sort"
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
	dataCtx, companyYearKeys := mustLoadAllData()
	fmt.Printf("Number of company-year pairs in ALL datasets: %d\n", len(companyYearKeys))
	writeScoresCSV(OutputCSVPath, config, dataCtx, companyYearKeys)
}

func mustLoadConfig(path string) *ScoreConfig {
	config, err := LoadScoreConfig(path)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	return config
}

func mustLoadAllData() (ScoringContext, map[CompanyYear]bool) {
	disclosure := loadDisclosureData(DisclosureCSVPath)
	emissions := loadLatestPerYear(EmissionsCSVPath, "emissions")
	waste := loadLatestPerYear(WasteCSVPath, "waste")

	allKeys := map[CompanyYear]bool{}
	// Only include company-years present in all three datasets
	for k := range disclosure {
		if _, okE := emissions[k]; okE {
			if _, okW := waste[k]; okW {
				allKeys[k] = true
			}
		}
	}

	ctx := ScoringContext{
		"disclosure": disclosure,
		"emissions":  emissions,
		"waste":      waste,
	}
	return ctx, allKeys
}

// Sorts keys by company_id, then year, and writes output CSV.
func writeScoresCSV(filename string, config *ScoreConfig, ctx ScoringContext, allKeys map[CompanyYear]bool) {
	// Extract and sort keys
	var sortedKeys []CompanyYear
	for k := range allKeys {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Slice(sortedKeys, func(i, j int) bool {
		if sortedKeys[i].CompanyID == sortedKeys[j].CompanyID {
			return sortedKeys[i].Year < sortedKeys[j].Year
		}
		return sortedKeys[i].CompanyID < sortedKeys[j].CompanyID
	})

	outFile, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Failed to create output file: %v", err)
	}
	defer outFile.Close()
	w := csv.NewWriter(outFile)
	defer w.Flush()

	// Write header
	header := []string{"company_id", "year"}
	for _, m := range config.Metrics {
		header = append(header, m.Name)
	}
	if err := w.Write(header); err != nil {
		log.Fatalf("Failed to write header: %v", err)
	}

	// Write rows sorted
	for _, k := range sortedKeys {
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
		if err := w.Write(row); err != nil {
			log.Printf("Warning: could not write row for company %d year %d", companyID, year)
		}
	}
	fmt.Printf("Output written to %s (sorted)\n", filename)
}
