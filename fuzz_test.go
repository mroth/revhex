package revhex

import (
	"bytes"
	"testing"
)

func FuzzEncodeDecode(f *testing.F) {
	f.Add([]byte{})
	f.Add([]byte("The quick brown fox jumps over the lazy dog"))
	f.Add([]byte("abcdefghijklmnopqrstuvwxyz0123456789"))
	f.Add([]byte("Hello world!"))
	f.Add([]byte("안녕 세상아!"))

	f.Fuzz(func(t *testing.T, data []byte) {
		encoded := make([]byte, EncodedLen(len(data)))
		Encode(encoded, data)

		decoded := make([]byte, DecodedLen(len(encoded)))
		_, err := Decode(decoded, encoded)
		if err != nil {
			t.Fatalf("Decode error: %v", err)
		}

		if !bytes.Equal(data, decoded) {
			t.Fatalf("Round trip failed: original %q, decoded %q", data, decoded)
		}
	})
}
