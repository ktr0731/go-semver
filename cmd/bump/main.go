package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"strconv"

	semver "github.com/ktr0731/go-semver"
	"github.com/pkg/errors"
)

var (
	pkg = `"github.com/ktr0731/go-semver"`

	write = flag.Bool("w", false, "write to source")

	version = semver.MustParse("0.1.0")
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), usageFormat, os.Args[0])
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	if flag.NArg() != 2 {
		exitWithUsage(1)
	}

	args := flag.Args()

	var typ semver.VersionType
	switch args[0] {
	case "major":
		typ = semver.VersionTypeMajor
	case "minor":
		typ = semver.VersionTypeMinor
	case "patch":
		typ = semver.VersionTypePatch
	case "show":
		// do nothing
	case "of":
		if args[1] == "chicken" {
			fmt.Println("見えないものを見ようとして望遠鏡を覗き込んだ")
			return
		}
		exitWithUsage(1)
	default:
		exitWithUsage(1)
	}

	fname := args[1]

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, fname, nil, parser.Mode(0))
	if err != nil {
		fatalf("failed to parse file: %s", err)
	}

	var is *ast.ImportSpec
	for _, i := range f.Imports {
		if i.Path.Value == pkg {
			is = i
		}
	}

	if is == nil {
		fatalf("package %s not imported\n", pkg)
	}

	var lit ast.Expr
	ast.Inspect(f, func(n ast.Node) bool {
		// found
		if lit != nil {
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

		ident, ok := selExpr.X.(*ast.Ident)
		if !ok || ident.Name != is.Name.Name {
			return true
		}

		// semver expr

		// Parse or MustParse?
		if selExpr.Sel.Name != "MustParse" && selExpr.Sel.Name != "Parse" {
			return true
		}

		if len(expr.Args) != 1 {
			fatalf("number of semver.Parse args must be one")
		}

		lit = expr.Args[0]

		return false
	})

	if lit == nil {
		fatalf("not found")
	}

	ver, err := processExpr(lit)
	if err != nil {
		fatalf("failed to process expr: %s", err)
	}

	// if show command, only show current version
	if args[0] == "show" {
		fmt.Println(ver)
		return
	}

	ver.Bump(typ)

	out := os.Stdout
	if *write {
		f, err := os.Create(fname)
		if err != nil {
			fatalf("failed to write bumped source to file: %s", err)
		}
		defer f.Close()
		out = f
	}

	p := &printer.Config{
		Mode:     printer.UseSpaces | printer.TabIndent,
		Tabwidth: 8,
	}
	err = p.Fprint(out, fset, f)
	if err != nil {
		fatalf("failed to print fileset: %s", err)
	}
}

func processExpr(e ast.Expr) (ver *semver.Version, err error) {
	switch l := e.(type) {
	case *ast.BasicLit:
		ver, err = processBasicLit(l)
	case *ast.Ident:
		ver, err = processObject(l.Obj)
	default:
		panic("not supported")
	}
	return
}

func processBasicLit(l *ast.BasicLit) (*semver.Version, error) {
	if l.Kind != token.STRING {
		return nil, errors.Errorf("arg of semver.Parse must be string literal, passed %T", l.Kind)
	}

	// trim double-quotes
	sv, err := strconv.Unquote(l.Value)
	if err != nil {
		return nil, errors.Errorf("failed to unquote literal: %s", err)
	}
	ver := semver.Parse(sv)
	l.Value = strconv.Quote(ver.String())
	return ver, ver.Error()
}

func processObject(o *ast.Object) (*semver.Version, error) {
	switch s := o.Decl.(type) {
	case *ast.ValueSpec:
		if len(s.Values) != 1 {
			return nil, errors.Errorf("expect just one value, actual %d", len(s.Values))
		}
		return processExpr(s.Values[0])
	default:
		return nil, errors.Errorf("unsupported type %T", s)
	}
	return nil, errors.Errorf("unsupported type %s", o.Kind)
}

func exitWithUsage(status int) {
	flag.Usage()
	os.Exit(status)
}

func fatalf(format string, a ...interface{}) {
	fmt.Fprintf(flag.CommandLine.Output(), format+"\n", a...)
	os.Exit(1)
}

const usageFormat = `
Usage: %s [-w] <command> <filename>

Commands:
	major	bump up major version
	minor	bump up minor version
	patch	bump up patch version
	show	show current version

Options:
`
