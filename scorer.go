package main

import (
	"math"
	"strings"
)

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

func computeMetrics(companyID, year int, config *ScoreConfig, ctx ScoringContext) map[string]float64 {
	computed := map[string]float64{}
	for _, metric := range config.Metrics {
		op := metric.Operation
		var result *float64
		switch op.Type {
		case "sum":
			sum := 0.0
			for _, param := range op.Parameters {
				v := getValue(ctx, companyID, year, param.Source, computed)
				if v != nil {
					sum += *v
				}
			}
			result = &sum
		case "or":
			x := getValue(ctx, companyID, year, op.Parameters[0].Source, computed)
			y := getValue(ctx, companyID, year, op.Parameters[1].Source, computed)
			if x != nil {
				result = x
			} else if y != nil {
				result = y
			}
		case "divide":
			num := getValue(ctx, companyID, year, op.Parameters[0].Source, computed)
			denom := getValue(ctx, companyID, year, op.Parameters[1].Source, computed)
			if num != nil && denom != nil && *denom != 0 {
				r := *num / *denom
				result = &r
			}
		}
		if result != nil && !math.IsNaN(*result) {
			computed[metric.Name] = *result
		}
	}
	return computed
}
