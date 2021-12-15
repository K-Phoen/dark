package converter

import (
	"bytes"
	"testing"

	grabana "github.com/K-Phoen/grabana/decoder"
	"github.com/K-Phoen/sdk"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestConvertInvalidJSONToYAML(t *testing.T) {
	req := require.New(t)

	converter := NewJSON(zap.NewNop())
	err := converter.ToYAML(bytes.NewBufferString(""), bytes.NewBufferString(""))

	req.Error(err)
}

func TestConvertValidJSONToYaml(t *testing.T) {
	req := require.New(t)

	converter := NewJSON(zap.NewNop())
	err := converter.ToYAML(bytes.NewBufferString("{}"), bytes.NewBufferString(""))

	req.NoError(err)
}

func TestConvertInvalidJSONToK8SManifest(t *testing.T) {
	req := require.New(t)

	converter := NewJSON(zap.NewNop())
	err := converter.ToK8SManifest(bytes.NewBufferString(""), bytes.NewBufferString(""), K8SManifestOptions{Name: "test-dashboard", Folder: "Folder"})

	req.Error(err)
}

func TestConvertValidJSONK8SManifest(t *testing.T) {
	req := require.New(t)

	converter := NewJSON(zap.NewNop())
	err := converter.ToK8SManifest(bytes.NewBufferString("{}"), bytes.NewBufferString(""), K8SManifestOptions{Name: "test-dashboard", Folder: "Folder"})

	req.NoError(err)
}

func TestConvertK8SManifestWithNoFolder(t *testing.T) {
	req := require.New(t)

	converter := NewJSON(zap.NewNop())
	err := converter.ToK8SManifest(bytes.NewBufferString("{}"), bytes.NewBufferString(""), K8SManifestOptions{Name: "name"})

	req.Error(err)
}

func TestConvertK8SManifestWithNoName(t *testing.T) {
	req := require.New(t)

	converter := NewJSON(zap.NewNop())
	err := converter.ToK8SManifest(bytes.NewBufferString("{}"), bytes.NewBufferString(""), K8SManifestOptions{Folder: "not empty"})

	req.Error(err)
}

func TestConvertGeneralSettings(t *testing.T) {
	req := require.New(t)

	board := &sdk.Board{}
	board.Title = "title"
	board.SharedCrosshair = true
	board.Editable = true
	board.Tags = []string{"tag", "other"}
	board.Refresh = &sdk.BoolString{
		Flag:  true,
		Value: "5s",
	}

	converter := NewJSON(zap.NewNop())

	dashboard := &grabana.DashboardModel{}
	converter.convertGeneralSettings(board, dashboard)

	req.Equal("title", dashboard.Title)
	req.Equal("5s", dashboard.AutoRefresh)
	req.Equal([]string{"tag", "other"}, dashboard.Tags)
	req.True(dashboard.Editable)
	req.True(dashboard.SharedCrosshair)
}

func strPtr(input string) *string {
	return &input
}

func boolPtr(input bool) *bool {
	return &input
}
