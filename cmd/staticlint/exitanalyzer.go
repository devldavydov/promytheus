package main

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
)

var ExitAnalyzer = &analysis.Analyzer{
	Name: "exitanalyzer",
	Doc:  `check for os.Exit in package main "main" function`,
	Run:  runExitAnalyzer,
}

func runExitAnalyzer(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		parseFile(pass, file)
	}
	return nil, nil
}

func parseFile(pass *analysis.Pass, n *ast.File) {
	// if not "main" file - skip
	if n.Name.Name != "main" {
		return
	}

	// check all func declarations of "main" file
	for _, decl := range n.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		// not main func - skip
		if funcDecl.Name.Name != "main" {
			continue
		}

		// find in main func
		if pos, ok := findInMainFunc(funcDecl); ok {
			pass.Reportf(pos, `cannot use "os.Exit" in package main "main" function`)
			break
		}
	}
}

func findInMainFunc(mainFNode *ast.FuncDecl) (token.Pos, bool) {
	var pos token.Pos
	var res bool
	ast.Inspect(mainFNode, func(node ast.Node) bool {
		fSelExpr, ok := node.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		if fSelExpr.X.(*ast.Ident).Name == "os" && fSelExpr.Sel.Name == "Exit" {
			pos, res = fSelExpr.Pos(), true
			return false
		}
		return true
	})

	return pos, res
}
