package semver

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSemVer(t *testing.T) {
	cases := []struct {
		version  string
		expected *Version
		hasErr   bool
	}{
		{"0.1.0", &Version{Minor: 1}, false},
		{"3.0.1.0", nil, true},
		{"0.-1.0", nil, true},
		{"0.01.0", nil, true},
		{"-1.0", nil, true},
		{"0", nil, true},
	}
	for _, c := range cases {
		v := Parse(c.version)
		if c.hasErr {
			assert.Error(t, v.Error())
		} else {
			assert.NoError(t, v.Error())
			assert.Equal(t, c.expected, v)
		}
	}
}

func TestSemVer_Equal(t *testing.T) {
	cases := []struct {
		v1, v2 string
		equal  bool
	}{
		{"0.0.1", "0.0.1", true},
		{"0.0.1", "0.1.0", false},
		{"0.0.1", "1.0.0", false},
	}
	for _, c := range cases {
		v1, v2 := MustParse(c.v1), MustParse(c.v2)
		assert.Equal(t, v1.Equal(v2), c.equal)
	}
}

func TestSemVer_LessThan(t *testing.T) {
	cases := []struct {
		v1, v2   string
		lessThan bool
	}{
		{"1.0.0", "0.0.1", false},
		{"0.1.0", "0.0.1", false},
		{"0.0.1", "0.0.1", false},
		{"0.0.1", "0.0.2", true},
		{"0.0.1", "0.1.0", true},
		{"0.0.1", "1.0.0", true},
	}
	for _, c := range cases {
		v1, v2 := MustParse(c.v1), MustParse(c.v2)
		assert.Equal(t, v1.LessThan(v2), c.lessThan)
	}
}

func TestSemVer_GreaterThan(t *testing.T) {
	cases := []struct {
		v1, v2      string
		greaterThan bool
	}{
		{"1.0.0", "0.0.1", true},
		{"0.1.0", "0.0.1", true},
		{"0.0.2", "0.0.1", true},
		{"0.0.1", "0.0.1", false},
		{"0.0.1", "0.1.0", false},
		{"0.0.1", "1.0.0", false},
	}
	for _, c := range cases {
		v1, v2 := MustParse(c.v1), MustParse(c.v2)
		assert.Equal(t, v1.GreaterThan(v2), c.greaterThan)
	}
}

func TestSemVer_Bump(t *testing.T) {
	cases := []struct {
		version string
		bumped  string
		major   bool
		minor   bool
		patch   bool
	}{
		{version: "0.2.3", bumped: "1.0.0", major: true},
		{version: "0.1.9", bumped: "0.2.0", minor: true},
		{version: "0.1.9", bumped: "0.1.10", patch: true},
	}
	for _, c := range cases {
		v := Parse(c.version)
		require.NoError(t, v.Error())
		switch {
		case c.major:
			v.Bump(VersionTypeMajor)
		case c.minor:
			v.Bump(VersionTypeMinor)
		case c.patch:
			v.Bump(VersionTypePatch)
		}
		assert.Equal(t, c.bumped, v.String())
	}
}
