package main

import (
	"bytes"
	"go/parser"
	"go/token"
	"strings"
	"testing"
)

func Test_realMain(t *testing.T) {
	cases := map[string]struct {
		src      string
		show     bool
		bumpType bumpType

		hasErr bool
		assert func(t *testing.T, out *bytes.Buffer)
	}{
		"package go-semver is not imported": {
			src:    `package main`,
			hasErr: true,
		},
		"both of Parse and MustParse are not found": {
			src:    `package main; import semver "github.com/ktr0731/go-semver"`,
			hasErr: true,
		},
		"number of MustParse is invalid": {
			src:    `package main; import semver "github.com/ktr0731/go-semver"; var v = semver.MustParse("foo", "bar")`,
			hasErr: true,
		},
		"show enabled realMain writes the current version": {
			src:      `package main; import semver "github.com/ktr0731/go-semver"; var v = semver.MustParse("0.1.2")`,
			show:     true,
			bumpType: bumpTypeNoop,
			assert: func(t *testing.T, out *bytes.Buffer) {
				if !strings.Contains(out.String(), "0.1.2") {
					t.Errorf("the result of realMain must contain the current version 0.1.2, but it is not found")
				}
			},
		},
		"patch": {
			src:      `package main; import semver "github.com/ktr0731/go-semver"; var v = semver.MustParse("0.1.2")`,
			bumpType: bumpTypePatch,
			assert: func(t *testing.T, out *bytes.Buffer) {
				if !strings.Contains(out.String(), "0.1.3") {
					t.Errorf("the result of realMain must contain a new version 0.1.3, but it is not found")
				}
			},
		},
		"minor": {
			src:      `package main; import semver "github.com/ktr0731/go-semver"; var v = semver.MustParse("0.1.2")`,
			bumpType: bumpTypeMinor,
			assert: func(t *testing.T, out *bytes.Buffer) {
				if !strings.Contains(out.String(), "0.2.0") {
					t.Errorf("the result of realMain must contain a new version 0.2.0, but it is not found")
				}
			},
		},
		"major": {
			src:      `package main; import semver "github.com/ktr0731/go-semver"; var v = semver.MustParse("0.1.2")`,
			bumpType: bumpTypeMajor,
			assert: func(t *testing.T, out *bytes.Buffer) {
				if !strings.Contains(out.String(), "1.0.0") {
					t.Errorf("the result of realMain must contain a new version 1.0.0, but it is not found")
				}
			},
		},
		"semver.MustParse uses a const": {
			src:      `package main; import semver "github.com/ktr0731/go-semver"; const v = "0.1.2"; var ver = semver.MustParse(v)`,
			show:     true,
			bumpType: bumpTypeNoop,
			assert: func(t *testing.T, out *bytes.Buffer) {
				if !strings.Contains(out.String(), "0.1.2") {
					t.Errorf("the result of realMain must contain a new version 1.0.0, but it is not found")
				}
			},
		},
		"semver.MustParse uses a var": {
			src:      `package main; import semver "github.com/ktr0731/go-semver"; var (v = "0.1.2"; ver = semver.MustParse(v))`,
			show:     true,
			bumpType: bumpTypeNoop,
			assert: func(t *testing.T, out *bytes.Buffer) {
				if !strings.Contains(out.String(), "0.1.2") {
					t.Errorf("the result of realMain must contain a new version 1.0.0, but it is not found")
				}
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, "", c.src, parser.Mode(0))
			if err != nil {
				t.Fatalf("failed to parse file: %s", err)
			}
			var out bytes.Buffer
			err = realMain(c.show, fset, f, c.bumpType, &out)
			if c.hasErr {
				if err == nil {
					t.Fatal("realMain must return an error, but got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("realMain must not return an error, but got %s", err)
				}
				c.assert(t, &out)
			}
		})
	}
}
