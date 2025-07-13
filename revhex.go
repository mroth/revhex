// Portions of this file are derived from Go standard library encoding/hex/hex.go
// (https://github.com/golang/go/blob/master/src/encoding/hex/hex.go).
//
// Copyright 2009 The Go Authors. All rights reserved.
// Use of that source code is governed by a BSD-style license that can be found
// in the LICENSE file.

// Package revhex implements reverse hexadecimal encoding and decoding.
//
// The reverse hex format is just normal 16-base hexadecimal encoding but using
// `z-k` as the alphabet instead of `0-9a-f`.
//
// The implementations in this package are closely modeled after the standard library
// [encoding/hex] package, and have identical APIs and performance characteristics.
package revhex

import (
	"errors"
	"fmt"
	"io"
	"slices"
)

const (
	revhexTable        = "zyxwvutsrqponmlk"
	reverseRevhexTable = "" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" + // 0x00-0x0f
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" + // 0x10-0x1f
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" + // 0x20-0x2f
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" + // 0x30-0x3f
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\x0f\x0e\x0d\x0c\x0b" + // 0x40-0x4f, 0x4B-0x4F = K-O
		"\x0a\x09\x08\x07\x06\x05\x04\x03\x02\x01\x00\xff\xff\xff\xff\xff" + // 0x50-0x5f, 0x50-0x5A = P-Z
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\x0f\x0e\x0d\x0c\x0b" + // 0x60-0x6f, 0x6B-0x6F = k-o
		"\x0a\x09\x08\x07\x06\x05\x04\x03\x02\x01\x00\xff\xff\xff\xff\xff" + // 0x70-0x7f, 0x70-0x7A = p-z
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" + // 0x80-0x8f
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" + // 0x90-0x9f
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" + // 0xa0-0xaf
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" + // 0xb0-0xbf
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" + // 0xc0-0xcf
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" + // 0xd0-0xdf
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" + // 0xe0-0xef
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" // 0xf0-0xff
)

/* ------------------------------------ ENCODING ------------------------------------ */

// EncodedLen returns the length of an encoding of n source bytes.
// Specifically, it returns n * 2.
func EncodedLen(n int) int { return n * 2 }

// Encode encodes src into [EncodedLen](len(src)) bytes of dst.
// As a convenience, it returns the number of bytes written to dst, but this value is always [EncodedLen](len(src)).
// Encode implements reverse hexadecimal encoding.
func Encode(dst, src []byte) int {
	j := 0
	for _, v := range src {
		dst[j] = revhexTable[v>>4]
		dst[j+1] = revhexTable[v&0x0f]
		j += 2
	}
	return len(src) * 2
}

// AppendEncode appends the reverse hexadecimally encoded src to dst and returns the extended buffer.
func AppendEncode(dst, src []byte) []byte {
	n := EncodedLen(len(src))
	dst = slices.Grow(dst, n)
	Encode(dst[len(dst):][:n], src)
	return dst[:len(dst)+n]
}

// EncodeToString returns the reverse hexadecimal encoding of src.
func EncodeToString(src []byte) string {
	dst := make([]byte, EncodedLen(len(src)))
	Encode(dst, src)
	return string(dst)
}

/* ------------------------------------ DECODING ------------------------------------ */

// ErrLength reports an attempt to decode an odd-length input using [Decode] or [DecodeString].
// The stream-based Decoder returns [io.ErrUnexpectedEOF] instead of ErrLength.
var ErrLength = errors.New("encoding/revhex: odd length reverse hex string")

// InvalidByteError values describe errors resulting from an invalid byte in a reverse hex string.
type InvalidByteError byte

func (e InvalidByteError) Error() string {
	return fmt.Sprintf("encoding/revhex: invalid byte: %#U", rune(e))
}

// DecodedLen returns the length of a decoding of x source bytes.
// Specifically, it returns x / 2.
func DecodedLen(x int) int { return x / 2 }

// Decode decodes src into [DecodedLen](len(src)) bytes, returning the actual number of bytes written to dst.
//
// Decode expects that src contains only reverse hexadecimal characters and that src has even length.
// If the input is malformed, Decode returns the number of bytes decoded before the error.
func Decode(dst, src []byte) (int, error) {
	i, j := 0, 1
	for ; j < len(src); j += 2 {
		p := src[j-1]
		q := src[j]

		a := reverseRevhexTable[p]
		b := reverseRevhexTable[q]
		if a > 0x0f {
			return i, InvalidByteError(p)
		}
		if b > 0x0f {
			return i, InvalidByteError(q)
		}
		dst[i] = (a << 4) | b
		i++
	}
	if len(src)%2 == 1 {
		// Check for invalid char before reporting bad length,
		// since the invalid char (if present) is an earlier problem.
		if reverseRevhexTable[src[j-1]] > 0x0f {
			return i, InvalidByteError(src[j-1])
		}
		return i, ErrLength
	}
	return i, nil
}

// AppendDecode appends the reverse hexadecimally decoded src to dst and returns the extended buffer.
// If the input is malformed, it returns the partially decoded src and an error.
func AppendDecode(dst, src []byte) ([]byte, error) {
	n := DecodedLen(len(src))
	dst = slices.Grow(dst, n)
	n, err := Decode(dst[len(dst):][:n], src)
	return dst[:len(dst)+n], err
}

// DecodeString returns the bytes represented by the reverse hexadecimal string s.
//
// DecodeString expects that src contains only reverse hexadecimal characters and that src has even length.
// If the input is malformed, DecodeString returns the bytes decoded before the error.
func DecodeString(s string) ([]byte, error) {
	dst := make([]byte, DecodedLen(len(s)))
	n, err := Decode(dst, []byte(s))
	return dst[:n], err
}

/* -------------------------------- ENCODER/DECODER -------------------------------- */

// bufferSize is the number of reverse hexadecimal characters to buffer in encoder and decoder.
const bufferSize = 1024

type encoder struct {
	w   io.Writer
	err error
	out [bufferSize]byte // output buffer
}

// NewEncoder returns an [io.Writer] that writes lowercase reverse hexadecimal characters to w.
func NewEncoder(w io.Writer) io.Writer {
	return &encoder{w: w}
}

func (e *encoder) Write(p []byte) (n int, err error) {
	for len(p) > 0 && e.err == nil {
		chunkSize := bufferSize / 2
		if len(p) < chunkSize {
			chunkSize = len(p)
		}

		var written int
		encoded := Encode(e.out[:], p[:chunkSize])
		written, e.err = e.w.Write(e.out[:encoded])
		n += written / 2
		p = p[chunkSize:]
	}
	return n, e.err
}

type decoder struct {
	r   io.Reader
	err error
	in  []byte           // input buffer (encoded form)
	arr [bufferSize]byte // backing array for in
}

// NewDecoder returns an [io.Reader] that decodes reverse hexadecimal characters from r.
// NewDecoder expects that r contain only an even number of reverse hexadecimal characters.
func NewDecoder(r io.Reader) io.Reader {
	return &decoder{r: r}
}

func (d *decoder) Read(p []byte) (n int, err error) {
	// Fill internal buffer with sufficient bytes to decode
	if len(d.in) < 2 && d.err == nil {
		var numCopy, numRead int
		numCopy = copy(d.arr[:], d.in) // Copies either 0 or 1 bytes
		numRead, d.err = d.r.Read(d.arr[numCopy:])
		d.in = d.arr[:numCopy+numRead]
		if d.err == io.EOF && len(d.in)%2 != 0 {

			if a := reverseRevhexTable[d.in[len(d.in)-1]]; a > 0x0f {
				d.err = InvalidByteError(d.in[len(d.in)-1])
			} else {
				d.err = io.ErrUnexpectedEOF
			}
		}
	}

	// Decode internal buffer into output buffer
	if numAvail := len(d.in) / 2; len(p) > numAvail {
		p = p[:numAvail]
	}
	numDec, err := Decode(p, d.in[:len(p)*2])
	d.in = d.in[2*numDec:]
	if err != nil {
		d.in, d.err = nil, err // Decode error; discard input remainder
	}

	if len(d.in) < 2 {
		return numDec, d.err // Only expose errors when buffer fully consumed
	}
	return numDec, nil
}
