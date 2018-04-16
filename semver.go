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

func New(in string) *Version {
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
