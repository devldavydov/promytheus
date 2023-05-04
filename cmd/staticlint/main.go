package main

import (
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
)

func main() {
	analyzers := make([]*analysis.Analyzer, 0)

	analyzers = append(analyzers, standardAnalyzers...)

	multichecker.Main(analyzers...)
}
