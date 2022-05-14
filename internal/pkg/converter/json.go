package converter

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	v1 "github.com/K-Phoen/dark/api/v1"
	grabana "github.com/K-Phoen/grabana/decoder"
	"github.com/K-Phoen/sdk"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

type k8sDashboard struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string
	Metadata   map[string]string
	Folder     string
	Spec       *grabana.DashboardModel
}

type K8SManifestOptions struct {
	Folder    string
	Name      string
	Namespace string
}

func (options K8SManifestOptions) validate() error {
	if options.Folder == "" {
		return fmt.Errorf("folder name is required")
	}

	if options.Name == "" {
		return fmt.Errorf("dashboard name is required")
	}

	return nil
}

type JSON struct {
	logger *zap.Logger
}

func NewJSON(logger *zap.Logger) *JSON {
	return &JSON{
		logger: logger,
	}
}

func (converter *JSON) ToYAML(input io.Reader, output io.Writer) error {
	dashboard, err := converter.parseInput(input)
	if err != nil {
		converter.logger.Error("could parse input", zap.Error(err))
		return err
	}

	converted, err := yaml.Marshal(dashboard)
	if err != nil {
		converter.logger.Error("could marshall dashboard to yaml", zap.Error(err))
		return err
	}

	_, err = output.Write(converted)

	return err
}

func (converter *JSON) ToK8SManifest(input io.Reader, output io.Writer, options K8SManifestOptions) error {
	if err := options.validate(); err != nil {
		return err
	}

	dashboard, err := converter.parseInput(input)
	if err != nil {
		converter.logger.Error("could parse input", zap.Error(err))
		return err
	}

	manifest := k8sDashboard{
		APIVersion: v1.GroupVersion.String(),
		Kind:       "GrafanaDashboard",
		Metadata:   map[string]string{"name": options.Name},
		Folder:     options.Folder,
		Spec:       dashboard,
	}

	if options.Namespace != "" {
		manifest.Metadata["namespace"] = options.Namespace
	}

	converted, err := yaml.Marshal(manifest)
	if err != nil {
		converter.logger.Error("could marshall dashboard to yaml", zap.Error(err))
		return err
	}

	_, err = output.Write(converted)

	return err
}

func (converter *JSON) parseInput(input io.Reader) (*grabana.DashboardModel, error) {
	content, err := ioutil.ReadAll(input)
	if err != nil {
		converter.logger.Error("could not read input", zap.Error(err))
		return nil, err
	}

	board := &sdk.Board{}
	if err := json.Unmarshal(content, board); err != nil {
		converter.logger.Error("could not unmarshall dashboard", zap.Error(err))
		return nil, err
	}

	dashboard := &grabana.DashboardModel{}

	converter.convertGeneralSettings(board, dashboard)
	converter.convertVariables(board.Templating.List, dashboard)
	converter.convertAnnotations(board.Annotations.List, dashboard)
	converter.convertExternalLinks(board.Links, dashboard)
	converter.convertPanels(board.Panels, dashboard)

	return dashboard, nil
}

func (converter *JSON) convertGeneralSettings(board *sdk.Board, dashboard *grabana.DashboardModel) {
	dashboard.Title = board.Title
	dashboard.SharedCrosshair = board.SharedCrosshair
	dashboard.Tags = board.Tags
	dashboard.Editable = board.Editable
	dashboard.Time = [2]string{board.Time.From, board.Time.To}
	dashboard.Timezone = board.Timezone

	if board.Refresh != nil {
		dashboard.AutoRefresh = board.Refresh.Value
	}
}

func (converter *JSON) convertPanels(panels []*sdk.Panel, dashboard *grabana.DashboardModel) {
	var currentRow *grabana.DashboardRow

	for _, panel := range panels {
		if panel.Type == "row" {
			if currentRow != nil {
				dashboard.Rows = append(dashboard.Rows, *currentRow)
			}

			currentRow = converter.convertRow(*panel)

			for _, rowPanel := range panel.Panels {
				convertedPanel, ok := converter.convertDataPanel(rowPanel)
				if ok {
					currentRow.Panels = append(currentRow.Panels, convertedPanel)
				}
			}
			continue
		}

		if currentRow == nil {
			currentRow = &grabana.DashboardRow{Name: "Overview"}
		}

		convertedPanel, ok := converter.convertDataPanel(*panel)
		if ok {
			currentRow.Panels = append(currentRow.Panels, convertedPanel)
		}
	}

	if currentRow != nil {
		dashboard.Rows = append(dashboard.Rows, *currentRow)
	}
}

func (converter *JSON) convertDataPanel(panel sdk.Panel) (grabana.DashboardPanel, bool) {
	switch panel.Type {
	case "graph":
		return converter.convertGraph(panel), true
	case "heatmap":
		return converter.convertHeatmap(panel), true
	case "singlestat":
		return converter.convertSingleStat(panel), true
	case "stat":
		return converter.convertStat(panel), true
	case "table":
		return converter.convertTable(panel), true
	case "text":
		return converter.convertText(panel), true
	case "timeseries":
		return converter.convertTimeSeries(panel), true
	default:
		converter.logger.Warn("unhandled panel type: skipped", zap.String("type", panel.Type), zap.String("title", panel.Title))
	}

	return grabana.DashboardPanel{}, false
}

func panelSpan(panel sdk.Panel) float32 {
	span := panel.Span
	if span == 0 && panel.GridPos.H != nil {
		span = float32(*panel.GridPos.W / 2) // 24 units per row to 12
	}

	return span
}

func defaultOption(opt sdk.Current) string {
	if opt.Value == nil {
		return ""
	}

	return opt.Value.(string)
}
