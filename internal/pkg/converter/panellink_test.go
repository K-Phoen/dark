package converter

import (
	"testing"

	"github.com/K-Phoen/sdk"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestConvertPanelLink(t *testing.T) {
	req := require.New(t)

	sdkLink := sdk.Link{
		Title:       "joe",
		TargetBlank: boolPtr(true),
		URL:         strPtr("http://lala"),
	}

	converter := NewJSON(zap.NewNop())

	convertedLinks := converter.convertPanelLinks([]sdk.Link{sdkLink})

	req.Len(convertedLinks, 1)

	link := convertedLinks[0]

	req.Equal("joe", link.Title)
	req.Equal("http://lala", link.URL)
	req.True(link.OpenInNewTab)
}
