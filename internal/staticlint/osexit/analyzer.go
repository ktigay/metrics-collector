package osexit

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

const (
	passReport = "expression returns direct call of os.Exit"
)

// Analyzer проверяет функцию main пакета main на прямой вызов os.Exit
var Analyzer = &analysis.Analyzer{
	Name: "osexit",
	Doc:  "check for direct use of os.Exit",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	expr := func(x *ast.SelectorExpr) {
		if ident, ok := x.X.(*ast.Ident); ok {
			if ident.Name != "os" || x.Sel.Name != "Exit" {
				return
			}

			pass.Reportf(x.Pos(), passReport)
		}
	}
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			switch x := node.(type) {
			case *ast.File:
				if x.Name.Name != "main" {
					return false
				}
			case *ast.FuncDecl:
				if x.Name.Name != "main" {
					return false
				}
			case *ast.SelectorExpr:
				expr(x)
			}
			return true
		})
	}

	return nil, nil
}
