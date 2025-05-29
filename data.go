package main

import (
	"encoding/csv"
	"os"
	"strconv"
	"time"
)

func loadDisclosureData(path string) Dataset {
	file, _ := os.Open(path)
	defer file.Close()
	reader := csv.NewReader(file)
	rows, _ := reader.ReadAll()
	header := rows[0]
	data := make(Dataset)
	for _, row := range rows[1:] {
		company_id, _ := strconv.Atoi(row[0])
		year, _ := strconv.Atoi(row[1])
		values := make(MetricValues)
		for i, h := range header[2:] {
			v, err := strconv.ParseFloat(row[i+2], 64)
			if err == nil {
				values[h] = v
			}
		}
		key := CompanyYear{company_id, year}
		data[key] = values
	}
	return data
}

func loadLatestPerYear(path string, prefix string) Dataset {
	file, _ := os.Open(path)
	defer file.Close()
	reader := csv.NewReader(file)
	rows, _ := reader.ReadAll()
	header := rows[0]
	idxDate := 1
	data := make(Dataset)
	dates := make(map[CompanyYear]string)
	for _, row := range rows[1:] {
		cid, _ := strconv.Atoi(row[0])
		t, _ := time.Parse("2006-01-02", row[idxDate])
		year := t.Year()
		key := CompanyYear{cid, year}
		if d, ok := dates[key]; ok && d > row[idxDate] {
			continue // already have later date
		}
		vals := make(MetricValues)
		for i, h := range header[2:] {
			v, err := strconv.ParseFloat(row[i+2], 64)
			if err == nil {
				vals[h] = v
				vals[prefix+"."+h] = v // also store with prefix
			}
		}
		data[key] = vals
		dates[key] = row[idxDate]
	}
	return data
}
