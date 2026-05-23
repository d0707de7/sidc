package app6b

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
			name:     "zero value renders as 15 dashes",
			sidc:     SIDC{},
			expected: "---------------",
		},
		{
			name: "warfighting friend air fighter renders correctly",
			sidc: SIDC{
				CodingScheme:    CodingSchemeWarfighting,
				Affiliation:     AffiliationFriend,
				BattleDimension: BattleDimensionAir,
				Status:          StatusPresent,
				FunctionID:      FunctionID{'M', 'F', 'F', 0, 0, 0},
			},
			expected: "SFAPMFF--------",
		},
		{
			name: "country code is placed at positions 12 and 13",
			sidc: SIDC{
				CodingScheme:    CodingSchemeWarfighting,
				Affiliation:     AffiliationFriend,
				BattleDimension: BattleDimensionGround,
				Status:          StatusPresent,
				FunctionID:      FunctionID{'U', 'C', 'I', 0, 0, 0},
				CountryCode:     [2]byte{'U', 'S'},
			},
			expected: "SFGPUCI-----US-",
		},
		{
			name: "order of battle is placed at position 14",
			sidc: SIDC{
				CodingScheme:    CodingSchemeWarfighting,
				Affiliation:     AffiliationFriend,
				BattleDimension: BattleDimensionGround,
				Status:          StatusPresent,
				FunctionID:      FunctionID{'U', 'C', 'I', 0, 0, 0},
				OrderOfBattle:   'A',
			},
			expected: "SFGPUCI-------A",
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
			name:  "15 dashes parses to zero fields containing dash bytes",
			input: "---------------",
			expected: SIDC{
				CodingScheme:    '-',
				Affiliation:     '-',
				BattleDimension: '-',
				Status:          '-',
				FunctionID:      FunctionID{'-', '-', '-', '-', '-', '-'},
				SymbolModifier1: '-',
				SymbolModifier2: '-',
				CountryCode:     [2]byte{'-', '-'},
				OrderOfBattle:   '-',
			},
		},
		{
			name:  "warfighting friend air fighter parses correctly",
			input: "SFAPMFF--------",
			expected: SIDC{
				CodingScheme:    CodingSchemeWarfighting,
				Affiliation:     AffiliationFriend,
				BattleDimension: BattleDimensionAir,
				Status:          StatusPresent,
				FunctionID:      FunctionID{'M', 'F', 'F', '-', '-', '-'},
				SymbolModifier1: '-',
				SymbolModifier2: '-',
				CountryCode:     [2]byte{'-', '-'},
				OrderOfBattle:   '-',
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
		{name: "empty string fails length check", input: "", expectedErr: ErrInvalidLength},
		{name: "14 characters fails length check", input: "SFAPMFF-------", expectedErr: ErrInvalidLength},
		{name: "16 characters fails length check", input: "SFAPMFF---------", expectedErr: ErrInvalidLength},
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
		"---------------",
		"SFAPMFF--------",
		"SFGPUCI-----US-",
		"SFGPUCI-------A",
		"SHGPUCD-----RU-",
		"GFGPGTB--------",
	}
	for _, input := range inputs {
		t.Run(input, func(t *testing.T) {
			parsed, err := Parse(input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			got := parsed.String()
			if got != input {
				t.Errorf("round trip changed value: parsed %q, got back %q", input, got)
			}
		})
	}
}

func TestAffiliation_Family(t *testing.T) {
	tests := []struct {
		name        string
		affiliation Affiliation
		family      AffiliationFamily
		exercise    bool
	}{
		{name: "Friend is friend family, not exercise", affiliation: AffiliationFriend, family: AffiliationFamilyFriend},
		{name: "Hostile is hostile family, not exercise", affiliation: AffiliationHostile, family: AffiliationFamilyHostile},
		{name: "Neutral is neutral family, not exercise", affiliation: AffiliationNeutral, family: AffiliationFamilyNeutral},
		{name: "Unknown is unknown family, not exercise", affiliation: AffiliationUnknown, family: AffiliationFamilyUnknown},
		{name: "Joker is hostile family and exercise", affiliation: AffiliationJoker, family: AffiliationFamilyHostile, exercise: true},
		{name: "Faker is hostile family and exercise", affiliation: AffiliationFaker, family: AffiliationFamilyHostile, exercise: true},
		{name: "Exercise friend is friend family and exercise", affiliation: AffiliationExerciseFriend, family: AffiliationFamilyFriend, exercise: true},
		{name: "Suspect is hostile family", affiliation: AffiliationSuspect, family: AffiliationFamilyHostile},
		{name: "Assumed friend is friend family", affiliation: AffiliationAssumedFriend, family: AffiliationFamilyFriend},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.affiliation.Family(); got != tt.family {
				t.Errorf("Family() = %v, expected %v", got, tt.family)
			}
			if got := tt.affiliation.IsExercise(); got != tt.exercise {
				t.Errorf("IsExercise() = %v, expected %v", got, tt.exercise)
			}
		})
	}
}
