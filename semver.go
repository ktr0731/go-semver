package semver

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type Version struct {
	Major, Minor, Patch int
}

func New(v string) (*Version, error) {
	sp := strings.Split(v, ".")
	if len(sp) < 3 {
		return nil, errors.Errorf("passed string is not following to semver: %s", v)
	}

	const msg = "failed to parse string as int"
	major, err := strconv.ParseInt(sp[0], 10, 32)
	if err != nil {
		return nil, errors.Wrapf(err, msg)
	}
	minor, err := strconv.ParseInt(sp[1], 10, 32)
	if err != nil {
		return nil, errors.Wrapf(err, msg)
	}
	patch, err := strconv.ParseInt(sp[2], 10, 32)
	if err != nil {
		return nil, errors.Wrapf(err, msg)
	}
	return &Version{
		Major: int(major),
		Minor: int(minor),
		Patch: int(patch),
	}, nil
}

func (v *Version) String() string {
	return fmt.Sprintf("%s.%s.%s", v.Major, v.Minor, v.Patch)
}
