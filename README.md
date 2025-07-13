# revhex ⏪

[![PkgGoDev](https://pkg.go.dev/badge/github.com/mroth/revhex)](https://pkg.go.dev/github.com/mroth/revhex)

Go module that implements *reverse hexadecimal* encoding and decoding, as used
in the [Jujutsu version control system](https://jj-vcs.github.io/jj/latest/).

The reverse hex format is normal 16-base hexadecimal encoding but using `z-k` as
the alphabet instead of `0-9a-f`.

```go
func Example() {
	sample := []byte{187, 132, 192, 163, 222, 186, 197, 248}
	fmt.Println("hex:    ", hex.EncodeToString(sample))
	fmt.Println("revhex: ", revhex.EncodeToString(sample))

	// Output:
	// hex:     bb84c0a3debac5f8
	// revhex:  oorvnzpwmlopnukr
}
```

The implementations in this package are modeled after the standard library
`encoding/hex` package, and have identical APIs and performance characteristics.
