package converter

import (
	"testing"

	"github.com/K-Phoen/sdk"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestConvertDashboardLink(t *testing.T) {
	req := require.New(t)

	externalLink := sdk.Link{
		Title:       "joe",
		Tags:        []string{"my-service"},
		Type:        "dashboards",
		IncludeVars: true,
		AsDropdown:  boolPtr(true),
		KeepTime:    boolPtr(true),
		TargetBlank: boolPtr(true),
	}
	converter := NewJSON(zap.NewNop())
	link := converter.convertDashboardLink(externalLink)

	req.Equal("joe", link.Title)
	req.ElementsMatch([]string{"my-service"}, link.Tags)
	req.True(link.OpenInNewTab)
	req.True(link.IncludeTimeRange)
	req.True(link.AsDropdown)
	req.True(link.IncludeVariableValues)
}
