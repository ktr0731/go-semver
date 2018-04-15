package semver

import "testing"

func TestSemVer(t *testing.T) {
	cases := []struct {
		version string
	}{
		{"0.1.0"},
	}
	for _, c := range cases {
		ver := New(c.version)
	}
}
