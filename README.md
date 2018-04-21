# go-semver
[![GoDoc](https://godoc.org/github.com/ktr0731/go-updater?status.svg)](https://godoc.org/github.com/ktr0731/go-updater)
[![CircleCI](https://circleci.com/gh/ktr0731/go-semver.svg?style=svg)](https://circleci.com/gh/ktr0731/go-semver)  

## Usage

### from command-line
main.go
``` go
package main

import (
	semver "github.com/ktr0731/go-semver"
)

var version = semver.MustParse("0.1.1")

// something...
```

if you want to write the result directly, use `-w` option.  
``` sh
$ bump minor main.go

package main

import (
	semver "github.com/ktr0731/go-semver"
)

var version = semver.MustParse("0.2.0")

// something...
```

### as library
``` go
package main

import (
	"fmt"

	semver "github.com/ktr0731/go-semver"
)

var v = semver.MustParse("0.1.1")

func main() {
	fmt.Printf("[%s] major: %d, minor: %d, patch: %d\n", v, v.Major, v.Minor, v.Patch)

	// bump up
	v.Bump(semver.VersionTypeMinor)
	fmt.Printf("[%s] major: %d, minor: %d, patch: %d\n", v, v.Major, v.Minor, v.Patch)
}
```

``` sh
$ go run main.go

[0.1.1] major: 0, minor: 1, patch: 1
[0.2.0] major: 0, minor: 2, patch: 0
```
