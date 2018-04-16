package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"

	semver "github.com/ktr0731/go-semver"
)

var (
	pkg = `"github.com/ktr0731/go-semver"`

	w = flag.Bool("write", false, "write to source")

	major = flag.Bool("major", false, "bump major")
	minor = flag.Bool("minor", false, "bump minor")
	patch = flag.Bool("patch", false, "bump patch")

	version = semver.New("0.1.0")
)

func main() {
	flag.Parse()

	args := flag.Args()

	if len(args) != 1 {
		return
	}

	if !*major && !*minor && !*patch {
		fmt.Println("usage: bump -major")
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

	var found bool
	ast.Inspect(f, func(n ast.Node) bool {
		if found {
			return false
		}

		expr, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		selExpr, ok := expr.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		// ast.Print(fset, expr)

		ident, ok := selExpr.X.(*ast.Ident)
		if !ok || ident.Name != is.Name.Name {
			return true
		}

		// semver expr

		if len(expr.Args) != 1 {
			fatal("number of semver.New args must be one")
		}

		// is string?
		lit, ok := expr.Args[0].(*ast.BasicLit)
		if !ok {
			// TODO: 変数を解釈する
			fatal("arg of semver.New must be string literal, passed arg is not BasicLit")
		}

		if lit.Kind != token.STRING {
			fatal("arg of semver.New must be string literal")
		}

		val := lit.Value[1 : len(lit.Value)-1]
		ver := semver.New(val)
		if ver.Error() != nil {
			fatal(fmt.Sprintf("%s", ver.Error()))
		}

		bump(ver)
		lit.Value = fmt.Sprintf(`"%s"`, ver.String())

		found = true

		return false
	})

	err = printer.Fprint(os.Stdout, fset, f)
	if err != nil {
		panic(err)
	}
}

func bump(v *semver.Version) {
	switch {
	case *major:
		v.Bump(semver.VersionTypeMajor)
	case *minor:
		v.Bump(semver.VersionTypeMinor)
	case *patch:
		v.Bump(semver.VersionTypePatch)
	}
}

func fatal(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}
