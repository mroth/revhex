package revhex_test

import (
	"encoding/hex"
	"fmt"
	"log"

	"github.com/mroth/revhex"
)

func ExampleEncode() {
	src := []byte("Hello Gopher!")

	dst := make([]byte, revhex.EncodedLen(len(src)))
	revhex.Encode(dst, src)

	fmt.Printf("%s\n", dst)

	// Output:
	// vrtutntntkxzvstksztrtusxxy
}

func ExampleDecode() {
	src := []byte("vrtutntntkxzvstksztrtusxxy")

	dst := make([]byte, revhex.DecodedLen(len(src)))
	n, err := revhex.Decode(dst, src)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", dst[:n])

	// Output:
	// Hello Gopher!
}

func ExampleDecodeString() {
	const s = "vrtutntntkxzvstksztrtusxxy"
	decoded, err := revhex.DecodeString(s)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", decoded)

	// Output:
	// Hello Gopher!
}

func ExampleEncodeToString() {
	src := []byte("Hello Gopher!")
	encodedStr := revhex.EncodeToString(src)

	fmt.Printf("%s\n", encodedStr)

	// Output:
	// vrtutntntkxzvstksztrtusxxy
}

func Example() {
	sample := []byte{187, 132, 192, 163, 222, 186, 197, 248}
	fmt.Println("bytes:  ", sample)
	fmt.Println("hex:    ", hex.EncodeToString(sample))
	fmt.Println("revhex: ", revhex.EncodeToString(sample))

	// Output:
	// bytes:   [187 132 192 163 222 186 197 248]
	// hex:     bb84c0a3debac5f8
	// revhex:  oorvnzpwmlopnukr
}
