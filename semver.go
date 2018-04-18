package semver

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type VersionType int

const (
	VersionTypeMajor = iota
	VersionTypeMinor
	VersionTypePatch
)

type Version struct {
	Major, Minor, Patch int
	err                 error
}

func (v *Version) Bump(typ VersionType) {
	switch typ {
	case VersionTypeMajor:
		v.Major++
		v.Minor = 0
		v.Patch = 0
	case VersionTypeMinor:
		v.Minor++
		v.Patch = 0
	case VersionTypePatch:
		v.Patch++
	}
}

func (v *Version) Compare(v2 *Version) int {
	if v.String() == v2.String() {
		return 0
	}

	// major
	if v.Major < v2.Major {
		return -1
	} else if v.Major > v2.Minor {
		return 1
	}

	// minor
	if v.Minor < v2.Minor {
		return -1
	} else if v.Minor > v2.Minor {
		return 1
	}

	// patch
	if v.Patch < v2.Patch {
		return -1
	}
	return 1
}

func (v *Version) Equal(v2 *Version) bool {
	return v.Compare(v2) == 0
}

func (v *Version) LessThan(v2 *Version) bool {
	return v.Compare(v2) == -1
}

func (v *Version) GreaterThan(v2 *Version) bool {
	return v.Compare(v2) == 1
}

func (v *Version) Error() error {
	return v.err
}

func (v *Version) major(in string) {
	if v.err != nil {
		return
	}
	v.Major, v.err = toInt(in)
}

func (v *Version) minor(in string) {
	if v.err != nil {
		return
	}
	v.Minor, v.err = toInt(in)
}

func (v *Version) patch(in string) {
	if v.err != nil {
		return
	}
	v.Patch, v.err = toInt(in)
}

func MustParse(in string) *Version {
	v := Parse(in)
	if v.Error() != nil {
		panic(v.Error())
	}
	return v
}

// Parse parses passed string as a semantic version
// if parsing will be failed, its error is stored to *Version.Error()
func Parse(in string) *Version {
	v := &Version{}

	sp := strings.Split(in, ".")
	if len(sp) != 3 {
		v.err = errors.Errorf("passed string is not following to semver: %s", in)
		return v
	}

	v.major(sp[0])
	v.minor(sp[1])
	v.patch(sp[2])

	return v
}

func (v *Version) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

func toInt(in string) (int, error) {
	n, err := strconv.ParseInt(in, 10, 32)
	if err != nil {
		return -1, errors.Wrapf(err, "failed to parse string as int: %s", err)
	}
	if n < 0 {
		return -1, errors.New("version must not negative")
	}
	// e.g. 01
	if len(in) > 1 && in[0] == '0' {
		return -1, errors.New("version must not have zero as a prefix")
	}
	return int(n), nil
}
