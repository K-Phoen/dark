package converter

import (
	grabana "github.com/K-Phoen/grabana/decoder"
	"github.com/K-Phoen/sdk"
)

func (converter *JSON) convertRow(panel sdk.Panel) *grabana.DashboardRow {
	repeat := ""
	if panel.Repeat != nil {
		repeat = *panel.Repeat
	}
	collapse := false
	if panel.RowPanel != nil && panel.RowPanel.Collapsed {
		collapse = true
	}

	return &grabana.DashboardRow{
		Name:     panel.Title,
		Repeat:   repeat,
		Collapse: collapse,
		Panels:   nil,
	}
}
