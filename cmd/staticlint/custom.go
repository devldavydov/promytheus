package main

import (
	"github.com/devldavydov/promytheus/pkg/exitanalyzer"
	"golang.org/x/tools/go/analysis"
)

var customAnalyzers = []*analysis.Analyzer{
	exitanalyzer.ExitAnalyzer,
}
