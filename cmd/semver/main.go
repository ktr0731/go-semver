package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"

	semver "github.com/ktr0731/go-semver"
)

var (
	pkg = `"github.com/ktr0731/go-semver"`

	w       = flag.Bool("write", false, "write to source")
	version = semver.New("0.1.0")
)

func main() {
	flag.Parse()

	args := flag.Args()

	if len(args) != 1 {
		return
	}

	fname := args[0]

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, fname, nil, parser.Mode(0))
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	var is *ast.ImportSpec
	for _, i := range f.Imports {
		if i.Path.Value == pkg {
			is = i
		}
	}

	if is == nil {
		fmt.Fprintf(os.Stderr, "not found")
		return
	}

	ast.Inspect(f, func(n ast.Node) bool {
		switch expr := n.(type) {
		case *ast.CallExpr:
			selExpr, ok := expr.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}

			ident, ok := selExpr.X.(*ast.Ident)
			if !ok || ident.Name != is.Name.Name {
				return true
			}

			ast.Print(fset, expr)

			// semver expr

			if len(expr.Args) != 1 {
				fatal("number of semver.New args must be one")
			}

			// is string?
			lit, ok := expr.Args[0].(*ast.BasicLit)
			if !ok {
				// TODO: 変数を解釈する
				fatal("arg of semver.New must be string literal")
			}

			if lit.Kind != token.STRING {
				fatal("arg of semver.New must be string literal")
			}

			ver := semver.New(lit.Value)
			if ver.Err() != nil {
				fatal(ver.Err().Error())
			}
		}
		return true
	})
}

func fatal(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}
