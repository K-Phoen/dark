package converter

import (
	grabana "github.com/K-Phoen/grabana/decoder"
	"github.com/K-Phoen/sdk"
)

func (converter *JSON) convertPanelLinks(links []sdk.Link) []grabana.DashboardPanelLink {
	converted := make([]grabana.DashboardPanelLink, 0, len(links))

	for _, link := range links {
		converted = append(converted, converter.convertPanelLink(link))
	}

	return converted
}

func (converter *JSON) convertPanelLink(link sdk.Link) grabana.DashboardPanelLink {
	return grabana.DashboardPanelLink{
		Title:        link.Title,
		URL:          *link.URL,
		OpenInNewTab: link.TargetBlank != nil && *link.TargetBlank,
	}
}
