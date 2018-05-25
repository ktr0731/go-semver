package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
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

		ast.Print(fset, expr.Args)

		// Parse or MustParse?
		if selExpr.Sel.Name != "MustParse" && selExpr.Sel.Name != "Parse" {
			return true
		}

		if len(expr.Args) != 1 {
			fatalf("number of semver.Parse args must be one")
		}

		var err error
		lit, err = extractString(expr.Args[0])
		if err != nil {
			log.Println(err)
			return true
		}

		return false
	})

	if lit == nil {
		fatalf("not found")
	}

	var ver *semver.Version
	switch l := lit.(type) {
	case *ast.BasicLit:
		ver, err = processBasicLit(l)
	default:
		panic("not supported")
	}
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

func processBasicLit(l *ast.BasicLit) (*semver.Version, error) {
	// trim double-quotes
	sv, err := strconv.Unquote(l.Value)
	if err != nil {
		return nil, errors.Errorf("failed to unquote literal: %s", err)
	}
	ver := semver.MustParse(sv)
	l.Value = strconv.Quote(ver.String())
	return ver, nil
}

func extractString(e ast.Expr) (ast.Expr, error) {
	l, ok := e.(*ast.BasicLit)
	if !ok {
		// TODO: 変数を解釈する
		return nil, errors.Errorf("arg of semver.Parse must be string literal, passed %T", e)
	}

	if l.Kind != token.STRING {
		return nil, errors.Errorf("arg of semver.Parse must be string literal, passed %T", l.Kind)
	}
	return l, nil
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
