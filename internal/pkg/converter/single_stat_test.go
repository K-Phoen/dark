package converter

import (
	"reflect"
	"testing"

	"github.com/K-Phoen/sdk"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestConvertSingleStatPanel(t *testing.T) {
	req := require.New(t)

	converter := NewJSON(zap.NewNop())
	height := "200px"
	datasource := "prometheus"

	singlestatPanel := sdk.Panel{
		CommonPanel: sdk.CommonPanel{
			Title:       "Singlestat panel",
			Description: strPtr("panel desc"),
			Type:        "singlestat",
			Transparent: true,
			Height:      &height,
			Datasource:  &sdk.DatasourceRef{LegacyName: datasource},
		},
		SinglestatPanel: &sdk.SinglestatPanel{
			Format:          "none",
			ValueName:       "current",
			ValueFontSize:   "120%",
			PrefixFontSize:  strPtr("80%"),
			PostfixFontSize: strPtr("80%"),
			Colors:          []string{"blue", "red", "green"},
			ColorBackground: true,
			ColorValue:      true,
		},
	}

	converted, ok := converter.convertDataPanel(singlestatPanel)

	req.True(ok)
	req.True(converted.SingleStat.Transparent)
	req.Equal("Singlestat panel", converted.SingleStat.Title)
	req.Equal("panel desc", converted.SingleStat.Description)
	req.Equal("none", converted.SingleStat.Unit)
	req.Equal("current", converted.SingleStat.ValueType)
	req.Equal("120%", converted.SingleStat.ValueFontSize)
	req.Equal("80%", converted.SingleStat.PrefixFontSize)
	req.Equal("80%", converted.SingleStat.PostfixFontSize)
	req.Equal(height, converted.SingleStat.Height)
	req.Equal(datasource, converted.SingleStat.Datasource)
	req.True(reflect.DeepEqual(converted.SingleStat.Colors, [3]string{
		"blue", "red", "green",
	}))
	req.True(reflect.DeepEqual(converted.SingleStat.Color, []string{
		"background", "value",
	}))
}
