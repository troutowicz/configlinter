package main

import (
	"golang.org/x/tools/go/analysis"

	"github.com/troutowicz/configlinter"
)

// New is the entry point for the golangci-lint Go plugin system
func New(conf any) ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{configlinter.Analyzer}, nil
}
