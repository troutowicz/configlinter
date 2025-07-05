package configlinter

import (
	"encoding/json"

	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"
)

// Settings for the configlinter plugin
type Settings struct {
	// No specific settings needed for now, but this allows for future configuration
}

// plugin implements the register.LinterPlugin interface
type plugin struct{}

// New creates a new instance of the configlinter plugin
func New(settings any) (register.LinterPlugin, error) {
	var s Settings
	if settings != nil {
		if data, err := json.Marshal(settings); err != nil {
			return nil, err
		} else if err := json.Unmarshal(data, &s); err != nil {
			return nil, err
		}
	}
	return &plugin{}, nil
}

// BuildAnalyzers returns the analyzers for the configlinter plugin
func (p *plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{Analyzer}, nil
}

// GetLoadMode returns the load mode for the analyzer
func (p *plugin) GetLoadMode() string {
	return register.LoadModeSyntax
}

func init() {
	register.Plugin("configlinter", New)
}
