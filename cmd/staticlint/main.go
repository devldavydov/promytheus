package main

import (
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
)

func main() {
	analyzers := make([]*analysis.Analyzer, 0)

	// Standard analzyers.
	analyzers = append(analyzers, standardAnalyzers...)

	// Staticcheck
	analyzers = append(analyzers, getStaticcheckAnalyzers()...)

	multichecker.Main(analyzers...)
}
