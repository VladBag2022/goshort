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
			// call
			if call, ok := node.(*ast.CallExpr); ok {
				// function
				if fun, ok := call.Fun.(*ast.SelectorExpr); ok {
					// from package "os"
					if id, ok := fun.X.(*ast.Ident); ok {
						if id.Name == "os" {
							// with name "Exit"
							if sel, ok := fun.X.(*ast.Ident); ok {
								if sel.Name == "Exit" {
									pass.Reportf(call.Pos(), "os.Exit() call from main")
								}
							}
						}
					}
				}
			}
			return true
		})
	}
	return nil, nil
}
