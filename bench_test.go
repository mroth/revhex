package revhex

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"
)

func BenchmarkEncode(b *testing.B) {
	for _, size := range []int{256, 1024, 4096, 16384} {
		b.Run(fmt.Sprintf("%v", size), func(b *testing.B) {
			src := bytes.Repeat([]byte{2, 3, 5, 7, 9, 11, 13, 17}, size/8)
			sink := make([]byte, 2*size)
			b.SetBytes(int64(size))
			for b.Loop() {
				Encode(sink, src)
			}
		})
	}
}

func BenchmarkDecode(b *testing.B) {
	for _, size := range []int{256, 1024, 4096, 16384} {
		b.Run(fmt.Sprintf("%v", size), func(b *testing.B) {
			src := bytes.Repeat([]byte{'z', 'y', 'x', 'w', 'v', 'u', 't', 's'}, size/8)
			sink := make([]byte, size/2)
			b.SetBytes(int64(size))
			for b.Loop() {
				Decode(sink, src)
			}
		})
	}
}

func BenchmarkDecodeString(b *testing.B) {
	for _, size := range []int{256, 1024, 4096, 16384} {
		b.Run(fmt.Sprintf("%v", size), func(b *testing.B) {
			src := strings.Repeat("zyxwvuts", size/8)
			b.SetBytes(int64(size))
			for b.Loop() {
				DecodeString(src)
			}
		})
	}
}

func BenchmarkEncoder(b *testing.B) {
	for _, size := range []int{256, 1024, 4096, 16384} {
		encoder := NewEncoder(io.Discard)
		b.Run(fmt.Sprintf("%v", size), func(b *testing.B) {
			src := bytes.Repeat([]byte{2, 3, 5, 7, 9, 11, 13, 17}, size/8)
			b.SetBytes(int64(size))
			for b.Loop() {
				encoder.Write(src)
			}
		})
	}
}
