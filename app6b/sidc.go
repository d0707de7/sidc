package app6b

import (
	"errors"
	"fmt"
)

// SIDCLength is the canonical length of an APP-6 B/C letter-based SIDC.
const SIDCLength = 15

// Placeholder is the character used in the canonical string when a single
// position is unspecified.
const Placeholder = '-'

// SIDC is a parsed 15-character letter-based SIDC for APP-6 B and C and
// MIL-STD-2525 B/C. Each field stores its raw byte value; the zero byte is
// rendered as '-' in String and treated as "unspecified".
type SIDC struct {
	CodingScheme    CodingScheme
	Affiliation     Affiliation
	BattleDimension BattleDimension
	Status          Status
	FunctionID      FunctionID
	SymbolModifier1 byte
	SymbolModifier2 byte
	CountryCode     [2]byte
	OrderOfBattle   byte
}

// FunctionID is the 6-character function ID at positions 4-9.
type FunctionID [6]byte

// ErrInvalidLength indicates the input was not 15 characters.
var ErrInvalidLength = errors.New("sidc/app6b: invalid length, expected 15 characters")

// String returns the canonical 15-character representation. Unset (zero)
// bytes render as '-'.
func (s SIDC) String() string {
	var buf [SIDCLength]byte
	buf[0] = orDash(byte(s.CodingScheme))
	buf[1] = orDash(byte(s.Affiliation))
	buf[2] = orDash(byte(s.BattleDimension))
	buf[3] = orDash(byte(s.Status))
	for i := range 6 {
		buf[4+i] = orDash(s.FunctionID[i])
	}
	buf[10] = orDash(s.SymbolModifier1)
	buf[11] = orDash(s.SymbolModifier2)
	buf[12] = orDash(s.CountryCode[0])
	buf[13] = orDash(s.CountryCode[1])
	buf[14] = orDash(s.OrderOfBattle)
	return string(buf[:])
}

func orDash(b byte) byte {
	if b == 0 {
		return Placeholder
	}
	return b
}

// Parse parses a 15-character letter-based SIDC. Placeholder '-' characters
// are accepted at any position and stored as the byte '-'.
func Parse(s string) (sidc SIDC, err error) {
	if len(s) != SIDCLength {
		return SIDC{}, fmt.Errorf("%w: got %d", ErrInvalidLength, len(s))
	}
	sidc = SIDC{
		CodingScheme:    CodingScheme(s[0]),
		Affiliation:     Affiliation(s[1]),
		BattleDimension: BattleDimension(s[2]),
		Status:          Status(s[3]),
		SymbolModifier1: s[10],
		SymbolModifier2: s[11],
		CountryCode:     [2]byte{s[12], s[13]},
		OrderOfBattle:   s[14],
	}
	for i := range 6 {
		sidc.FunctionID[i] = s[4+i]
	}
	return sidc, nil
}
