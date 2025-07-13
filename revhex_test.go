package revhex

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

type encDecTest struct {
	enc string
	dec []byte
}

// Round-trip tests for encoding and decoding.
var encDecTests = []encDecTest{
	// testcase byte values from encoding/hex tests
	{"", []byte{}},
	{"zzzyzxzwzvzuztzs", []byte{0, 1, 2, 3, 4, 5, 6, 7}},
	{"zrzqzpzoznzmzlzk", []byte{8, 9, 10, 11, 12, 13, 14, 15}},
	{"kzkykxkwkvkuktks", []byte{0xf0, 0xf1, 0xf2, 0xf3, 0xf4, 0xf5, 0xf6, 0xf7}},
	{"krkqkpkoknkmklkk", []byte{0xf8, 0xf9, 0xfa, 0xfb, 0xfc, 0xfd, 0xfe, 0xff}},
	{"ts", []byte{'g'}},
	{"lwpy", []byte{0xe3, 0xa1}},

	// test cases from https://github.com/jj-vcs/jj/blob/main/lib/src/hex_util.rs
	{"zyxwvutsrqponmlk", []byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef}},
}

// Decode only tests. Ensure Decode is case-insensitive.
var decOnlyTests = []encDecTest{
	{"ZZZYZXZWZVZUZTZS", []byte{0, 1, 2, 3, 4, 5, 6, 7}},
	{"ZZZYZXZWzvzuztzs", []byte{0, 1, 2, 3, 4, 5, 6, 7}},
	{"zzzyzxzWZVZUZTZs", []byte{0, 1, 2, 3, 4, 5, 6, 7}},
}

// Decoding tests (combined round-trip tests and decode-only tests).
var decTests = append(encDecTests, decOnlyTests...)

func TestEncode(t *testing.T) {
	for _, tt := range encDecTests {
		t.Run(tt.enc, func(t *testing.T) {
			dst := make([]byte, EncodedLen(len(tt.dec)))
			n := Encode(dst, tt.dec)
			if n != len(dst) {
				t.Errorf("Encode(%q) = %d bytes, want %d", tt.dec, n, len(dst))
			}
			if got := string(dst); got != tt.enc {
				t.Errorf("Encode(%q) = %q, want %q", tt.dec, got, tt.enc)
			}
		})
	}
}

func TestAppendEncode(t *testing.T) {
	for _, tt := range encDecTests {
		t.Run(tt.enc, func(t *testing.T) {
			dst := []byte("lead")
			dst = AppendEncode(dst, tt.dec)
			if got, want := string(dst), "lead"+tt.enc; got != want {
				t.Errorf("AppendEncode(lead, %q) = %q, want %q", tt.dec, got, want)
			}
		})
	}
}

func TestEncodeToString(t *testing.T) {
	for _, tt := range encDecTests {
		t.Run(tt.enc, func(t *testing.T) {
			if got := EncodeToString(tt.dec); got != tt.enc {
				t.Errorf("EncodeToString(%q) = %s, want %s", tt.dec, got, tt.enc)
			}
		})
	}
}

func TestDecode(t *testing.T) {
	for _, tt := range decTests {
		t.Run(tt.enc, func(t *testing.T) {
			dst := make([]byte, DecodedLen(len(tt.enc)))
			n, err := Decode(dst, []byte(tt.enc))
			if err != nil {
				t.Errorf("Decode(%q) unexpected error: %v", tt.enc, err)
			}
			if n != len(dst) {
				t.Errorf("Decode(%q) = %d bytes, want %d", tt.enc, n, len(dst))
			}
			if !bytes.Equal(dst, tt.dec) {
				t.Errorf("Decode(%q) = %v, want %v", tt.enc, dst, tt.dec)
			}
		})
	}
}

func TestAppendDecode(t *testing.T) {
	for _, tt := range decTests {
		t.Run(tt.enc, func(t *testing.T) {
			dst := []byte("lead")
			// var err error
			dst, err := AppendDecode(dst, []byte(tt.enc))
			if err != nil {
				t.Errorf("AppendDecode(lead, %q) unexpected error: %v", tt.enc, err)
			}
			if got, want := string(dst), "lead"+string(tt.dec); got != want {
				t.Errorf("AppendDecode(lead, %q) = %q, want %q", tt.enc, got, want)
			}
		})
	}
}

func TestDecodeString(t *testing.T) {
	for _, tt := range decTests {
		t.Run(tt.enc, func(t *testing.T) {
			got, err := DecodeString(tt.enc)
			if err != nil {
				t.Errorf("DecodeString(%q) unexpected error: %v", tt.enc, err)
			}
			if !bytes.Equal(got, tt.dec) {
				t.Errorf("DecodeString(%q) = %v, want %v", tt.enc, got, tt.dec)
			}
		})
	}
}

func TestEncoderDecoder(t *testing.T) {
	for _, multiplier := range []int{1, 128, 192} {
		for _, test := range encDecTests {
			input := bytes.Repeat(test.dec, multiplier)
			output := strings.Repeat(test.enc, multiplier)

			var buf bytes.Buffer
			enc := NewEncoder(&buf)
			r := struct{ io.Reader }{bytes.NewReader(input)} // io.Reader only; not io.WriterTo
			if n, err := io.CopyBuffer(enc, r, make([]byte, 7)); n != int64(len(input)) || err != nil {
				t.Errorf("encoder.Write(%q*%d) = (%d, %v), want (%d, nil)", test.dec, multiplier, n, err, len(input))
				continue
			}

			if encDst := buf.String(); encDst != output {
				t.Errorf("buf(%q*%d) = %v, want %v", test.dec, multiplier, encDst, output)
				continue
			}

			dec := NewDecoder(&buf)
			var decBuf bytes.Buffer
			w := struct{ io.Writer }{&decBuf} // io.Writer only; not io.ReaderFrom
			if _, err := io.CopyBuffer(w, dec, make([]byte, 7)); err != nil || decBuf.Len() != len(input) {
				t.Errorf("decoder.Read(%q*%d) = (%d, %v), want (%d, nil)", test.enc, multiplier, decBuf.Len(), err, len(input))
			}

			if !bytes.Equal(decBuf.Bytes(), input) {
				t.Errorf("decBuf(%q*%d) = %v, want %v", test.dec, multiplier, decBuf.Bytes(), input)
				continue
			}
		}
	}
}

var decErrTests = []struct {
	in  string
	out string
	err error
}{
	{"", "", nil},
	{"z", "", ErrLength},
	{"9mvpp", "", InvalidByteError('9')},
	{"mvpp9", "\xd4\xaa", InvalidByteError('9')},
	{"zzzzz", "\x00\x00", ErrLength},
	{"zj", "", InvalidByteError('j')},
	{"zzjj", "\x00", InvalidByteError('j')},
	{"z\x01", "", InvalidByteError('\x01')},
	{"kkllm", "\xff\xee", ErrLength},
}

func TestDecodeErr(t *testing.T) {
	for _, tt := range decErrTests {
		t.Run(tt.in, func(t *testing.T) {
			out := make([]byte, DecodedLen(len(tt.in)))
			n, err := Decode(out, []byte(tt.in))
			if string(out[:n]) != tt.out {
				t.Errorf("Decode(%q) = %q, want %q", tt.in, string(out[:n]), tt.out)
			}
			if err != tt.err {
				t.Errorf("Decode(%q) error = %v, want %v", tt.in, err, tt.err)
			}
		})
	}
}

func TestDecodeStringErr(t *testing.T) {
	for _, tt := range decErrTests {
		t.Run(tt.in, func(t *testing.T) {
			out, err := DecodeString(tt.in)
			if string(out) != tt.out {
				t.Errorf("DecodeString(%q) = %q, want %q", tt.in, string(out), tt.out)
			}
			if err != tt.err {
				t.Errorf("DecodeString(%q) error = %v, want %v", tt.in, err, tt.err)
			}
		})
	}
}

func TestDecoderErr(t *testing.T) {
	for _, tt := range decErrTests {
		dec := NewDecoder(strings.NewReader(tt.in))
		out, err := io.ReadAll(dec)
		wantErr := tt.err
		// Decoder is reading from stream, so it reports io.ErrUnexpectedEOF instead of ErrLength.
		if wantErr == ErrLength {
			wantErr = io.ErrUnexpectedEOF
		}
		if string(out) != tt.out || err != wantErr {
			t.Errorf("NewDecoder(%q) = %q, %v, want %q, %v", tt.in, out, err, tt.out, wantErr)
		}
	}
}
