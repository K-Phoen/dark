package converter

import (
	grabana "github.com/K-Phoen/grabana/decoder"
	"github.com/K-Phoen/sdk"
)

func (converter *JSON) convertText(panel sdk.Panel) grabana.DashboardPanel {
	text := &grabana.DashboardText{
		Title:       panel.Title,
		Span:        panelSpan(panel),
		Transparent: panel.Transparent,
	}

	if panel.Description != nil {
		text.Description = *panel.Description
	}
	if panel.Height != nil {
		text.Height = *(panel.Height).(*string)
	}

	if panel.TextPanel.Options.Mode == "markdown" {
		text.Markdown = panel.TextPanel.Options.Content
	} else {
		text.HTML = panel.TextPanel.Options.Content
	}

	return grabana.DashboardPanel{Text: text}
}
