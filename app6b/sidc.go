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
// rendered as '-' by Value and treated as "unspecified".
//
// SIDC deliberately does not implement fmt.Stringer. Use Value to render
// the wire-format string; it validates first and returns an error if any
// field byte is non-printable.
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

// ErrInvalidCharacter indicates a non-printable-ASCII byte was found in the input.
var ErrInvalidCharacter = errors.New("sidc/app6b: invalid character, expected printable ASCII")

// Value returns the canonical 15-character SIDC representation. It calls
// Validate first; if Validate returns an error, Value returns that error
// and the empty string. On success the result is guaranteed to be exactly
// 15 printable ASCII bytes.
func (s SIDC) Value() (string, error) {
	if err := s.Validate(); err != nil {
		return "", err
	}
	return s.render(), nil
}

// MarshalText implements encoding.TextMarshaler. Encoding to JSON or XML
// will fail if the SIDC is not valid.
func (s SIDC) MarshalText() ([]byte, error) {
	v, err := s.Value()
	if err != nil {
		return nil, err
	}
	return []byte(v), nil
}

// UnmarshalText implements encoding.TextUnmarshaler. The input must be a
// structurally valid SIDC (Parse). Callers wanting stricter semantic
// validation should call Validate on the result.
func (s *SIDC) UnmarshalText(text []byte) error {
	parsed, err := Parse(string(text))
	if err != nil {
		return err
	}
	*s = parsed
	return nil
}

// render produces the canonical encoding without validating. Unexported so
// the only public path to a SIDC string forces validation.
func (s SIDC) render() string {
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

// Parse parses a 15-character letter-based SIDC. Each byte must be printable
// ASCII (0x20-0x7E). Placeholder '-' characters are accepted at any position
// and stored as the byte '-'.
func Parse(s string) (sidc SIDC, err error) {
	if len(s) != SIDCLength {
		return SIDC{}, fmt.Errorf("%w: got %d", ErrInvalidLength, len(s))
	}
	for i := range len(s) {
		if s[i] < 0x20 || s[i] > 0x7E {
			return SIDC{}, fmt.Errorf("%w: position %d has byte 0x%02x", ErrInvalidCharacter, i, s[i])
		}
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
