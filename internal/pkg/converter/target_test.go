package converter

import (
	"testing"

	"github.com/K-Phoen/sdk"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestConvertTargetFailsIfNoValidTargetIsGiven(t *testing.T) {
	req := require.New(t)
	converter := NewJSON(zap.NewNop())

	convertedTarget := converter.convertTarget(sdk.Target{})
	req.Nil(convertedTarget)
}

func TestConvertTargetWithPrometheusTarget(t *testing.T) {
	req := require.New(t)

	converter := NewJSON(zap.NewNop())

	target := sdk.Target{
		Expr:         "prometheus_query",
		LegendFormat: "{{ field }}",
		RefID:        "A",
	}

	convertedTarget := converter.convertTarget(target)

	req.NotNil(convertedTarget)
	req.Nil(convertedTarget.Stackdriver)
	req.Equal("prometheus_query", convertedTarget.Prometheus.Query)
	req.Equal("{{ field }}", convertedTarget.Prometheus.Legend)
	req.Equal("A", convertedTarget.Prometheus.Ref)
}

func TestConvertTargetWithGraphiteTarget(t *testing.T) {
	req := require.New(t)

	converter := NewJSON(zap.NewNop())

	target := sdk.Target{
		Target: "graphite_query",
		RefID:  "A",
		Hide:   true,
	}

	convertedTarget := converter.convertTarget(target)

	req.NotNil(convertedTarget)
	req.NotNil(convertedTarget.Graphite)
	req.Equal("graphite_query", convertedTarget.Graphite.Query)
	req.Equal("A", convertedTarget.Graphite.Ref)
	req.True(convertedTarget.Graphite.Hidden)
}

func TestConvertTargetWithInfluxDBTarget(t *testing.T) {
	req := require.New(t)

	converter := NewJSON(zap.NewNop())

	target := sdk.Target{
		Measurement: "influxdb_query",
		RefID:       "A",
		Hide:        true,
	}

	convertedTarget := converter.convertTarget(target)

	req.NotNil(convertedTarget)
	req.NotNil(convertedTarget.InfluxDB)
	req.Equal("influxdb_query", convertedTarget.InfluxDB.Query)
	req.Equal("A", convertedTarget.InfluxDB.Ref)
	req.True(convertedTarget.InfluxDB.Hidden)
}

func TestConvertTargetWithStackdriverTargetFailsIfNoMetricKind(t *testing.T) {
	req := require.New(t)
	converter := NewJSON(zap.NewNop())

	target := sdk.Target{
		MetricType: "pubsub.googleapis.com/subscription/ack_message_count",
	}

	convertedTarget := converter.convertTarget(target)

	req.Nil(convertedTarget)
}

func TestConvertTargetWithStackdriverTargetIgnoresUnknownCrossSeriesReducer(t *testing.T) {
	req := require.New(t)
	converter := NewJSON(zap.NewNop())

	target := sdk.Target{
		MetricKind:         "DELTA",
		MetricType:         "pubsub.googleapis.com/subscription/ack_message_count",
		CrossSeriesReducer: "unknown",
	}

	convertedTarget := converter.convertTarget(target)

	req.NotNil(convertedTarget)
	req.NotNil(convertedTarget.Stackdriver)
	req.Empty(convertedTarget.Stackdriver.Aggregation)
}

func TestConvertTargetWithStackdriverTargetIgnoresUnknownAligner(t *testing.T) {
	req := require.New(t)
	converter := NewJSON(zap.NewNop())

	target := sdk.Target{
		MetricKind:       "DELTA",
		MetricType:       "pubsub.googleapis.com/subscription/ack_message_count",
		PerSeriesAligner: "unknown",
	}

	convertedTarget := converter.convertTarget(target)

	req.NotNil(convertedTarget)
	req.NotNil(convertedTarget.Stackdriver)
	req.Empty(convertedTarget.Stackdriver.Alignment)
}

func TestConvertTargetWithStackdriverTarget(t *testing.T) {
	req := require.New(t)

	converter := NewJSON(zap.NewNop())

	target := sdk.Target{
		MetricKind:         "DELTA",
		MetricType:         "pubsub.googleapis.com/subscription/ack_message_count",
		CrossSeriesReducer: "REDUCE_MEAN",
		PerSeriesAligner:   "ALIGN_DELTA",
		AlignmentPeriod:    "stackdriver-auto",
		GroupBys:           []string{"field"},
		AliasBy:            "legend",
		RefID:              "A",
		Filters: []string{
			"resource.label.subscription_id",
			"=",
			"subscription_name",
			"AND",
			"other-property",
			"!=",
			"other-value",
			"AND",
			"regex-property",
			"=~",
			"regex-value",
			"AND",
			"regex-not-property",
			"!=~",
			"regex-not-value",
		},
	}

	convertedTarget := converter.convertTarget(target)

	req.NotNil(convertedTarget)
	req.Nil(convertedTarget.Prometheus)
	req.NotNil(convertedTarget.Stackdriver)
	req.Equal("delta", convertedTarget.Stackdriver.Type)
	req.Equal("pubsub.googleapis.com/subscription/ack_message_count", convertedTarget.Stackdriver.Metric)
	req.Equal("mean", convertedTarget.Stackdriver.Aggregation)
	req.Equal("stackdriver-auto", convertedTarget.Stackdriver.Alignment.Period)
	req.Equal("delta", convertedTarget.Stackdriver.Alignment.Method)
	req.Equal("legend", convertedTarget.Stackdriver.Legend)
	req.Equal("A", convertedTarget.Stackdriver.Ref)
	req.ElementsMatch([]string{"field"}, convertedTarget.Stackdriver.GroupBy)
	req.EqualValues(map[string]string{"resource.label.subscription_id": "subscription_name"}, convertedTarget.Stackdriver.Filters.Eq)
	req.EqualValues(map[string]string{"other-property": "other-value"}, convertedTarget.Stackdriver.Filters.Neq)
	req.EqualValues(map[string]string{"regex-property": "regex-value"}, convertedTarget.Stackdriver.Filters.Matches)
	req.EqualValues(map[string]string{"regex-not-property": "regex-not-value"}, convertedTarget.Stackdriver.Filters.NotMatches)
}

func TestConvertTargetWithStackdriverGauge(t *testing.T) {
	req := require.New(t)

	converter := NewJSON(zap.NewNop())

	target := sdk.Target{
		MetricKind: "GAUGE",
		MetricType: "pubsub.googleapis.com/subscription/ack_message_count",
	}

	convertedTarget := converter.convertTarget(target)

	req.NotNil(convertedTarget)
	req.Nil(convertedTarget.Prometheus)
	req.NotNil(convertedTarget.Stackdriver)
	req.Equal("gauge", convertedTarget.Stackdriver.Type)
	req.Equal("pubsub.googleapis.com/subscription/ack_message_count", convertedTarget.Stackdriver.Metric)
}

func TestConvertTargetWithStackdriverCumulative(t *testing.T) {
	req := require.New(t)

	converter := NewJSON(zap.NewNop())

	target := sdk.Target{
		MetricKind: "CUMULATIVE",
		MetricType: "pubsub.googleapis.com/subscription/ack_message_count",
	}

	convertedTarget := converter.convertTarget(target)

	req.NotNil(convertedTarget)
	req.Nil(convertedTarget.Prometheus)
	req.NotNil(convertedTarget.Stackdriver)
	req.Equal("cumulative", convertedTarget.Stackdriver.Type)
	req.Equal("pubsub.googleapis.com/subscription/ack_message_count", convertedTarget.Stackdriver.Metric)
}
