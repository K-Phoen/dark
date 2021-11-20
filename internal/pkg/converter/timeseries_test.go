package converter

import (
	"testing"

	"github.com/K-Phoen/sdk"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestConvertTimeSeriesPanel(t *testing.T) {
	req := require.New(t)

	converter := NewJSON(zap.NewNop())
	height := "400px"
	datasource := "prometheus"

	panel := sdk.Panel{
		CommonPanel: sdk.CommonPanel{
			OfType:      sdk.TimeseriesType,
			Title:       "test timeseries",
			Type:        "timeseries",
			Description: strPtr("timeseries description"),
			Transparent: true,
			Height:      height,
			Datasource:  &datasource,
		},
		TimeseriesPanel: &sdk.TimeseriesPanel{},
	}

	converted, ok := converter.convertDataPanel(panel)

	req.True(ok)
	req.NotNil(converted.TimeSeries)

	convertedTs := converted.TimeSeries
	req.True(convertedTs.Transparent)
	req.Equal("test timeseries", convertedTs.Title)
	req.Equal("timeseries description", convertedTs.Description)
	req.Equal(height, convertedTs.Height)
	req.Equal(datasource, convertedTs.Datasource)
}

func TestConvertTimeSeriesLegendDisplay(t *testing.T) {
	testCases := []struct {
		display  string
		expected string
	}{
		{
			display:  "list",
			expected: "as_list",
		},
		{
			display:  "table",
			expected: "as_table",
		},
		{
			display:  "hidden",
			expected: "hide",
		},
	}

	for _, testCase := range testCases {
		tc := testCase

		t.Run(tc.display, func(t *testing.T) {
			req := require.New(t)

			rawLegend := sdk.TimeseriesLegendOptions{
				DisplayMode: tc.display,
			}

			converter := NewJSON(zap.NewNop())
			legend := converter.convertTimeSeriesLegend(rawLegend)

			req.Contains(legend, tc.expected)
		})
	}
}

func TestConvertTimeSeriesLegendPlacement(t *testing.T) {
	testCases := []struct {
		placement string
		expected  string
	}{
		{
			placement: "right",
			expected:  "to_the_right",
		},
		{
			placement: "bottom",
			expected:  "to_bottom",
		},
	}

	for _, testCase := range testCases {
		tc := testCase

		t.Run(tc.placement, func(t *testing.T) {
			req := require.New(t)

			rawLegend := sdk.TimeseriesLegendOptions{
				Placement: tc.placement,
			}

			converter := NewJSON(zap.NewNop())
			legend := converter.convertTimeSeriesLegend(rawLegend)

			req.Contains(legend, tc.expected)
		})
	}
}

func TestConvertTimeSeriesLegendCalculations(t *testing.T) {
	req := require.New(t)

	rawLegend := sdk.TimeseriesLegendOptions{
		Calcs: []string{
			"first",
			"firstNotNull",
			"last",
			"lastNotNull",
			"min",
			"max",
			"mean",
			"count",
			"sum",
			"range",
		},
	}

	converter := NewJSON(zap.NewNop())
	legend := converter.convertTimeSeriesLegend(rawLegend)

	expected := []string{
		"first",
		"first_non_null",
		"last",
		"last_non_null",
		"min",
		"max",
		"avg",
		"count",
		"total",
		"range",
	}
	for _, expectedItem := range expected {
		req.Contains(legend, expectedItem)
	}
}
