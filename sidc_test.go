package sidc

import (
	"strings"
	"testing"
)

func TestDetect(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantVersion Version
		wantOK      bool
	}{
		{name: "15 characters is APP-6 B", input: "SFAPMFF--------", wantVersion: VersionAPP6B, wantOK: true},
		{name: "15 dashes is APP-6 B", input: strings.Repeat("-", 15), wantVersion: VersionAPP6B, wantOK: true},
		{name: "20 digits starting 10 is APP-6 D", input: "10031000141211000000", wantVersion: VersionAPP6D, wantOK: true},
		{name: "20 digits starting 11 is APP-6 D", input: "11031000141211000000", wantVersion: VersionAPP6D, wantOK: true},
		{name: "20 digits starting 12 is APP-6 D", input: "12031000141211000000", wantVersion: VersionAPP6D, wantOK: true},
		{name: "20 digits starting 13 is APP-6 E", input: "13166000001100000000", wantVersion: VersionAPP6E, wantOK: true},
		{name: "20 digits starting 14 is APP-6 E", input: "14025002000000007199", wantVersion: VersionAPP6E, wantOK: true},
		{name: "20 digits with unknown version prefix is unrecognised", input: "99031000141211000000", wantVersion: VersionUnknown, wantOK: false},
		{name: "20 chars with non-digit is unrecognised", input: "1003100014A211000000", wantVersion: VersionUnknown, wantOK: false},
		{name: "empty string is unrecognised", input: "", wantVersion: VersionUnknown, wantOK: false},
		{name: "10 characters is unrecognised", input: strings.Repeat("0", 10), wantVersion: VersionUnknown, wantOK: false},
		{name: "30 characters is unrecognised", input: strings.Repeat("0", 30), wantVersion: VersionUnknown, wantOK: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			version, ok := Detect(tt.input)
			if version != tt.wantVersion {
				t.Errorf("got version %v, expected %v", version, tt.wantVersion)
			}
			if ok != tt.wantOK {
				t.Errorf("got ok=%v, expected %v", ok, tt.wantOK)
			}
		})
	}
}

func TestVersion_String(t *testing.T) {
	tests := []struct {
		version Version
		want    string
	}{
		{version: VersionUnknown, want: "Unknown"},
		{version: VersionAPP6B, want: "APP-6 B/C (letter-based, 15 chars)"},
		{version: VersionAPP6D, want: "APP-6 D (number-based, 20 chars)"},
		{version: VersionAPP6E, want: "APP-6 E (number-based, 20 chars)"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.version.String(); got != tt.want {
				t.Errorf("got %q, expected %q", got, tt.want)
			}
		})
	}
}
