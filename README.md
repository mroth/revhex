# revhex ⏪

[![PkgGoDev](https://pkg.go.dev/badge/github.com/mroth/revhex)](https://pkg.go.dev/github.com/mroth/revhex)

Go module that implements *reverse hexadecimal* encoding and decoding, as used
in the [Jujutsu version control system](https://jj-vcs.github.io/jj/).

The reverse hex format is normal 16-base hexadecimal encoding but using `z-k` as
the alphabet instead of `0-9a-f`.

```go
import (
	"encoding/hex"
	"fmt"

	"github.com/mroth/revhex"
)

func Example() {
	src := []byte("Hello!")
	fmt.Println("hex:    ", hex.EncodeToString(src))
	fmt.Println("revhex: ", revhex.EncodeToString(src))

	// Output:
	// hex:     48656c6c6f21
	// revhex:  vrtutntntkxy
}
```

The implementations in this package are modeled after the standard library
`encoding/hex` package, and have identical APIs and performance characteristics.
