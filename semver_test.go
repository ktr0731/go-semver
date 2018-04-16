package semver

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
