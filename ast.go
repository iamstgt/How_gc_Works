package main

import (
    "go/ast"
    "go/parser"
    "go/token"
)

func main() {
    fset := token.NewFileSet()
    f, _ := parser.ParseFile(fset, "main.go", nil, parser.Mode(0))

    for _, d := range f.Decls {
        ast.Print(fset, d)
    }
}
