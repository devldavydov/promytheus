package main

import "golang.org/x/tools/go/analysis"

var customAnalyzers = []*analysis.Analyzer{
	ExitAnalyzer,
}
