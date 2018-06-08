/*Package nctype provides a simple conversion from a string, utilising a
fixed string character set, to an unsigned integer. The string provided is
used in a Little Endian fashion so that adding more characters will increase
the integer value.

The general purpose of this is to allow users to define config keys which
will only ever use a fixed, reduced character set (subset of ASCII) and to store
it in an efficient manner within the DB. There will be scenarios where the
storage of an unsigned 32 bit integer is less efficient than the storage of a
certain amount of characters e.g. anything less than 4 bytes like 'CCY'.

*N.B. The current implementation only uses uint32 as the integer value.*

> e.g. We might want to store the code HELLO as a key in the DB. This would
utilise 5 bytes (assuming UTF-8 or other single byte encoding). However using
this conversion will store this as a 4 byte unsigned integer.

This is all based upon the idea of base conversion:
If we assume that we will only use the uppercase ASCII characters for the
definition of a code, then we can say that that character set represents an
integer value on a scale. This is analogous to hexadecimal where the character
set is 0-9 and A-F, representing the integers 0-15.

The math involved here is to check the maximum integer that can be stored with
a character set (of length l) for a given number of characters in the code (n).
The maximum number that can be stored is l<sup>n</sup>. This will be stored in a 32
bit (unsigned) integer which has a maximum value of 2<sup>32-1</sup>. Therefore the
constrain is simply l<sup>n</sup> < 2<sup>32-1</sup>.
*/
package nctype

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
)

// DefaultCharSet is the charset used. It allows codes up to 6 chars long
//  ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_/
const DefaultCharSet string = " ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_/"

var useCharSet = DefaultCharSet

// OverideCharSet will override the character set for encoding.
//The provided string set must be between 10 and 64 characters in length. The
//function will silently fail and revert back to the DefaultCharSet if this is
//not met. There is no real reason around these limits other than the powermap
//has been generated between these numbers.
func OverideCharSet(s string) {
	if len(s) < 10 || len(s) > 64 {
		// silent fail
		return
	}
	useCharSet = s
}

/*
MaxChars is now public. The original implementation of maxChars used the correct mathematical formula
to calculate the maximum number of characters allows in the code.
	func maxChars() int {
		return int(math.Floor(32.0 / math.Log2(float64(len(useCharSet)))))
	}
However, seeing as we need to store a powermap for the different lengths of
character sets, then we can simply look up the length of the slice of powers.
*/
func MaxChars() int {
	return len(powermap32[len(useCharSet)]) - 1
}

// Numericode holds the slice of uint8 (we've chosen bytes here) which represent
//the location within the character map.
type Numericode []byte

// FromString creates the Numericode from a string
func FromString(s string) (Numericode, error) {
	if len(s) > MaxChars() {
		return nil, fmt.Errorf("invalid string length (%d) %s", len(s), s)
	}
	var n Numericode
	for i := 0; i < len(s); i++ {
		idx := strings.Index(useCharSet, string(s[i]))
		if idx == -1 {
			return nil, fmt.Errorf("Invalid character at position %d", i)
		}
		n = append(n, uint8(idx))
	}
	return n, nil
}

// FromUint32 creates the numericode from a Uint32
func FromUint32(i uint32) (Numericode, error) {
	l := len(useCharSet)
	str := []byte{}

	//fmt.Printf("FromUint32: %d\n", i)
	num := i

	for num > 0 {
		x := math.Mod(float64(num), float64(l))
		num = num / uint32(l)

		//fmt.Printf("NUM: %d, REM: %v, %s\n", num, int(x), string(useCharSet[int(x)]))
		str = append(str, uint8(x))
	}
	return Numericode(str), nil
}

// ToUint32 creates a uint32 from the byte slice
func (n Numericode) ToUint32() (uint32, error) {
	if len(n) > MaxChars() {
		return 0, fmt.Errorf("uint32 invalid - truncation will occur")
	}

	tot := uint32(n[0])
	//fmt.Printf(" %d", tot)
	for i := 1; i < len(n); i++ {
		//fmt.Printf(" +  %d * %d", uint32(n[i]), powermap32[len(useCharSet)][i-1])
		tot += uint32(n[i]) * powermap32[len(useCharSet)][i]
	}
	//fmt.Printf(" = %d\n", tot)
	return tot, nil
}

// String implements the stringe interface
func (n Numericode) String() string {
	var str string
	for i := 0; i < len(n); i++ {
		str = str + string(useCharSet[int(n[i])])
	}
	return str
}

// Value implements the value interface
func (n Numericode) Value() uint32 {
	val, _ := n.ToUint32()
	return val
}

// RawString returns the byte slice string
func (n Numericode) RawString() string {
	return fmt.Sprintf("%v", []byte(n))
}

// MarshalJSON for sending data via JSON (uses string method)
func (n Numericode) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.String())
}
