package main

import (
	"golang.org/x/tools/go/analysis"
	"honnef.co/go/tools/simple"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"
)

func getStaticcheckAnalyzers() []*analysis.Analyzer {
	lstA := make([]*analysis.Analyzer, 0)

	// All staticcheck analyzers
	for _, sa := range staticcheck.Analyzers {
		lstA = append(lstA, sa.Analyzer)
	}

	// Simple analyzers
	for _, s := range simple.Analyzers {
		switch s.Analyzer.Name {
		// Simplify error construction with fmt.Errorf
		case "S1028":
		default:
			continue
		}
		lstA = append(lstA, s.Analyzer)
	}

	// Stylecheck analyzers
	for _, st := range stylecheck.Analyzers {
		switch st.Analyzer.Name {
		// Poorly chosen name for error variable
		case "ST1012":
		default:
			continue
		}
		lstA = append(lstA, st.Analyzer)
	}

	return lstA
}
