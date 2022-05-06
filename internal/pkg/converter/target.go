package converter

import (
	"strings"

	grabana "github.com/K-Phoen/grabana/decoder"
	"github.com/K-Phoen/grabana/target/stackdriver"
	"github.com/K-Phoen/sdk"
	"go.uber.org/zap"
)

func (converter *JSON) convertTarget(target sdk.Target) *grabana.Target {
	// looks like a prometheus target
	if target.Expr != "" {
		return converter.convertPrometheusTarget(target)
	}

	// looks like graphite
	if target.Target != "" {
		return converter.convertGraphiteTarget(target)
	}

	// looks like influxdb
	if target.Measurement != "" {
		return converter.convertInfluxDBTarget(target)
	}

	// looks like stackdriver
	if target.MetricType != "" {
		return converter.convertStackdriverTarget(target)
	}

	converter.logger.Warn("unhandled target type: skipped", zap.Any("target", target))

	return nil
}

func (converter *JSON) convertPrometheusTarget(target sdk.Target) *grabana.Target {
	return &grabana.Target{
		Prometheus: &grabana.PrometheusTarget{
			Query:          target.Expr,
			Legend:         target.LegendFormat,
			Hidden:         target.Hide,
			Format:         target.Format,
			Instant:        target.Instant,
			IntervalFactor: &target.IntervalFactor,
		},
	}
}

func (converter *JSON) convertGraphiteTarget(target sdk.Target) *grabana.Target {
	return &grabana.Target{
		Graphite: &grabana.GraphiteTarget{
			Query:  target.Target,
			Hidden: target.Hide,
		},
	}
}

func (converter *JSON) convertInfluxDBTarget(target sdk.Target) *grabana.Target {
	return &grabana.Target{
		InfluxDB: &grabana.InfluxDBTarget{
			Query:  target.Measurement,
			Hidden: target.Hide,
		},
	}
}

func (converter *JSON) convertStackdriverTarget(target sdk.Target) *grabana.Target {
	switch strings.ToLower(target.MetricKind) {
	case "cumulative":
	case "gauge":
	case "delta":
	default:
		converter.logger.Warn("unhandled stackdriver metric kind: target skipped", zap.Any("metricKind", target.MetricKind))
		return nil
	}

	var aggregation string
	if target.CrossSeriesReducer != "" {
		aggregationMap := map[string]string{
			string(stackdriver.ReduceNone):              "none",
			string(stackdriver.ReduceMean):              "mean",
			string(stackdriver.ReduceMin):               "min",
			string(stackdriver.ReduceMax):               "max",
			string(stackdriver.ReduceSum):               "sum",
			string(stackdriver.ReduceStdDev):            "stddev",
			string(stackdriver.ReduceCount):             "count",
			string(stackdriver.ReduceCountTrue):         "count_true",
			string(stackdriver.ReduceCountFalse):        "count_false",
			string(stackdriver.ReduceCountFractionTrue): "fraction_true",
			string(stackdriver.ReducePercentile99):      "percentile_99",
			string(stackdriver.ReducePercentile95):      "percentile_95",
			string(stackdriver.ReducePercentile50):      "percentile_50",
			string(stackdriver.ReducePercentile05):      "percentile_05",
		}

		if agg, ok := aggregationMap[target.CrossSeriesReducer]; ok {
			aggregation = agg
		} else {
			converter.logger.Warn("unhandled stackdriver crossSeriesReducer: target skipped", zap.Any("crossSeriesReducer", target.CrossSeriesReducer))
		}
	}

	var alignment *grabana.StackdriverAlignment
	if target.PerSeriesAligner != "" {
		alignmentMethodMap := map[string]string{
			string(stackdriver.AlignNone):          "none",
			string(stackdriver.AlignDelta):         "delta",
			string(stackdriver.AlignRate):          "rate",
			string(stackdriver.AlignInterpolate):   "interpolate",
			string(stackdriver.AlignNextOlder):     "next_older",
			string(stackdriver.AlignMin):           "min",
			string(stackdriver.AlignMax):           "max",
			string(stackdriver.AlignMean):          "mean",
			string(stackdriver.AlignCount):         "count",
			string(stackdriver.AlignSum):           "sum",
			string(stackdriver.AlignStdDev):        "stddev",
			string(stackdriver.AlignCountTrue):     "count_true",
			string(stackdriver.AlignCountFalse):    "count_false",
			string(stackdriver.AlignFractionTrue):  "fraction_true",
			string(stackdriver.AlignPercentile99):  "percentile_99",
			string(stackdriver.AlignPercentile95):  "percentile_95",
			string(stackdriver.AlignPercentile50):  "percentile_50",
			string(stackdriver.AlignPercentile05):  "percentile_05",
			string(stackdriver.AlignPercentChange): "percent_change",
		}

		if method, ok := alignmentMethodMap[target.PerSeriesAligner]; ok {
			alignment = &grabana.StackdriverAlignment{
				Period: target.AlignmentPeriod,
				Method: method,
			}
		} else {
			converter.logger.Warn("unhandled stackdriver perSeriesAligner: target skipped", zap.Any("perSeriesAligner", target.PerSeriesAligner))
		}
	}

	return &grabana.Target{
		Stackdriver: &grabana.StackdriverTarget{
			Project:     target.ProjectName,
			Type:        strings.ToLower(target.MetricKind),
			Metric:      target.MetricType,
			Filters:     converter.convertStackdriverFilters(target),
			Aggregation: aggregation,
			Alignment:   alignment,
			GroupBy:     target.GroupBys,
			Legend:      target.AliasBy,
			Hidden:      target.Hide,
		},
	}
}

func (converter *JSON) convertStackdriverFilters(target sdk.Target) grabana.StackdriverFilters {
	filters := grabana.StackdriverFilters{
		Eq:         map[string]string{},
		Neq:        map[string]string{},
		Matches:    map[string]string{},
		NotMatches: map[string]string{},
	}

	var leftOperand, rightOperand, operator *string
	for i := range target.Filters {
		if target.Filters[i] == "AND" {
			continue
		}

		if leftOperand == nil {
			leftOperand = &target.Filters[i]
			continue
		}
		if operator == nil {
			operator = &target.Filters[i]
			continue
		}
		if rightOperand == nil {
			rightOperand = &target.Filters[i]
		}

		if leftOperand != nil && operator != nil && rightOperand != nil {
			switch *operator {
			case "=":
				filters.Eq[*leftOperand] = *rightOperand
			case "!=":
				filters.Neq[*leftOperand] = *rightOperand
			case "=~":
				filters.Matches[*leftOperand] = *rightOperand
			case "!=~":
				filters.NotMatches[*leftOperand] = *rightOperand
			default:
				converter.logger.Warn("unhandled stackdriver filter operator: filter skipped", zap.Any("operator", *operator))
			}

			leftOperand = nil
			rightOperand = nil
			operator = nil
		}
	}

	return filters
}
