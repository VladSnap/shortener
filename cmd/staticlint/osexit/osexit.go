// Package osexit запрещает использовать прямой вызов os.Exit в функции main пакета main.
package osexit

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// Analyzer запрещает использование os.Exit в main.main.
var Analyzer = &analysis.Analyzer{
	Name: "osexit",
	Doc:  "Запрещает использовать прямой вызов os.Exit в функции main пакета main.",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		// Проверяем, является ли файл частью пакета main
		if pass.Pkg.Name() != "main" {
			continue
		}

		ast.Inspect(file, func(n ast.Node) bool {
			callExpr, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}

			ident, ok := selExpr.X.(*ast.Ident)
			if !ok {
				return true
			}

			if ident.Name == "os" && selExpr.Sel.Name == "Exit" {
				// Теперь проверяем, находится ли это внутри функции main
				for _, decl := range file.Decls {
					if fn, ok := decl.(*ast.FuncDecl); ok && fn.Name.Name == "main" {
						pass.Reportf(callExpr.Lparen, "запрещено использовать os.Exit напрямую в main.main")
					}
				}
			}

			return true
		})
	}
	//nolint:nilnil // This is written in the examples from the theory
	return nil, nil
}
