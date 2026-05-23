package app6d

import (
	"errors"
	"testing"
)

func TestSIDC_String(t *testing.T) {
	tests := []struct {
		name     string
		sidc     SIDC
		expected string
	}{
		{
			name:     "zero value encodes as 20 zeros",
			sidc:     SIDC{},
			expected: "00000000000000000000",
		},
		{
			name: "APP-6 D version 10 land unit infantry encodes correctly",
			sidc: SIDC{
				Version:     VersionD10,
				Context:     ContextReality,
				Affiliation: AffiliationFriend,
				SymbolSet:   SymbolSetLandUnit,
				Status:      StatusPresent,
				HQTFD:       HQTFDNone,
				Amplifier:   AmplifierPlatoonDetachment,
				Entity:      121100,
				Modifier1:   0,
				Modifier2:   0,
			},
			expected: "10031000141211000000",
		},
		{
			name: "APP-6 E exercise hostile cyberspace encodes correctly",
			sidc: SIDC{
				Version:     VersionE13,
				Context:     ContextExercise,
				Affiliation: AffiliationHostile,
				SymbolSet:   SymbolSetCyberspace,
				Status:      StatusPresent,
				HQTFD:       HQTFDNone,
				Amplifier:   AmplifierNone,
				Entity:      110000,
				Modifier1:   0,
				Modifier2:   0,
			},
			expected: "13166000001100000000",
		},
		{
			name: "modifiers at the end are placed at positions 16-19",
			sidc: SIDC{
				Version:   VersionE13,
				SymbolSet: SymbolSetAir,
				Modifier1: 7,
				Modifier2: 42,
			},
			expected: "13000100000000000742",
		},
		{
			name: "HQTFD combination of headquarters and task force encodes as 6",
			sidc: SIDC{
				Version:   VersionD10,
				SymbolSet: SymbolSetLandUnit,
				HQTFD:     HQTFDHeadquartersTaskForce,
			},
			expected: "10001006000000000000",
		},
		{
			name: "every echelon amplifier value round trips",
			sidc: SIDC{
				Version:   VersionD10,
				SymbolSet: SymbolSetLandUnit,
				Amplifier: AmplifierCommand,
			},
			expected: "10001000260000000000",
		},
		{
			name: "entity uses six digits and pads with zeros",
			sidc: SIDC{
				Version:   VersionD10,
				SymbolSet: SymbolSetAir,
				Entity:    1,
			},
			expected: "10000100000000010000",
		},
		{
			name: "values exceeding field width wrap rather than panic",
			sidc: SIDC{
				Version: 99,
				Entity:  1234567,
			},
			expected: "99000000002345670000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.sidc.String()
			if got != tt.expected {
				t.Errorf("got %q, expected %q", got, tt.expected)
			}
			if len(got) != SIDCLength {
				t.Errorf("got length %d, expected %d", len(got), SIDCLength)
			}
		})
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected SIDC
	}{
		{
			name:     "20 zeros parses to zero value",
			input:    "00000000000000000000",
			expected: SIDC{},
		},
		{
			name:  "APP-6 D land unit infantry friend parses correctly",
			input: "10031000141211000000",
			expected: SIDC{
				Version:     VersionD10,
				Context:     ContextReality,
				Affiliation: AffiliationFriend,
				SymbolSet:   SymbolSetLandUnit,
				Status:      StatusPresent,
				HQTFD:       HQTFDNone,
				Amplifier:   AmplifierPlatoonDetachment,
				Entity:      121100,
				Modifier1:   0,
				Modifier2:   0,
			},
		},
		{
			name:  "APP-6 E exercise hostile cyberspace parses correctly",
			input: "13166000001100000000",
			expected: SIDC{
				Version:     VersionE13,
				Context:     ContextExercise,
				Affiliation: AffiliationHostile,
				SymbolSet:   SymbolSetCyberspace,
				Status:      StatusPresent,
				HQTFD:       HQTFDNone,
				Amplifier:   AmplifierNone,
				Entity:      110000,
				Modifier1:   0,
				Modifier2:   0,
			},
		},
		{
			name:  "modifiers at end are placed at the right positions",
			input: "13000100000000000742",
			expected: SIDC{
				Version:   VersionE13,
				SymbolSet: SymbolSetAir,
				Modifier1: 7,
				Modifier2: 42,
			},
		},
		{
			name:  "HQTFD bit values decode into the right combinations",
			input: "10001007000000000000",
			expected: SIDC{
				Version:   VersionD10,
				SymbolSet: SymbolSetLandUnit,
				HQTFD:     HQTFDFeintDummyHeadquartersTaskForce,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.expected {
				t.Errorf("got %+v, expected %+v", got, tt.expected)
			}
		})
	}
}

func TestParse_Errors(t *testing.T) {
	tests := []struct {
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
			name:        "19 characters fails length check",
			input:       "1003100014121100000",
			expectedErr: ErrInvalidLength,
		},
		{
			name:        "21 characters fails length check",
			input:       "100310001412110000000",
			expectedErr: ErrInvalidLength,
		},
		{
			name:        "letter in the middle fails character check",
			input:       "1003100014A211000000",
			expectedErr: ErrInvalidCharacter,
		},
		{
			name:        "hyphen at start fails character check",
			input:       "-0031000141211000000",
			expectedErr: ErrInvalidCharacter,
		},
		{
			name:        "space embedded fails character check",
			input:       "10031000 41211000000",
			expectedErr: ErrInvalidCharacter,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse(tt.input)
			if err == nil {
				t.Fatalf("expected error %v, got nil", tt.expectedErr)
			}
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("got error %v, expected to wrap %v", err, tt.expectedErr)
			}
		})
	}
}

func TestParseStringRoundTrip(t *testing.T) {
	inputs := []string{
		"00000000000000000000",
		"10031000141211000000",
		"13166000001100000000",
		"14025002000000007199",
		"12060006260000000000",
		"11041000351502010000",
		"13000100000000000742",
	}
	for _, input := range inputs {
		t.Run(input, func(t *testing.T) {
			parsed, err := Parse(input)
			if err != nil {
				t.Fatalf("unexpected error parsing %q: %v", input, err)
			}
			got := parsed.String()
			if got != input {
				t.Errorf("round trip changed value: parsed %q, got back %q", input, got)
			}
		})
	}
}
