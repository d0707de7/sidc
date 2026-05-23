package app6d

import (
	"errors"
	"strings"
	"testing"
)

// landUnitInfantryPlatoon is a fully-populated valid SIDC used across tests as a
// realistic example: an APP-6 D, friendly, present, platoon-sized infantry unit.
var landUnitInfantryPlatoon = SIDC{
	Version:     VersionD10,
	Context:     ContextReality,
	Affiliation: AffiliationFriend,
	SymbolSet:   SymbolSetLandUnit,
	Status:      StatusPresent,
	HQTFD:       HQTFDNone,
	Amplifier:   AmplifierPlatoonDetachment,
	Entity:      EntityLandUnit_MovementAndManeuverInfantry,
}

// airFighterWithModifiers shows how modifier constants slot into the SIDC.
// The Air symbol set's modifier 1 sector identifies the air vehicle type;
// modifier 2 identifies a re-fuelling capability for tankers.
var airFighterWithModifiers = SIDC{
	Version:     VersionE13,
	Context:     ContextReality,
	Affiliation: AffiliationFriend,
	SymbolSet:   SymbolSetAir,
	Entity:      EntityAir_MilitaryFixedWingFighter,
	Modifier1:   Modifier1Air_Fighter,
	Modifier2:   Modifier2Air_BoomOnly,
}

// mustValue is a helper used in tests to render a SIDC that the test author
// knows is valid. Use it only after constructing from constants the test
// expects to be valid; if Validate fails, the test fails loudly.
func mustValue(t *testing.T, s SIDC) string {
	t.Helper()
	v, err := s.Value()
	if err != nil {
		t.Fatalf("expected valid SIDC, got Validate error: %v", err)
	}
	return v
}

func TestSIDC_Value(t *testing.T) {
	tests := []struct {
		name string
		sidc SIDC
		// expected is the one place we anchor the encoding to specific
		// characters. It is the canonical reference for what each field
		// position looks like.
		expected string
	}{
		{
			name:     "zero value encodes as 20 zeros",
			sidc:     SIDC{},
			expected: strings.Repeat("0", SIDCLength),
		},
		{
			name:     "landUnitInfantryPlatoon encodes per the field layout",
			sidc:     landUnitInfantryPlatoon,
			expected: "10031000141211000000",
		},
		{
			name:     "airFighterWithModifiers places modifiers in the last four positions",
			sidc:     airFighterWithModifiers,
			expected: "13030100001101040404",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.sidc.Value()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.expected {
				t.Errorf("got %q, expected %q", got, tt.expected)
			}
			if len(got) != SIDCLength {
				t.Errorf("got length %d, expected %d", len(got), SIDCLength)
			}
		})
	}
}

func TestSIDC_Value_RefusesInvalid(t *testing.T) {
	cases := []struct {
		name      string
		sidc      SIDC
		wantField string
	}{
		{name: "unknown version", sidc: mutate(landUnitInfantryPlatoon, func(s *SIDC) { s.Version = 99 }), wantField: "Version"},
		{name: "unknown context", sidc: mutate(landUnitInfantryPlatoon, func(s *SIDC) { s.Context = 5 }), wantField: "Context"},
		{name: "Entity above six digits is rejected (as unknown entity for symbol set)", sidc: SIDC{Entity: 1_000_000}, wantField: "Entity"},
		{name: "Modifier1 above two digits", sidc: mutate(landUnitInfantryPlatoon, func(s *SIDC) { s.Modifier1 = 100 }), wantField: "Modifier1"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.sidc.Value()
			if err == nil {
				t.Fatalf("expected Value to refuse invalid SIDC, got %q", got)
			}
			if got != "" {
				t.Errorf("expected empty string on error, got %q", got)
			}
			var ve *ValidationError
			if !errors.As(err, &ve) {
				t.Fatalf("expected *ValidationError, got %T: %v", err, err)
			}
			if ve.Field != tc.wantField {
				t.Errorf("got error for field %q, expected %q", ve.Field, tc.wantField)
			}
		})
	}
}

// TestParse_RoundTrip uses Value to produce a canonical 20-digit string from
// a SIDC built from constants, parses it, and asserts equality. This proves
// Parse and Value round-trip without hardcoding 20-digit magic strings.
func TestParse_RoundTrip(t *testing.T) {
	cases := []struct {
		name string
		sidc SIDC
	}{
		{name: "zero value", sidc: SIDC{}},
		{name: "land unit infantry platoon", sidc: landUnitInfantryPlatoon},
		{name: "air fighter with modifiers", sidc: airFighterWithModifiers},
		{
			name: "every echelon amplifier value",
			sidc: SIDC{
				Version:   VersionD10,
				SymbolSet: SymbolSetLandUnit,
				Amplifier: AmplifierCommand,
			},
		},
		{
			name: "every HQTFD bit combined",
			sidc: SIDC{
				Version:   VersionD10,
				SymbolSet: SymbolSetLandUnit,
				HQTFD:     HQTFDFeintDummyHeadquartersTaskForce,
			},
		},
		{
			name: "APP-6 E cyberspace unit",
			sidc: SIDC{
				Version:     VersionE13,
				Context:     ContextExercise,
				Affiliation: AffiliationHostile,
				SymbolSet:   SymbolSetCyberspace,
				Entity:      EntityCyberspace_CyberspaceUnitCyberspaceUnitNonSpecified,
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			encoded := mustValue(t, tc.sidc)
			parsed, err := Parse(encoded)
			if err != nil {
				t.Fatalf("Parse(%q) returned unexpected error: %v", encoded, err)
			}
			if parsed != tc.sidc {
				t.Errorf("round trip mismatch:\n  original: %+v\n  parsed:   %+v\n  encoded:  %q", tc.sidc, parsed, encoded)
			}
		})
	}
}

func TestParse_RejectsInvalidStructure(t *testing.T) {
	cases := []struct {
		name        string
		input       string
		expectedErr error
	}{
		{
			name:        "empty string fails length check",
			input:       "",
			expectedErr: ErrInvalidLength,
		},
		{
			name:        "one short of SIDCLength fails length check",
			input:       strings.Repeat("0", SIDCLength-1),
			expectedErr: ErrInvalidLength,
		},
		{
			name:        "one over SIDCLength fails length check",
			input:       strings.Repeat("0", SIDCLength+1),
			expectedErr: ErrInvalidLength,
		},
		{
			name:        "letter at the middle fails character check",
			input:       "1003100014A211000000",
			expectedErr: ErrInvalidCharacter,
		},
		{
			name:        "hyphen at start fails character check",
			input:       "-" + strings.Repeat("0", SIDCLength-1),
			expectedErr: ErrInvalidCharacter,
		},
		{
			name:        "space embedded fails character check",
			input:       "10031000 41211000000",
			expectedErr: ErrInvalidCharacter,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Parse(tc.input)
			if err == nil {
				t.Fatalf("expected error %v, got nil", tc.expectedErr)
			}
			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("got error %v, expected to wrap %v", err, tc.expectedErr)
			}
		})
	}
}

// TestParse_AcceptsStructurallyValidButSemanticallyInvalid documents the
// split: Parse accepts any 20 digits, even if the resulting field values are
// not meaningful. Validate is the gate for that. This is deliberate so a
// caller can inspect or repair a malformed SIDC rather than losing it at the
// Parse step.
func TestParse_AcceptsStructurallyValidButSemanticallyInvalid(t *testing.T) {
	cases := []struct {
		name      string
		input     string
		wantField string
	}{
		{name: "version 99 is structurally fine, semantically invalid", input: "99031000141211000000", wantField: "Version"},
		{name: "context 5 is structurally fine, semantically invalid", input: "10531000141211000000", wantField: "Context"},
		{name: "affiliation 9 is structurally fine, semantically invalid", input: "10091000141211000000", wantField: "Affiliation"},
		{name: "symbol set 99 is structurally fine, semantically invalid", input: "10039900141211000000", wantField: "SymbolSet"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			parsed, err := Parse(tc.input)
			if err != nil {
				t.Fatalf("Parse(%q) returned unexpected error: %v", tc.input, err)
			}

			err = parsed.Validate()
			if err == nil {
				t.Fatalf("Validate accepted semantically invalid SIDC, expected error for field %q", tc.wantField)
			}
			var ve *ValidationError
			if !errors.As(err, &ve) {
				t.Fatalf("expected *ValidationError, got %T: %v", err, err)
			}
			if ve.Field != tc.wantField {
				t.Errorf("got error for field %q, expected %q", ve.Field, tc.wantField)
			}
		})
	}
}

// mutate is a tiny helper to copy a SIDC and tweak one field, used to build
// test inputs from a known-valid baseline without restating every field.
func mutate(base SIDC, f func(*SIDC)) SIDC {
	f(&base)
	return base
}
