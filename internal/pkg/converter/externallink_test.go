package converter

import (
	"testing"

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
	converter := NewJSON(zap.NewNop())
	link := converter.convertExternalLink(externalLink)

	req.Equal("joe", link.Title)
	req.Equal("description", link.Description)
	req.Equal("cloud", link.Icon)
	req.Equal("http://lala", link.URL)
	req.True(link.OpenInNewTab)
	req.True(link.IncludeTimeRange)
	req.True(link.IncludeVariableValues)
}
