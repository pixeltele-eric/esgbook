package main

import (
	"math"
	"strings"
)

func sumOp(params []string, ctx ScoringContext, companyID, year int, computed map[string]float64) *float64 {
	var sum float64
	found := false
	for _, src := range params {
		v := getValue(ctx, companyID, year, src, computed)
		if v != nil {
			sum += *v
			found = true
		}
	}
	if found {
		return &sum
	}
	return nil
}

func orOp(params []string, ctx ScoringContext, companyID, year int, computed map[string]float64) *float64 {
	for _, src := range params {
		if v := getValue(ctx, companyID, year, src, computed); v != nil {
			return v
		}
	}
	return nil
}

func divideOp(params []string, ctx ScoringContext, companyID, year int, computed map[string]float64) *float64 {
	if len(params) < 2 {
		return nil
	}
	num := getValue(ctx, companyID, year, params[0], computed)
	denom := getValue(ctx, companyID, year, params[1], computed)
	if num != nil && denom != nil && *denom != 0 {
		res := *num / *denom
		return &res
	}
	return nil
}

func computeMetrics(companyID, year int, config *ScoreConfig, ctx ScoringContext) map[string]float64 {
	computed := map[string]float64{}
	for _, metric := range config.Metrics {
		opType := metric.Operation.Type
		var paramSources []string
		for _, p := range metric.Operation.Parameters {
			paramSources = append(paramSources, p.Source)
		}

		var v *float64
		switch opType {
		case "sum":
			v = sumOp(paramSources, ctx, companyID, year, computed)
		case "or":
			v = orOp(paramSources, ctx, companyID, year, computed)
		case "divide":
			v = divideOp(paramSources, ctx, companyID, year, computed)
		}
		if v != nil && !math.IsNaN(*v) {
			computed[metric.Name] = *v
		}
	}
	return computed
}

func getValue(
	ctx ScoringContext,
	companyID, year int,
	source string,
	computed map[string]float64,
) *float64 {
	if strings.HasPrefix(source, "self.") {
		name := strings.TrimPrefix(source, "self.")
		if v, ok := computed[name]; ok {
			return &v
		}
		return nil
	}
	parts := strings.Split(source, ".")
	if len(parts) != 2 {
		return nil
	}
	dataset, metric := parts[0], parts[1]
	key := CompanyYear{CompanyID: companyID, Year: year}
	if datasetMap, ok := ctx[dataset]; ok {
		if vals, ok := datasetMap[key]; ok {
			if v, ok := vals[metric]; ok {
				return &v
			}
			if v, ok := vals[dataset+"."+metric]; ok {
				return &v
			}
		}
	}
	return nil
}
