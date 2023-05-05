package main

import (
	"github.com/jingyugao/rowserrcheck/passes/rowserr"
	"github.com/ryanrolds/sqlclosecheck/pkg/analyzer"
	"golang.org/x/tools/go/analysis"
)

func getThirdPartyAnalyzers() []*analysis.Analyzer {
	return []*analysis.Analyzer{
		// checks whether sql.Rows.Err is correctly checked
		rowserr.NewAnalyzer(),
		// checks if SQL rows/statements are closed
		analyzer.NewAnalyzer(),
	}
}
