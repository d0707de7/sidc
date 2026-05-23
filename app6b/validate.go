package app6b

import "fmt"

// ValidationError reports the first invalid field encountered when validating
// a SIDC. It identifies the field by name and includes a human-readable reason.
type ValidationError struct {
	Field  string
	Reason string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("sidc/app6b: invalid %s: %s", e.Field, e.Reason)
}

// Validate checks that every field byte is printable ASCII (0x20-0x7E) or
// the unset zero byte. Letter-based SIDCs have a very permissive grammar at
// the byte level — any printable ASCII byte is syntactically legal — so
// Validate only catches the failure mode of "this byte cannot be rendered".
//
// Callers that want to check the bytes name a meaningful symbol per
// APP-6 B/C should consult the appropriate appendix and the FunctionID
// hierarchy; that's beyond the scope of this package, which models the
// string format rather than the standard's full semantic rules.
func (s SIDC) Validate() error {
	checks := []struct {
		field string
		b     byte
	}{
		{"CodingScheme", byte(s.CodingScheme)},
		{"Affiliation", byte(s.Affiliation)},
		{"BattleDimension", byte(s.BattleDimension)},
		{"Status", byte(s.Status)},
		{"SymbolModifier1", s.SymbolModifier1},
		{"SymbolModifier2", s.SymbolModifier2},
		{"OrderOfBattle", s.OrderOfBattle},
	}
	for _, c := range checks {
		if !isUnsetOrPrintable(c.b) {
			return &ValidationError{Field: c.field, Reason: fmt.Sprintf("byte 0x%02x is not printable ASCII", c.b)}
		}
	}
	for i, b := range s.FunctionID {
		if !isUnsetOrPrintable(b) {
			return &ValidationError{Field: fmt.Sprintf("FunctionID[%d]", i), Reason: fmt.Sprintf("byte 0x%02x is not printable ASCII", b)}
		}
	}
	for i, b := range s.CountryCode {
		if !isUnsetOrPrintable(b) {
			return &ValidationError{Field: fmt.Sprintf("CountryCode[%d]", i), Reason: fmt.Sprintf("byte 0x%02x is not printable ASCII", b)}
		}
	}
	return nil
}

func isUnsetOrPrintable(b byte) bool {
	return b == 0 || (b >= 0x20 && b <= 0x7E)
}
