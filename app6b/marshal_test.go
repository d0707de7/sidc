package app6b

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"strings"
	"testing"
)

func TestMarshalText(t *testing.T) {
	tests := []struct {
		name     string
		sidc     SIDC
		expected string
	}{
		{
			name:     "zero value",
			sidc:     SIDC{},
			expected: "---------------",
		},
		{
			name: "warfighting friend air fighter",
			sidc: SIDC{
				CodingScheme:    CodingSchemeWarfighting,
				Affiliation:     AffiliationFriend,
				BattleDimension: BattleDimensionAir,
				Status:          StatusPresent,
				FunctionID:      FunctionID{'M', 'F', 'F', 0, 0, 0},
			},
			expected: "SFAPMFF--------",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.sidc.MarshalText()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if string(got) != tt.expected {
				t.Errorf("got %q, expected %q", got, tt.expected)
			}
		})
	}
}

func TestMarshalText_RefusesInvalid(t *testing.T) {
	invalid := SIDC{CodingScheme: 0x01}
	got, err := invalid.MarshalText()
	if err == nil {
		t.Fatalf("expected MarshalText to refuse SIDC with non-printable byte, got %q", got)
	}
	if got != nil {
		t.Errorf("expected nil bytes on error, got %q", got)
	}
	var ve *ValidationError
	if !errors.As(err, &ve) {
		t.Fatalf("expected *ValidationError, got %T: %v", err, err)
	}
}

func TestUnmarshalText(t *testing.T) {
	var s SIDC
	if err := s.UnmarshalText([]byte("SFAPMFF--------")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := SIDC{
		CodingScheme:    CodingSchemeWarfighting,
		Affiliation:     AffiliationFriend,
		BattleDimension: BattleDimensionAir,
		Status:          StatusPresent,
		FunctionID:      FunctionID{'M', 'F', 'F', '-', '-', '-'},
		SymbolModifier1: '-',
		SymbolModifier2: '-',
		CountryCode:     [2]byte{'-', '-'},
		OrderOfBattle:   '-',
	}
	if s != want {
		t.Errorf("UnmarshalText produced %+v, expected %+v", s, want)
	}
}

func TestJSONMarshal(t *testing.T) {
	got, err := json.Marshal(SIDC{
		CodingScheme:    CodingSchemeWarfighting,
		Affiliation:     AffiliationFriend,
		BattleDimension: BattleDimensionAir,
		Status:          StatusPresent,
		FunctionID:      FunctionID{'M', 'F', 'F', 0, 0, 0},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := `"SFAPMFF--------"`
	if string(got) != want {
		t.Errorf("got %s, expected %s", got, want)
	}
}

func TestJSONMarshal_RefusesInvalid(t *testing.T) {
	invalid := SIDC{CodingScheme: 0x01}
	_, err := json.Marshal(invalid)
	if err == nil {
		t.Fatal("expected json.Marshal to refuse invalid SIDC, got nil error")
	}
}

func TestJSONRoundTrip(t *testing.T) {
	type message struct {
		Name string `json:"name"`
		SIDC SIDC   `json:"sidc"`
	}
	original := message{
		Name: "alpha",
		SIDC: SIDC{
			CodingScheme:    CodingSchemeWarfighting,
			Affiliation:     AffiliationFriend,
			BattleDimension: BattleDimensionAir,
			Status:          StatusPresent,
			FunctionID:      FunctionID{'M', 'F', 'F', 0, 0, 0},
		},
	}
	encoded, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}
	wantSubstring := `"sidc":"SFAPMFF--------"`
	if !strings.Contains(string(encoded), wantSubstring) {
		t.Errorf("JSON output %q did not contain %q", encoded, wantSubstring)
	}

	var decoded message
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	// Note: decoded.SIDC has '-' bytes where original has 0 (Parse normalised
	// the placeholders), so direct equality won't hold. Compare the rendered
	// form instead.
	originalEncoded, _ := original.SIDC.Value()
	decodedEncoded, _ := decoded.SIDC.Value()
	if originalEncoded != decodedEncoded {
		t.Errorf("round trip mismatch:\n  original: %q\n  decoded:  %q", originalEncoded, decodedEncoded)
	}
}

func TestXMLRoundTrip(t *testing.T) {
	type message struct {
		XMLName xml.Name `xml:"message"`
		SIDC    SIDC     `xml:"sidc"`
	}
	original := message{SIDC: SIDC{
		CodingScheme:    CodingSchemeWarfighting,
		Affiliation:     AffiliationFriend,
		BattleDimension: BattleDimensionAir,
		Status:          StatusPresent,
		FunctionID:      FunctionID{'M', 'F', 'F', 0, 0, 0},
	}}
	encoded, err := xml.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}
	wantSubstring := "<sidc>SFAPMFF--------</sidc>"
	if !strings.Contains(string(encoded), wantSubstring) {
		t.Errorf("XML output %q did not contain %q", encoded, wantSubstring)
	}

	var decoded message
	if err := xml.Unmarshal(encoded, &decoded); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	originalEncoded, _ := original.SIDC.Value()
	decodedEncoded, _ := decoded.SIDC.Value()
	if originalEncoded != decodedEncoded {
		t.Errorf("round trip mismatch: original %q, decoded %q", originalEncoded, decodedEncoded)
	}
}
