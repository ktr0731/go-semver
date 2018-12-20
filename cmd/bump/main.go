package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/ktr0731/dept/logger"
	semver "github.com/ktr0731/go-semver"
	"github.com/pkg/errors"
)

var (
	pkg     = `"github.com/ktr0731/go-semver"`
	write   = flag.Bool("w", false, "write to source")
	verbose = flag.Bool("v", false, "verbose output")

	version = semver.MustParse("0.1.0")
)

type bumpType int

const (
	bumpTypeUnknown bumpType = iota
	bumpTypeMajor
	bumpTypeMinor
	bumpTypePatch
	bumpTypeNoop // Used to show the current version.
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

	var typ bumpType
	switch args[0] {
	case "major":
		typ = bumpTypeMajor
	case "minor":
		typ = bumpTypeMinor
	case "patch":
		typ = bumpTypePatch
	case "show":
		typ = bumpTypeNoop
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
	logger.Printf("target file: %s\n", fname)

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, fname, nil, parser.Mode(0))
	if err != nil {
		fatalf("failed to parse file: %s", err)
	}

	out := os.Stdout
	if *write {
		f, err := os.Create(fname)
		if err != nil {
			fatalf("failed to write bumped source to file: %s", err)
		}
		defer f.Close()
		out = f
	}

	realMain(args[0] == "show", fset, f, typ, out)
}

func realMain(show bool, fset *token.FileSet, f *ast.File, typ bumpType, w io.Writer) error {
	logger := newLogger(*verbose)

	expr, err := findVersionExpr(fset, f)
	if err != nil {
		return err
	}

	ver, err := processExpr(expr, typ)
	if err != nil {
		return errors.Wrap(err, "failed to process expr")
	}

	logger.Printf("current version found: %s\n", ver)

	// If show command, show returned version.
	// It is equal to the current version.
	if show {
		fmt.Fprintln(w, ver)
		return nil
	}

	p := &printer.Config{
		Mode:     printer.UseSpaces | printer.TabIndent,
		Tabwidth: 8,
	}
	err = p.Fprint(w, fset, f)
	if err != nil {
		return errors.Wrap(err, "failed to print fileset")
	}
	return nil
}

// findVersionExpr finds an ast.Expr that contains version string
// such that `semver.MustParse("0.3.4")`.
// If it is not found, findVersionExpr returns nil.
func findVersionExpr(fset *token.FileSet, f *ast.File) (ast.Expr, error) {
	var is *ast.ImportSpec
	for _, i := range f.Imports {
		if i.Path.Value == pkg {
			is = i
		}
	}

	if is == nil {
		return nil, errors.Errorf("package %s not imported\n", pkg)
	}

	var targetExpr ast.Expr
	var err error
	ast.Inspect(f, func(n ast.Node) bool {
		// found
		if targetExpr != nil {
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
			err = errors.New("number of semver.Parse args must be one")
			return false
		}

		targetExpr = expr.Args[0]

		return false
	})
	if err != nil {
		return nil, err
	}
	if targetExpr == nil {
		return nil, errors.New("version not found")
	}
	return targetExpr, nil
}

func processExpr(e ast.Expr, typ bumpType) (ver *semver.Version, err error) {
	switch l := e.(type) {
	case *ast.BasicLit:
		ver, err = processBasicLit(l, typ)
	case *ast.Ident:
		ver, err = processObject(l.Obj, typ)
	default:
		panic("not supported")
	}
	return
}

func processBasicLit(l *ast.BasicLit, typ bumpType) (*semver.Version, error) {
	if l.Kind != token.STRING {
		return nil, errors.Errorf("arg of semver.Parse must be string literal, passed %T", l.Kind)
	}

	// trim double-quotes
	sv, err := strconv.Unquote(l.Value)
	if err != nil {
		return nil, errors.Errorf("failed to unquote literal: %s", err)
	}
	ver := semver.Parse(sv)
	if err := ver.Error(); err != nil {
		return nil, err
	}

	switch typ {
	case bumpTypeMajor:
		ver.Bump(semver.VersionTypeMajor)
	case bumpTypeMinor:
		ver.Bump(semver.VersionTypeMinor)
	case bumpTypePatch:
		ver.Bump(semver.VersionTypePatch)
	case bumpTypeNoop:
		// No-op
	case bumpTypeUnknown:
		panic("unknown type passed, uninitialized?")
	default:
		panic(fmt.Sprintf("unknown type: %d", typ))
	}

	l.Value = strconv.Quote(ver.String())

	return ver, nil
}

func processObject(o *ast.Object, typ bumpType) (*semver.Version, error) {
	switch s := o.Decl.(type) {
	case *ast.ValueSpec:
		if len(s.Values) != 1 {
			return nil, errors.Errorf("expect just one value, actual %d", len(s.Values))
		}
		return processExpr(s.Values[0], typ)
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

func newLogger(verbose bool) *log.Logger {
	if verbose {
		return log.New(os.Stderr, "[bump] ", log.LstdFlags|log.Lshortfile)
	}
	return log.New(ioutil.Discard, "", log.LstdFlags)
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
