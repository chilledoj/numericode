package numericode

import (
	"encoding/binary"
	"math"
	"strings"
)

const stdCharSet = " ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_/"

// Encoding holds the charset for encoding
type Encoding struct {
	charset string
}

// NewEncoding returns a new Encoding charset
func NewEncoding(charset string) *Encoding {
	e := new(Encoding)
	if len(charset) < 10 || len(charset) > 64 {
		// silent fail
		return e
	}
	e.charset = charset
	return e
}

// StdEncoding returns a standard charset encoding
var StdEncoding = NewEncoding(stdCharSet)

/*
MaxChars returns the maximum number of characters that can be stored in the 32
bit integer using the stored charset. If the length (l) of the charset is not within
the bounds
	10 < l <= 64
then the charset is invalid and MaxChars returns 0.
*/
func (e *Encoding) MaxChars() int {
	if len(e.charset) < 10 || len(e.charset) > 64 {
		return 0
	}
	return len(powermap32[len(e.charset)]) - 1
}

/*
Encode will convert characters from src into LittleEndian uint32 bytes in dst.
If the src has a character not in the encoding character set then the result
will be 0.
*/
func (e *Encoding) Encode(dst, src []byte) {
	if len(src) > e.MaxChars() {
		return
	}
	var tot uint32
	for i := 0; i < len(src); i++ {
		idx := strings.Index(e.charset, string(src[i]))
		if idx == -1 {
			return
		}
		//temp = append(temp, uint8(idx))
		tot += uint32(uint8(idx)) * powermap32[len(e.charset)][i]
	}
	binary.LittleEndian.PutUint32(dst, tot)
}

// EncodeToUint32 will encode direct to an uint32
func (e *Encoding) EncodeToUint32(src []byte) uint32 {
	dst := make([]byte, 4)
	e.Encode(dst, src)
	return binary.LittleEndian.Uint32(dst)
}

/*
Decode will decode an unsigned interger into a destination byte slice.

N.B. the zerofill boolean will determine whether to pad the destination slice
with a null string byte (value 0) or the zero value of the charset.

*/
func (e *Encoding) Decode(dst []byte, i uint32, zerofill bool) {
	l := uint32(len(e.charset))
	num := i
	counter := 0
	for counter < e.MaxChars() {
		if num > 0 {
			x := math.Mod(float64(num), float64(l))
			num = num / l
			dst[counter] = e.charset[uint8(x)]
		} else if num <= 0 && zerofill {
			dst[counter] = e.charset[0]
		} else if num <= 0 && !zerofill {
			break // no need for pointless cycles
		}
		counter++
	}
}
