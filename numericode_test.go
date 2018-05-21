package numericode_test

import (
	"fmt"
	"testing"

	"github.com/chilledoj/numericode"
)

func TestNumericode(t *testing.T) {

	tsts := []struct {
		cde string
		exp uint32
	}{
		{"A", 1},
		{"CODE", 327003},       //  3 + 15 * 40^1 +  4 * 40^2 +  5 * 40^3
		{"//////", 4095999999}, // 39 +  39 * 40 +  39 * 1600 +  39 * 64000 +  39 * 2560000 +  39 * 102400000
	}
	for _, tst := range tsts {
		act, err := numericode.FromString(tst.cde)
		if err != nil {
			t.Fatalf("Returned error from .FromString: %+v", err)
		}
		i, err := act.ToUint32()
		if err != nil {
			t.Fatalf("Returned error from .ToUint32: %+v", err)
		}
		if i != tst.exp {
			t.Errorf("Actual (%v) != (%v) Expected ", i, tst.exp)
		}
		if v := act.Value(); v != tst.exp {
			t.Errorf("Value (%v) != (%v) Expected ", v, tst.exp)
		}
	}
}

func TestNumericode_ToUint32(t *testing.T) {
	tsts := []struct {
		ip     []byte
		expErr bool
		exp    uint32
	}{
		{[]byte{30, 40, 50, 60, 70, 80, 90, 100, 110}, true, 0},
		{[]byte{1, 1}, false, 41},
		{[]byte{39, 39, 39, 39, 39, 39}, false, 4095999999},
	}

	for _, tst := range tsts {
		n := numericode.Numericode(tst.ip)
		i, err := n.ToUint32()
		if err != nil && !tst.expErr {
			t.Fatalf("Return error from .ToUint32: %+v", err)
		}
		if tst.expErr && err == nil {
			t.Errorf("Expected error but returned nil: %+v", tst.ip)
		}
		if !tst.expErr && i != tst.exp {
			t.Errorf("Actual (%d) != (%d) Expected", i, tst.exp)
		}
	}
}
func TestNumericode_MaxChars(t *testing.T) {
	tsts := []struct {
		charset   string
		expMaxLen int
	}{
		{"1234567890123456789012345678901234567890", 6},
		{"12345678901234567890123456789012345678901", 5},
	}

	for _, tst := range tsts {
		numericode.OverideCharSet(tst.charset)
		act := numericode.MaxChars()
		if act != tst.expMaxLen {
			t.Errorf("Actual (%d) != Expected (%d) max chars. length(%d)", act, tst.expMaxLen, len(tst.charset))
		}
	}
	numericode.OverideCharSet(numericode.DefaultCharSet)
}

func TestNumericode_FromString(t *testing.T) {
	tsts := []struct {
		ip     string
		experr bool
	}{
		{"A", false},
		{"CODE", false},
		{"LARGECODE", true},
		{"Invalid code", true},
		{"inval", true},
		{"      ", false},
		{"       ", true},
	}
	for _, tst := range tsts {
		n, err := numericode.FromString(tst.ip)
		if err != nil && !tst.experr {
			t.Fatalf("Return error from .FromString: %+v", err)
		}
		if tst.experr && err == nil {
			t.Errorf("Expected error for code %s, but none returned: %+v", tst.ip, n)
		}
		if !tst.experr && n.String() != tst.ip {
			t.Errorf("Actual (%v) != (%v) Expected", n, tst.ip)
		}
	}
}

func TestNumericode_FromUint32(t *testing.T) {
	tsts := []struct {
		ip  uint32
		exp string
	}{
		{255, "FF"},
	}
	numericode.OverideCharSet("0123456789ABCDEF")
	for _, tst := range tsts {
		n, err := numericode.FromUint32(tst.ip)
		if err != nil {
			t.Fatalf("Return error from .FromUint32: %+v", err)
		}
		fmt.Println(n.RawString())
		if n.String() != tst.exp {
			t.Errorf("String method err:  Act(%s) != Input(%s)", n.String(), tst.exp)
		}
	}
	numericode.OverideCharSet(numericode.DefaultCharSet)
}

func TestNumericode_String(t *testing.T) {
	tsts := []string{"A", "CODE", "LOWER", "UPPER", "//////"}
	for _, tst := range tsts {
		n, err := numericode.FromString(tst)
		if err != nil {
			t.Fatalf("Return error from .FromString: %+v", err)
		}
		if n.String() != tst {
			t.Errorf("String method err:  Act(%s) != Input(%s)", n.String(), tst)
		}
	}
}

func TestNumericode_OverideCharSet(t *testing.T) {
	tsts := []struct {
		charset string
		ip      string
		exp     uint32
	}{
		{"0123456789ABCDEF", "FF", 255},
		{"0123456789ABCDEF", "FFFFFF", 16777215},
		{"01", "0", 27}, // Silent fallback to DefaultCharSet
	}

	for _, tst := range tsts {
		numericode.OverideCharSet(tst.charset)
		n, err := numericode.FromString(tst.ip)
		if err != nil {
			t.Errorf("Return error from .FromString: %+v", err)
		}
		if v := n.Value(); v != tst.exp {
			t.Errorf("Actual (%d) != (%d) Expected", v, tst.exp)
		}
		// RESET
		numericode.OverideCharSet(numericode.DefaultCharSet)
	}
}
