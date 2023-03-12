package converter

import (
	grabana "github.com/K-Phoen/grabana/decoder"
	"github.com/K-Phoen/sdk"
)

func (converter *JSON) convertDashboardLink(link sdk.Link) *grabana.DashboardInternalLink {
	extLink := &grabana.DashboardInternalLink{
		Title:                 link.Title,
		Tags:                  link.Tags,
		IncludeVariableValues: link.IncludeVars,
	}

	if link.AsDropdown != nil {
		extLink.AsDropdown = *link.AsDropdown
	}
	if link.TargetBlank != nil {
		extLink.OpenInNewTab = *link.TargetBlank
	}
	if link.KeepTime != nil {
		extLink.IncludeTimeRange = *link.KeepTime
	}

	return extLink
}
