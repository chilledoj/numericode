package numericode_test

import (
	"encoding/binary"
	"fmt"
	"strings"
	"testing"

	nc "github.com/chilledoj/numericode"
)

func TestEncoding_MaxChars(t *testing.T) {
	if n := nc.StdEncoding.MaxChars(); n != 6 {
		t.Errorf("Max Chars for stdCharSet returns %d not 6", n)
	}

	e := nc.NewEncoding("ABCDEFG")
	if n := e.MaxChars(); n != 0 {
		t.Errorf("Max Chars for invalid charset returns %d not 0", n)
	}
}

func TestEncoding_Encode(t *testing.T) {
	tsts := []struct {
		name     string
		encoding *nc.Encoding
		ip       string
		exp      uint32
	}{
		{"Std - 'A'", nc.StdEncoding, "A", 1},
		{"Std - 'Z'", nc.StdEncoding, "Z", 26},
		{"Std - '/'", nc.StdEncoding, "/", 39},
		{"Std - ' A'", nc.StdEncoding, " A", 40},
		{"Std - '  A'", nc.StdEncoding, "  A", 1600},
		// This should cause a silent fail - hence 0
		{"Std - '1234567'", nc.StdEncoding, "1234567", 0},
		// Note that the code is actually backwards to the actual hex represnetation because the code is effectively little Endian
		{"Hex - 'EF'", nc.NewEncoding("0123456789ABCDEF"), "EF", 254},
		// This should cause a silent fail - hence 0
		{"Hex - 'FG'", nc.NewEncoding("0123456789ABCDEF"), "FG", 0},
	}

	for _, tst := range tsts {
		dst := make([]byte, 4)
		tst.encoding.Encode(dst, []byte(tst.ip))
		u := binary.LittleEndian.Uint32(dst)
		if u != tst.exp {
			t.Errorf("%s: Actual(%d) != Expected(%d)", tst.name, u, tst.exp)
		}
	}
}

func TestEncoding_EncodeToUint32(t *testing.T) {
	tsts := []struct {
		name     string
		encoding *nc.Encoding
		ip       string
		exp      uint32
	}{
		{"Std - 'A'", nc.StdEncoding, "A", 1},
		{"Std - 'Z'", nc.StdEncoding, "Z", 26},
		{"Std - '/'", nc.StdEncoding, "/", 39},
		{"Std - ' A'", nc.StdEncoding, " A", 40},
		{"Std - '  A'", nc.StdEncoding, "  A", 1600},
		{"Std - '1234567'", nc.StdEncoding, "1234567", 0},
		{"Hex - 'EF'", nc.NewEncoding("0123456789ABCDEF"), "EF", 254},
		{"Hex - 'FG'", nc.NewEncoding("0123456789ABCDEF"), "FG", 0},
	}

	for _, tst := range tsts {
		u := tst.encoding.EncodeToUint32([]byte(tst.ip))
		if u != tst.exp {
			t.Errorf("%s: Actual(%d) != Expected(%d)", tst.name, u, tst.exp)
		}
	}
}

func TestEncoding_Decode(t *testing.T) {
	tsts := []struct {
		name     string
		encoding *nc.Encoding
		exp      string
		ip       uint32
		zerofill bool
	}{
		// Note tha the following tests include padding with the zero character - which is space
		{"Std - 'A'", nc.StdEncoding, "A     ", 1, true},
		{"Std - 'Z'", nc.StdEncoding, "Z     ", 26, true},
		{"Std - '/'", nc.StdEncoding, "/     ", 39, true},
		{"Std - ' A'", nc.StdEncoding, " A    ", 40, true},
		{"Std - '  A'", nc.StdEncoding, "  A   ", 1600, true},
		{"Std - '  A'", nc.StdEncoding, "  A", 1600, false},
		// Padding is with zero character which is 0 in hex
		{"Hex - 'EF'", nc.NewEncoding("0123456789ABCDEF"), "EF00000", 254, true},
		{"Hex - 'FF'", nc.NewEncoding("0123456789ABCDEF"), "FF", 255, false},
	}

	for _, tst := range tsts {
		dst := make([]byte, len(tst.exp)) // IS THIS A CHEAT? Should it reallocate the dst for each new byte?
		tst.encoding.Decode(dst, tst.ip, tst.zerofill)
		if n := strings.Compare(string(dst), tst.exp); n != 0 {
			t.Errorf("%s: Actual(%v) != Expected(%v); %d diff", tst.name, dst, []byte(tst.exp), n)
		}
	}
}

func ExampleEncoding_Encode() {
	e := nc.StdEncoding
	dst := make([]byte, 4)
	e.Encode(dst, []byte("A"))
	fmt.Println(binary.LittleEndian.Uint32(dst)) // Outputs: 1
	e.Decode(dst, 1, false)
	fmt.Println(dst) // Outputs: "A"
}

func ExampleEncoding_EncodeToUint32() {
	e := nc.StdEncoding

	i := e.EncodeToUint32([]byte("NCODE"))
	fmt.Println(i) // Outputs: 3
}

func ExampleEncoding_Decode() {
	e := nc.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	dst := make([]byte, 6)
	e.Decode(dst, 1, true)
	fmt.Println(dst) // Outputs: "BAAAAA"
	e.Decode(dst, 1, false)
	fmt.Println(dst) // Outputs: "B"
}
