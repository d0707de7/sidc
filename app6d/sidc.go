package app6d

import (
	"errors"
	"fmt"
	"strconv"
)

// SIDC is a parsed 20-digit Symbol Identification Code for APP-6 D or E.
//
// The zero value is the string "10000000000000000000": APP-6 D version 10,
// reality, pending affiliation, unknown symbol set, present, no amplifier,
// no entity, no modifiers. Use a composite literal for non-default values.
type SIDC struct {
	Version     Version
	Context     Context
	Affiliation Affiliation
	SymbolSet   SymbolSet
	Status      Status
	HQTFD       HQTFD
	Amplifier   Amplifier
	Entity      Entity
	Modifier1   Modifier1
	Modifier2   Modifier2
}

// SIDCLength is the canonical length of an APP-6 D/E SIDC string.
const SIDCLength = 20

// ErrInvalidLength indicates the input was not 20 characters.
var ErrInvalidLength = errors.New("sidc: invalid length, expected 20 characters")

// ErrInvalidCharacter indicates a non-digit was found in the input.
var ErrInvalidCharacter = errors.New("sidc: invalid character, expected digits only")

// String returns the canonical 20-digit SIDC representation. It never fails;
// values out of range for a field's width are masked to fit.
func (s SIDC) String() string {
	var buf [SIDCLength]byte
	writeDigits(buf[0:2], uint64(s.Version)%100)
	writeDigits(buf[2:3], uint64(s.Context)%10)
	writeDigits(buf[3:4], uint64(s.Affiliation)%10)
	writeDigits(buf[4:6], uint64(s.SymbolSet)%100)
	writeDigits(buf[6:7], uint64(s.Status)%10)
	writeDigits(buf[7:8], uint64(s.HQTFD)%10)
	writeDigits(buf[8:10], uint64(s.Amplifier)%100)
	writeDigits(buf[10:16], uint64(s.Entity)%1000000)
	writeDigits(buf[16:18], uint64(s.Modifier1)%100)
	writeDigits(buf[18:20], uint64(s.Modifier2)%100)
	return string(buf[:])
}

// writeDigits writes v zero-padded into b, least-significant digit last.
func writeDigits(b []byte, v uint64) {
	for i := len(b) - 1; i >= 0; i-- {
		b[i] = '0' + byte(v%10)
		v /= 10
	}
}

// Parse parses a 20-digit SIDC string into its component fields. The input
// must be exactly 20 ASCII digits. Parse performs structural validation only;
// it does not check that the entity exists in the symbol set, nor that the
// symbol set is valid for the version. Use Validate for those checks.
func Parse(s string) (sidc SIDC, err error) {
	if len(s) != SIDCLength {
		return SIDC{}, fmt.Errorf("%w: got %d", ErrInvalidLength, len(s))
	}
	for i := range len(s) {
		if s[i] < '0' || s[i] > '9' {
			return SIDC{}, fmt.Errorf("%w: position %d has %q", ErrInvalidCharacter, i, s[i])
		}
	}

	version, err := strconv.ParseUint(s[0:2], 10, 8)
	if err != nil {
		return SIDC{}, fmt.Errorf("sidc: parsing version: %w", err)
	}
	context, err := strconv.ParseUint(s[2:3], 10, 8)
	if err != nil {
		return SIDC{}, fmt.Errorf("sidc: parsing context: %w", err)
	}
	affiliation, err := strconv.ParseUint(s[3:4], 10, 8)
	if err != nil {
		return SIDC{}, fmt.Errorf("sidc: parsing affiliation: %w", err)
	}
	symbolSet, err := strconv.ParseUint(s[4:6], 10, 8)
	if err != nil {
		return SIDC{}, fmt.Errorf("sidc: parsing symbol set: %w", err)
	}
	status, err := strconv.ParseUint(s[6:7], 10, 8)
	if err != nil {
		return SIDC{}, fmt.Errorf("sidc: parsing status: %w", err)
	}
	hqtfd, err := strconv.ParseUint(s[7:8], 10, 8)
	if err != nil {
		return SIDC{}, fmt.Errorf("sidc: parsing HQTFD: %w", err)
	}
	amplifier, err := strconv.ParseUint(s[8:10], 10, 8)
	if err != nil {
		return SIDC{}, fmt.Errorf("sidc: parsing amplifier: %w", err)
	}
	entity, err := strconv.ParseUint(s[10:16], 10, 32)
	if err != nil {
		return SIDC{}, fmt.Errorf("sidc: parsing entity: %w", err)
	}
	modifier1, err := strconv.ParseUint(s[16:18], 10, 8)
	if err != nil {
		return SIDC{}, fmt.Errorf("sidc: parsing modifier1: %w", err)
	}
	modifier2, err := strconv.ParseUint(s[18:20], 10, 8)
	if err != nil {
		return SIDC{}, fmt.Errorf("sidc: parsing modifier2: %w", err)
	}

	return SIDC{
		Version:     Version(version),
		Context:     Context(context),
		Affiliation: Affiliation(affiliation),
		SymbolSet:   SymbolSet(symbolSet),
		Status:      Status(status),
		HQTFD:       HQTFD(hqtfd),
		Amplifier:   Amplifier(amplifier),
		Entity:      Entity(entity),
		Modifier1:   Modifier1(modifier1),
		Modifier2:   Modifier2(modifier2),
	}, nil
}
