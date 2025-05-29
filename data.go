package main

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"time"
)

const (
	companyIDCol   = 0
	dateOrYearCol  = 1
	firstMetricCol = 2
)

// Load disclosure dataset where year is a column and each row is unique
func loadDisclosureData(path string) Dataset {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Cannot open file %s: %v", path, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Cannot read csv %s: %v", path, err)
	}

	header := rows[0]
	data := make(Dataset)
	for _, row := range rows[1:] {
		companyID, _ := strconv.Atoi(row[companyIDCol])
		year, _ := strconv.Atoi(row[dateOrYearCol])
		vals := make(MetricValues)
		for i, h := range header[firstMetricCol:] {
			v, err := strconv.ParseFloat(row[i+firstMetricCol], 64)
			if err == nil {
				vals[h] = v
			}
		}
		key := CompanyYear{CompanyID: companyID, Year: year}
		data[key] = vals
	}
	return data
}

// Load a dataset where the date, "you should use the latest for the year for the input data".
func loadLatestPerYear(path string, prefix string) Dataset {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Cannot open file %s: %v", path, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Cannot read csv %s: %v", path, err)
	}

	header := rows[0]
	data := make(Dataset)
	latestDate := make(map[CompanyYear]time.Time)

	for _, row := range rows[1:] {
		companyID, _ := strconv.Atoi(row[companyIDCol])
		t, _ := time.Parse("2006-01-02", row[dateOrYearCol])
		year := t.Year()
		key := CompanyYear{CompanyID: companyID, Year: year}
		if prevDate, exists := latestDate[key]; exists && !t.After(prevDate) {
			continue // skip, a later date found earlier
		}

		vals := make(MetricValues)
		for i, h := range header[firstMetricCol:] {
			v, err := strconv.ParseFloat(row[i+firstMetricCol], 64)
			if err == nil {
				vals[h] = v
				vals[prefix+"."+h] = v // Also store with prefix
			}
		}
		data[key] = vals
		latestDate[key] = t
	}
	return data
}
