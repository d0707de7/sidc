package sidc

import (
	"errors"
	"testing"
)

func TestDetect(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Version
		wantErr bool
	}{
		{name: "15 characters is APP-6 B", input: "SFAPMFF--------", want: VersionAPP6B},
		{name: "15 dashes is APP-6 B", input: "---------------", want: VersionAPP6B},
		{name: "20 digits starting 10 is APP-6 D", input: "10031000141211000000", want: VersionAPP6D},
		{name: "20 digits starting 11 is APP-6 D", input: "11031000141211000000", want: VersionAPP6D},
		{name: "20 digits starting 12 is APP-6 D", input: "12031000141211000000", want: VersionAPP6D},
		{name: "20 digits starting 13 is APP-6 E", input: "13166000001100000000", want: VersionAPP6E},
		{name: "20 digits starting 14 is APP-6 E", input: "14025002000000007199", want: VersionAPP6E},
		{name: "20 digits with unknown version prefix is rejected", input: "99031000141211000000", want: VersionUnknown, wantErr: true},
		{name: "20 chars with non-digit is rejected", input: "1003100014A211000000", want: VersionUnknown, wantErr: true},
		{name: "empty string is rejected", input: "", want: VersionUnknown, wantErr: true},
		{name: "10 characters is rejected", input: "1234567890", want: VersionUnknown, wantErr: true},
		{name: "30 characters is rejected", input: "100310001412110000000000000000", want: VersionUnknown, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Detect(tt.input)
			if got != tt.want {
				t.Errorf("got version %v, expected %v", got, tt.want)
			}
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if !errors.Is(err, ErrUnknownVersion) {
					t.Errorf("got error %v, expected to wrap ErrUnknownVersion", err)
				}
				return
			}
			if err != nil {
				t.Errorf("got unexpected error: %v", err)
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
