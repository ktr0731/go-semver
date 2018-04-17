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
		v := New(c.version)
		if c.hasErr {
			assert.Error(t, v.Error())
		} else {
			assert.NoError(t, v.Error())
			assert.Equal(t, c.expected, v)
		}
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
		v := New(c.version)
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
