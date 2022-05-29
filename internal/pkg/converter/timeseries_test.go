package converter

import (
	"fmt"
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
			Datasource:  &sdk.DatasourceRef{LegacyName: datasource},
		},
		TimeseriesPanel: &sdk.TimeseriesPanel{
			Targets: []sdk.Target{
				{
					Expr:         "prometheus_query",
					LegendFormat: "{{ field }}",
					RefID:        "A",
				},
			},
		},
	}

	converted, ok := converter.convertDataPanel(panel)

	req.True(ok)
	req.NotNil(converted.TimeSeries)

	convertedTS := converted.TimeSeries
	req.True(convertedTS.Transparent)
	req.Equal("test timeseries", convertedTS.Title)
	req.Equal("timeseries description", convertedTS.Description)
	req.Equal(height, convertedTS.Height)
	req.Equal(datasource, convertedTS.Datasource)
	req.Len(convertedTS.Targets, 1)
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

func TestConvertTimeSeriesVisualizationGradient(t *testing.T) {
	testCases := []struct {
		mode     string
		expected string
	}{
		{
			mode:     "none",
			expected: "none",
		},
		{
			mode:     "hue",
			expected: "hue",
		},
		{
			mode:     "opacity",
			expected: "opacity",
		},
		{
			mode:     "scheme",
			expected: "scheme",
		},
	}

	for _, testCase := range testCases {
		tc := testCase

		t.Run(tc.mode, func(t *testing.T) {
			req := require.New(t)

			panel := sdk.Panel{
				CommonPanel: sdk.CommonPanel{},
				TimeseriesPanel: &sdk.TimeseriesPanel{
					Options: sdk.TimeseriesOptions{},
					FieldConfig: sdk.FieldConfig{
						Defaults: sdk.FieldConfigDefaults{
							Custom: sdk.FieldConfigCustom{
								GradientMode: tc.mode,
							},
						},
					},
				},
			}

			converter := NewJSON(zap.NewNop())
			tsViz := converter.convertTimeSeriesVisualization(panel)

			req.Equal(tc.expected, tsViz.GradientMode)
		})
	}
}

func TestConvertTimeSeriesVisualizationTooltipMode(t *testing.T) {
	testCases := []struct {
		mode     string
		expected string
	}{
		{
			mode:     "none",
			expected: "none",
		},
		{
			mode:     "single",
			expected: "single_series",
		},
		{
			mode:     "multi",
			expected: "all_series",
		},
	}

	for _, testCase := range testCases {
		tc := testCase

		t.Run(tc.mode, func(t *testing.T) {
			req := require.New(t)

			panel := sdk.Panel{
				CommonPanel: sdk.CommonPanel{},
				TimeseriesPanel: &sdk.TimeseriesPanel{
					Options: sdk.TimeseriesOptions{
						Tooltip: sdk.TimeseriesTooltipOptions{
							Mode: tc.mode,
						},
					},
					FieldConfig: sdk.FieldConfig{},
				},
			}

			converter := NewJSON(zap.NewNop())
			tsViz := converter.convertTimeSeriesVisualization(panel)

			req.Equal(tc.expected, tsViz.Tooltip)
		})
	}
}

func TestConvertTimeSeriesAxisPlacement(t *testing.T) {
	testCases := []struct {
		placement string
		expected  string
	}{
		{placement: "hidden", expected: "hidden"},
		{placement: "left", expected: "left"},
		{placement: "right", expected: "right"},
		{placement: "auto", expected: "auto"},
	}

	for _, testCase := range testCases {
		tc := testCase

		t.Run(tc.placement, func(t *testing.T) {
			req := require.New(t)

			panel := sdk.Panel{
				CommonPanel: sdk.CommonPanel{},
				TimeseriesPanel: &sdk.TimeseriesPanel{
					FieldConfig: sdk.FieldConfig{
						Defaults: sdk.FieldConfigDefaults{
							Custom: sdk.FieldConfigCustom{
								AxisPlacement: tc.placement,
							},
						},
					},
				},
			}

			converter := NewJSON(zap.NewNop())
			tsAxis := converter.convertTimeSeriesAxis(panel)

			req.Equal(tc.expected, tsAxis.Display)
		})
	}
}
func TestConvertTimeSeriesAxisOptions(t *testing.T) {
	req := require.New(t)

	panel := sdk.Panel{
		CommonPanel: sdk.CommonPanel{},
		TimeseriesPanel: &sdk.TimeseriesPanel{
			FieldConfig: sdk.FieldConfig{
				Defaults: sdk.FieldConfigDefaults{
					Unit:     "short",
					Decimals: intPtr(2),
					Min:      float64Ptr(1),
					Max:      float64Ptr(11),
					Custom: sdk.FieldConfigCustom{
						AxisLabel:   "label",
						AxisSoftMin: intPtr(0),
						AxisSoftMax: intPtr(10),
					},
				},
			},
		},
	}

	converter := NewJSON(zap.NewNop())
	tsAxis := converter.convertTimeSeriesAxis(panel)

	req.Equal("label", tsAxis.Label)
	req.Equal("short", tsAxis.Unit)
	req.Equal(2, *tsAxis.Decimals)
	req.Equal(float64(1), *tsAxis.Min)
	req.Equal(float64(11), *tsAxis.Max)
	req.Equal(0, *tsAxis.SoftMin)
	req.Equal(10, *tsAxis.SoftMax)
}

func TestConvertTimeSeriesAxisScale(t *testing.T) {
	testCases := []struct {
		scaleMode struct {
			Type string `json:"type"`
			Log  int    `json:"log,omitempty"`
		}
		expected string
	}{
		{
			scaleMode: struct {
				Type string `json:"type"`
				Log  int    `json:"log,omitempty"`
			}{
				Type: "linear",
			},
			expected: "linear",
		},
		{
			scaleMode: struct {
				Type string `json:"type"`
				Log  int    `json:"log,omitempty"`
			}{
				Type: "log",
				Log:  2,
			},
			expected: "log2",
		},
		{
			scaleMode: struct {
				Type string `json:"type"`
				Log  int    `json:"log,omitempty"`
			}{
				Type: "log",
				Log:  10,
			},
			expected: "log10",
		},
	}

	for _, testCase := range testCases {
		tc := testCase

		t.Run(fmt.Sprintf("%s %d", tc.scaleMode.Type, tc.scaleMode.Log), func(t *testing.T) {
			req := require.New(t)

			panel := sdk.Panel{
				CommonPanel: sdk.CommonPanel{},
				TimeseriesPanel: &sdk.TimeseriesPanel{
					FieldConfig: sdk.FieldConfig{
						Defaults: sdk.FieldConfigDefaults{
							Custom: sdk.FieldConfigCustom{
								ScaleDistribution: tc.scaleMode,
							},
						},
					},
				},
			}

			converter := NewJSON(zap.NewNop())
			tsAxis := converter.convertTimeSeriesAxis(panel)

			req.Equal(tc.expected, tsAxis.Scale)
		})
	}
}
