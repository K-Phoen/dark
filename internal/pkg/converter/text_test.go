package converter

import (
	"testing"

	"github.com/K-Phoen/sdk"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestConvertTextPanelWithMarkdown(t *testing.T) {
	req := require.New(t)

	converter := NewJSON(zap.NewNop())
	height := "200px"

	textPanel := sdk.Panel{
		CommonPanel: sdk.CommonPanel{
			Title:       "Text panel",
			Transparent: true,
			Height:      &height,
			Type:        "text",
		},
		TextPanel: &sdk.TextPanel{
			Options: struct {
				Content string `json:"content"`
				Mode    string `json:"mode"`
			}{Content: "# hello world", Mode: "markdown"},
		},
	}

	converted, ok := converter.convertDataPanel(textPanel)

	req.True(ok)
	req.True(converted.Text.Transparent)
	req.Equal("Text panel", converted.Text.Title)
	req.Equal("# hello world", converted.Text.Markdown)
	req.Equal(height, converted.Text.Height)
}

func TestConvertTextLinks(t *testing.T) {
	req := require.New(t)

	converter := NewJSON(zap.NewNop())
	sdkPanel := sdk.Panel{
		CommonPanel: sdk.CommonPanel{
			Type: "text",
			Links: []sdk.Link{
				{Title: "text title", URL: strPtr("text url")},
			},
		},
		TextPanel: &sdk.TextPanel{},
	}

	converted, ok := converter.convertDataPanel(sdkPanel)

	req.True(ok)
	req.NotNil(converted.Text)

	panel := converted.Text
	req.Len(panel.Links, 1)
	req.Equal("text title", panel.Links[0].Title)
	req.Equal("text url", panel.Links[0].URL)
}

func TestConvertTextPanelWithHTML(t *testing.T) {
	req := require.New(t)

	converter := NewJSON(zap.NewNop())

	textPanel := sdk.Panel{
		CommonPanel: sdk.CommonPanel{
			Title:       "Text panel html",
			Type:        "text",
			Description: strPtr("panel description"),
		},
		TextPanel: &sdk.TextPanel{
			Mode: "html",
			Options: struct {
				Content string `json:"content"`
				Mode    string `json:"mode"`
			}{Content: "<h1>hello world</h1>", Mode: "html"},
		},
	}

	converted, ok := converter.convertDataPanel(textPanel)

	req.True(ok)
	req.False(converted.Text.Transparent)
	req.Equal("Text panel html", converted.Text.Title)
	req.Equal("panel description", converted.Text.Description)
	req.Equal("<h1>hello world</h1>", converted.Text.HTML)
}
