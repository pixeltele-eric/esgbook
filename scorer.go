package main

import (
	"math"
	"strings"
)

func sumOp(params []string, ctx ScoringContext, cid, year int, computed map[string]float64) *float64 {
	s := 0.0
	ok := false
	for _, src := range params {
		v := getValue(ctx, cid, year, src, computed)
		if v != nil {
			s += *v
			ok = true
		}
	}
	if ok {
		return &s
	}
	return nil
}

func orOp(params []string, ctx ScoringContext, cid, year int, computed map[string]float64) *float64 {
	for _, src := range params {
		v := getValue(ctx, cid, year, src, computed)
		if v != nil {
			return v
		}
	}
	return nil
}

func divideOp(params []string, ctx ScoringContext, cid, year int, computed map[string]float64) *float64 {
	if len(params) < 2 {
		return nil
	}
	num := getValue(ctx, cid, year, params[0], computed)
	denom := getValue(ctx, cid, year, params[1], computed)
	if num != nil && denom != nil && *denom != 0 {
		res := *num / *denom
		return &res
	}
	return nil
}

func computeMetrics(cid, year int, config *ScoreConfig, ctx ScoringContext) map[string]float64 {
	computed := map[string]float64{}
	for _, m := range config.Metrics {
		op := m.Operation.Type
		var paramSources []string
		for _, p := range m.Operation.Parameters {
			paramSources = append(paramSources, p.Source)
		}
		var v *float64
		switch op {
		case "sum":
			v = sumOp(paramSources, ctx, cid, year, computed)
		case "or":
			v = orOp(paramSources, ctx, cid, year, computed)
		case "divide":
			v = divideOp(paramSources, ctx, cid, year, computed)
		}
		if v != nil && !math.IsNaN(*v) {
			computed[m.Name] = *v
		}
	}
	return computed
}

func getFirst(keys []string, vals MetricValues) *float64 {
	for _, k := range keys {
		if v, ok := vals[k]; ok {
			return &v
		}
	}
	return nil
}

func getValue(ctx ScoringContext, cid, year int, source string, computed map[string]float64) *float64 {
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
	key := CompanyYear{cid, year}
	if datasetMap, ok := ctx[dataset]; ok {
		if vals, ok := datasetMap[key]; ok {
			return getFirst([]string{metric, dataset + "." + metric}, vals)
		}
	}
	return nil
}
