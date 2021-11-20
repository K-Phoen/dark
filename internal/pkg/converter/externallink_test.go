package converter

import (
	"testing"

	grabana "github.com/K-Phoen/grabana/decoder"
	"github.com/K-Phoen/sdk"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestConvertExternalLink(t *testing.T) {
	req := require.New(t)

	externalLink := sdk.Link{
		Title:       "joe",
		Type:        "link",
		Icon:        strPtr("cloud"),
		IncludeVars: true,
		KeepTime:    boolPtr(true),
		TargetBlank: boolPtr(true),
		Tooltip:     strPtr("description"),
		URL:         strPtr("http://lala"),
	}
	dashLink := sdk.Link{
		Title: "not link",
	}

	converter := NewJSON(zap.NewNop())

	dashboard := &grabana.DashboardModel{}
	converter.convertExternalLinks([]sdk.Link{externalLink, dashLink}, dashboard)

	req.Len(dashboard.ExternalLinks, 1)

	link := dashboard.ExternalLinks[0]

	req.Equal("joe", link.Title)
	req.Equal("description", link.Description)
	req.Equal("cloud", link.Icon)
	req.Equal("http://lala", link.URL)
	req.True(link.OpenInNewTab)
	req.True(link.IncludeTimeRange)
	req.True(link.IncludeVariableValues)
}
