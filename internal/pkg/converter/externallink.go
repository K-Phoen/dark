package converter

import (
	grabana "github.com/K-Phoen/grabana/decoder"
	"github.com/K-Phoen/sdk"
	"go.uber.org/zap"
)

func (converter *JSON) convertExternalLink(link sdk.Link) *grabana.DashboardExternalLink {
	if link.URL == nil || *link.URL == "" {
		converter.logger.Warn("link URL empty: skipped", zap.String("title", link.Title))
		return nil
	}

	extLink := &grabana.DashboardExternalLink{
		Title:                 link.Title,
		URL:                   *link.URL,
		IncludeVariableValues: link.IncludeVars,
	}

	if link.Tooltip != nil && *link.Tooltip != "" {
		extLink.Description = *link.Tooltip
	}
	if link.Icon != nil && *link.Icon != "" {
		extLink.Icon = *link.Icon
	}
	if link.TargetBlank != nil {
		extLink.OpenInNewTab = *link.TargetBlank
	}
	if link.KeepTime != nil {
		extLink.IncludeTimeRange = *link.KeepTime
	}

	return extLink
}
