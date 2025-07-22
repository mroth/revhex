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
	src := []byte("Hello!")
	fmt.Println("hex:    ", hex.EncodeToString(src))
	fmt.Println("revhex: ", revhex.EncodeToString(src))

	// Output:
	// hex:     48656c6c6f21
	// revhex:  vrtutntntkxy
}
