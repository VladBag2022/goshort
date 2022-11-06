package exitcheck

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "exitcheck",
	Doc:  "check for os.Exit() call from main()",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		if file.Name.String() != "main" {
			continue
		}

		ast.Inspect(file, func(node ast.Node) bool {
			// function call
			var ok bool
			var call *ast.CallExpr
			if call, ok = node.(*ast.CallExpr); !ok {
				return true
			}

			// function with name "Exit"
			var fun *ast.SelectorExpr
			if fun, ok = call.Fun.(*ast.SelectorExpr); !ok || fun.Sel.Name != "Exit" {
				return true
			}

			// from package "os"
			var id *ast.Ident
			if id, ok = fun.X.(*ast.Ident); !ok || id.Name != "os" {
				return true
			}

			pass.Reportf(call.Pos(), "os.Exit() call from main")
			return true
		})
	}
	return nil, nil
}
