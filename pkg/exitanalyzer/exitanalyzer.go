// Exit analyzer checks for use os.Exit in package main "main" function
package exitanalyzer

import (
	"go/ast"

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

		findInMainFunc(pass, funcDecl)
	}
}

func findInMainFunc(pass *analysis.Pass, mainFNode *ast.FuncDecl) {
	ast.Inspect(mainFNode, func(node ast.Node) bool {
		fSelExpr, ok := node.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		if fSelExpr.X.(*ast.Ident).Name == "os" && fSelExpr.Sel.Name == "Exit" {
			pass.Reportf(fSelExpr.Pos(), `cannot use "os.Exit" in package main "main" function`)
		}
		return true
	})
}
