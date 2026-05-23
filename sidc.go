// Package sidc provides the top-level Detect function for identifying which
// version of Symbol Identification Code a string represents. Parsing and
// building is done in the version-specific subpackages:
//
//   - github.com/d0707de7/sidc/app6b — 15-character letter-based SIDCs
//     (STANAG APP-6 B and C, MIL-STD-2525 B and C).
//   - github.com/d0707de7/sidc/app6d — 20-character number-based SIDCs
//     (STANAG APP-6 D and E, MIL-STD-2525 D and E).
//
// Most callers should call Detect on input of unknown origin, then dispatch
// to the appropriate package's Parse function.
package sidc

import (
	"errors"
	"fmt"
)

//go:generate go run ./internal/tsvgen -tables tables -out app6d

// Version identifies which SIDC standard a string belongs to.
type Version uint8

const (
	// VersionUnknown indicates the input did not match any known SIDC layout.
	VersionUnknown Version = iota
	// VersionAPP6B is the 15-character letter-based encoding used by
	// STANAG APP-6 B and C and MIL-STD-2525 B and C.
	VersionAPP6B
	// VersionAPP6D is the 20-character number-based encoding used by
	// STANAG APP-6 D and MIL-STD-2525 D (version digits 10-12).
	VersionAPP6D
	// VersionAPP6E is the 20-character number-based encoding used by
	// STANAG APP-6 E and MIL-STD-2525 E (version digits 13-14).
	VersionAPP6E
)

func (v Version) String() string {
	switch v {
	case VersionUnknown:
		return "Unknown"
	case VersionAPP6B:
		return "APP-6 B/C (letter-based, 15 chars)"
	case VersionAPP6D:
		return "APP-6 D (number-based, 20 chars)"
	case VersionAPP6E:
		return "APP-6 E (number-based, 20 chars)"
	}
	return fmt.Sprintf("Version(%d)", uint8(v))
}

// ErrUnknownVersion indicates Detect could not classify the input.
var ErrUnknownVersion = errors.New("sidc: input does not match any known SIDC layout")

// Detect returns the SIDC version that the input string belongs to, based on
// its length and leading characters. It is a cheap structural check — it does
// not parse or validate the body of the SIDC.
//
// Rules:
//   - 15 characters of any composition → APP-6 B (letter-based).
//   - 20 characters, all ASCII digits, leading "10", "11", or "12" → APP-6 D.
//   - 20 characters, all ASCII digits, leading "13" or "14" → APP-6 E.
//   - Anything else → VersionUnknown with ErrUnknownVersion.
func Detect(s string) (Version, error) {
	switch len(s) {
	case 15:
		return VersionAPP6B, nil
	case 20:
		for i := range len(s) {
			if s[i] < '0' || s[i] > '9' {
				return VersionUnknown, fmt.Errorf("%w: 20-char input contains non-digit at position %d", ErrUnknownVersion, i)
			}
		}
		switch s[0:2] {
		case "10", "11", "12":
			return VersionAPP6D, nil
		case "13", "14":
			return VersionAPP6E, nil
		default:
			return VersionUnknown, fmt.Errorf("%w: 20-digit input has unrecognised version prefix %q", ErrUnknownVersion, s[0:2])
		}
	default:
		return VersionUnknown, fmt.Errorf("%w: length %d is neither 15 nor 20", ErrUnknownVersion, len(s))
	}
}
