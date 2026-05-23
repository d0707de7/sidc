package app6d

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
		{name: "zero value", sidc: SIDC{}, expected: strings.Repeat("0", SIDCLength)},
		{name: "land unit infantry platoon", sidc: landUnitInfantryPlatoon, expected: "10031000141211000000"},
		{name: "air fighter with modifiers", sidc: airFighterWithModifiers, expected: "13030100001101040404"},
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
	invalid := mutate(landUnitInfantryPlatoon, func(s *SIDC) { s.Version = 99 })
	got, err := invalid.MarshalText()
	if err == nil {
		t.Fatalf("expected MarshalText to refuse invalid SIDC, got %q", got)
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
	if err := s.UnmarshalText([]byte("10031000141211000000")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s != landUnitInfantryPlatoon {
		t.Errorf("UnmarshalText produced %+v, expected %+v", s, landUnitInfantryPlatoon)
	}
}

func TestUnmarshalText_RejectsInvalidStructure(t *testing.T) {
	var s SIDC
	err := s.UnmarshalText([]byte("not a SIDC"))
	if err == nil {
		t.Fatal("expected UnmarshalText to reject invalid input, got nil error")
	}
}

// TestJSONMarshal verifies the JSON encoder picks up MarshalText and emits
// the SIDC as a JSON string, not as a struct object.
func TestJSONMarshal(t *testing.T) {
	got, err := json.Marshal(landUnitInfantryPlatoon)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := `"10031000141211000000"`
	if string(got) != want {
		t.Errorf("got %s, expected %s", got, want)
	}
}

func TestJSONMarshal_RefusesInvalid(t *testing.T) {
	invalid := mutate(landUnitInfantryPlatoon, func(s *SIDC) { s.Version = 99 })
	_, err := json.Marshal(invalid)
	if err == nil {
		t.Fatal("expected json.Marshal to refuse invalid SIDC, got nil error")
	}
}

func TestJSONUnmarshal(t *testing.T) {
	var s SIDC
	if err := json.Unmarshal([]byte(`"10031000141211000000"`), &s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s != landUnitInfantryPlatoon {
		t.Errorf("json.Unmarshal produced %+v, expected %+v", s, landUnitInfantryPlatoon)
	}
}

// TestJSONRoundTrip exercises a typical embedding case: a SIDC field inside
// a larger struct, marshalled and unmarshalled back together.
func TestJSONRoundTrip(t *testing.T) {
	type message struct {
		Name string `json:"name"`
		SIDC SIDC   `json:"sidc"`
	}
	original := message{Name: "alpha", SIDC: landUnitInfantryPlatoon}
	encoded, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}
	wantSubstring := `"sidc":"10031000141211000000"`
	if !strings.Contains(string(encoded), wantSubstring) {
		t.Errorf("JSON output %q did not contain %q", encoded, wantSubstring)
	}

	var decoded message
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	if decoded != original {
		t.Errorf("round trip mismatch:\n  original: %+v\n  decoded:  %+v", original, decoded)
	}
}

func TestXMLRoundTrip(t *testing.T) {
	type message struct {
		XMLName xml.Name `xml:"message"`
		SIDC    SIDC     `xml:"sidc"`
	}
	original := message{SIDC: landUnitInfantryPlatoon}
	encoded, err := xml.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}
	wantSubstring := "<sidc>10031000141211000000</sidc>"
	if !strings.Contains(string(encoded), wantSubstring) {
		t.Errorf("XML output %q did not contain %q", encoded, wantSubstring)
	}

	var decoded message
	if err := xml.Unmarshal(encoded, &decoded); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	if decoded.SIDC != original.SIDC {
		t.Errorf("round trip mismatch: original %+v, decoded %+v", original.SIDC, decoded.SIDC)
	}
}
