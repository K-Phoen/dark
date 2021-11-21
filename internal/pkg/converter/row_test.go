package converter

import (
	"testing"

	"github.com/K-Phoen/sdk"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestConvertRow(t *testing.T) {
	req := require.New(t)

	converter := NewJSON(zap.NewNop())

	row := converter.convertRow(sdk.Panel{CommonPanel: sdk.CommonPanel{Title: "Row title"}})

	req.Equal("Row title", row.Name)
}

func TestConvertCollapsedRow(t *testing.T) {
	req := require.New(t)

	converter := NewJSON(zap.NewNop())

	row := converter.convertRow(sdk.Panel{
		CommonPanel: sdk.CommonPanel{
			Title: "Row title",
		},
		RowPanel: &sdk.RowPanel{
			Collapsed: true,
		},
	})

	req.Equal("Row title", row.Name)
	req.True(row.Collapse)
}
