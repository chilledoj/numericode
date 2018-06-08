# Numericode
## Type Version
**UPDATE**

This package provides a simple conversion from a string, utilising a
fixed string character set, to an unsigned integer. The string provided is
used in a Little Endian fashion so that adding more characters will increase
the integer value.

The general purpose of this is to allow users to define config keys which
will only ever use a fixed, reduced character set (subset of ASCII) and to store
it in an efficient manner within the DB. There will be scenarios where the
storage of an unsigned 32 bit integer is less efficient than the storage of a
certain amount of characters e.g. anything less than 4 bytes like 'CCY'.

*N.B. The current implementation only uses uint32 as the integer value.*

> e.g. We might want to store the code HELLO as a key in the DB. This would utilise 5 bytes (assuming UTF-8 or other single byte encoding). However using this conversion will store this as a 4 byte unsigned integer.

This is all based upon the idea of base conversion:
If we assume that we will only use the uppercase ASCII characters for the
definition of a code, then we can say that that character set represents an
integer value on a scale. This is analogous to hexadecimal where the character
set is 0-9 and A-F, representing the integers 0-15.

The math involved here is to check the maximum integer that can be stored with
a character set (of length l) for a given number of characters in the code (n).
The maximum number that can be stored is l<sup>n</sup>. This will be stored in a 32
bit (unsigned) integer which has a maximum value of 2<sup>32</sup>-1. Therefore the
constrain is simply l<sup>n</sup> < 2<sup>32</sup>-1.
