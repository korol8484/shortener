package mainexit

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"honnef.co/go/tools/analysis/code"
	"strings"
)

// Analyzer - Check direct call to os.Exit in the main function
var Analyzer = &analysis.Analyzer{
	Name:     "mainexit",
	Doc:      "Checking the use of a direct call to os.Exit in the main function of the main package",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	if pass.Pkg.Name() != "main" {
		return nil, nil
	}

	fn := func(node ast.Node) {
		switch v := node.(type) {
		case *ast.FuncDecl:
			if v.Body == nil {
				return
			}

			if v.Name.Name != "main" {
				return
			}

			if strings.HasSuffix(pass.Pkg.Path(), ".example") ||
				strings.HasSuffix(pass.Pkg.Path(), ".test") ||
				strings.HasSuffix(pass.Pkg.Path(), "benchmark") ||
				strings.HasSuffix(pass.Pkg.Path(), "fuzz") {
				return
			}

			ast.Inspect(v, func(n ast.Node) bool {
				ce, okCe := n.(*ast.CallExpr)
				if !okCe {
					return true
				}

				se, okSe := ce.Fun.(*ast.SelectorExpr)
				if !okSe {
					return true
				}

				if ident, ok := se.X.(*ast.Ident); ok && ident.Name == "os" && se.Sel.Name == "Exit" {
					pass.Reportf(se.Pos(), "avoid using os.Exit directly in main function")
				}

				return false
			})
		}
	}
	//
	needle := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	code.Preorder(pass, fn, needle...)

	return nil, nil
}
