# go-semver

## Usage

### from command-line
main.go
``` go
package main

import (
	semver "github.com/ktr0731/go-semver"
)

var version = semver.New("0.1.1")

// something...
```

``` sh
$ semver -minor main.go

package main

import (
	semver "github.com/ktr0731/go-semver"
)

var version = semver.New("0.2.0")

// something...
```

### as library
``` go
package main

import (
	"fmt"

	semver "github.com/ktr0731/go-semver"
)

var v = semver.New("0.1.1")

func init() {
	if v.Error() != nil {
		panic(v.Error())
	}
}

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
