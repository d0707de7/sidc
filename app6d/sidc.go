package app6d

import (
	"errors"
	"fmt"
	"strconv"
)

// SIDC is a parsed 20-digit Symbol Identification Code for APP-6 D or E.
//
// SIDC deliberately does not implement fmt.Stringer. Rendering a SIDC to
// its string form can fail (the field values must name a real symbol), and
// the API surfaces that by requiring callers to use the Value method, which
// returns (string, error). This means fmt.Println(s) and similar print the
// struct dump rather than a SIDC string; that's intentional — there is no
// safe way to render an arbitrary composite literal without first validating
// it.
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

// Value returns the canonical 20-digit SIDC representation. It calls
// Validate first; if Validate returns an error, Value returns that error
// and the empty string. On success the result is guaranteed to be exactly
// 20 ASCII digits.
//
// Use Value (not a String method, which doesn't exist) any time you want
// the wire-format SIDC. The two-result signature exists so that callers
// cannot accidentally ship an invalid SIDC.
func (s SIDC) Value() (string, error) {
	if err := s.Validate(); err != nil {
		return "", err
	}
	return s.render(), nil
}

// MarshalText implements encoding.TextMarshaler. It calls Value, so encoding
// a SIDC to JSON, XML, or any other text-based format will fail if the SIDC
// is not valid. encoding/json and encoding/xml both honour this interface,
// so no separate MarshalJSON is needed.
func (s SIDC) MarshalText() ([]byte, error) {
	v, err := s.Value()
	if err != nil {
		return nil, err
	}
	return []byte(v), nil
}

// UnmarshalText implements encoding.TextUnmarshaler. The input must be a
// structurally valid SIDC (Parse). Callers that want stricter semantic
// validation should call Validate on the result.
func (s *SIDC) UnmarshalText(text []byte) error {
	parsed, err := Parse(string(text))
	if err != nil {
		return err
	}
	*s = parsed
	return nil
}

// render produces the canonical encoding without validating. It is the raw
// rendering used by Value once validation has passed. Kept unexported so
// that the only public path to a SIDC string forces validation.
func (s SIDC) render() string {
	return fmt.Sprintf(
		"%02d%01d%01d%02d%01d%01d%02d%06d%02d%02d",
		uint8(s.Version),
		uint8(s.Context),
		uint8(s.Affiliation),
		uint8(s.SymbolSet),
		uint8(s.Status),
		uint8(s.HQTFD),
		uint8(s.Amplifier),
		uint32(s.Entity),
		uint8(s.Modifier1),
		uint8(s.Modifier2),
	)
}

// Parse parses a 20-digit SIDC string into its component fields. It checks
// the structure of the input — length, and that every byte is an ASCII digit —
// but does not check that field values are meaningful. Call Validate on the
// result to check that enum values are in range, that the entity is defined
// for its symbol set, and so on.
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
