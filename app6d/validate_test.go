package app6d

import (
	"errors"
	"testing"
)

func TestValidate(t *testing.T) {
	valid := SIDC{
		Version:     VersionD10,
		Context:     ContextReality,
		Affiliation: AffiliationFriend,
		SymbolSet:   SymbolSetLandUnit,
		Status:      StatusPresent,
		HQTFD:       HQTFDNone,
		Amplifier:   AmplifierPlatoonDetachment,
		Entity:      121100,
	}

	tests := []struct {
		name      string
		sidc      SIDC
		wantField string
	}{
		{name: "fully populated valid SIDC passes", sidc: valid, wantField: ""},
		{name: "zero value passes", sidc: SIDC{}, wantField: ""},
		{name: "unknown version is rejected", sidc: with(valid, func(s *SIDC) { s.Version = 99 }), wantField: "Version"},
		{name: "context above simulation is rejected", sidc: with(valid, func(s *SIDC) { s.Context = 5 }), wantField: "Context"},
		{name: "affiliation above hostile is rejected", sidc: with(valid, func(s *SIDC) { s.Affiliation = 9 }), wantField: "Affiliation"},
		{name: "status above full-to-capacity is rejected", sidc: with(valid, func(s *SIDC) { s.Status = 9 }), wantField: "Status"},
		{name: "HQTFD above 7 is rejected", sidc: with(valid, func(s *SIDC) { s.HQTFD = 8 }), wantField: "HQTFD"},
		{name: "amplifier in echelon range passes", sidc: with(valid, func(s *SIDC) { s.Amplifier = AmplifierBrigade }), wantField: ""},
		{name: "amplifier in mobility range passes", sidc: with(valid, func(s *SIDC) { s.Amplifier = AmplifierTracked }), wantField: ""},
		{name: "amplifier in leadership range passes", sidc: with(valid, func(s *SIDC) { s.Amplifier = AmplifierLeaderIndividual }), wantField: ""},
		{name: "amplifier in undefined range is rejected", sidc: with(valid, func(s *SIDC) { s.Amplifier = 99 }), wantField: "Amplifier"},
		{name: "unknown symbol set is rejected", sidc: with(valid, func(s *SIDC) { s.SymbolSet = 99 }), wantField: "SymbolSet"},
		{
			name:      "E-only symbol set in D version is rejected",
			sidc:      with(valid, func(s *SIDC) { s.SymbolSet = SymbolSetCyberspace; s.Version = VersionD10; s.Entity = 0 }),
			wantField: "SymbolSet",
		},
		{
			name:      "E-only symbol set in E version passes",
			sidc:      with(valid, func(s *SIDC) { s.SymbolSet = SymbolSetCyberspace; s.Version = VersionE13; s.Entity = 0 }),
			wantField: "",
		},
		{
			name:      "E-only symbol set with unspecified version passes",
			sidc:      with(valid, func(s *SIDC) { s.SymbolSet = SymbolSetCyberspace; s.Version = VersionUnspecified; s.Entity = 0 }),
			wantField: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.sidc.Validate()
			if tt.wantField == "" {
				if err != nil {
					t.Fatalf("expected nil error, got %v", err)
				}
				return
			}
			if err == nil {
				t.Fatalf("expected error for field %q, got nil", tt.wantField)
			}
			var ve *ValidationError
			if !errors.As(err, &ve) {
				t.Fatalf("expected *ValidationError, got %T: %v", err, err)
			}
			if ve.Field != tt.wantField {
				t.Errorf("got error for field %q, expected %q (full error: %v)", ve.Field, tt.wantField, err)
			}
		})
	}
}

func TestValidateStrict(t *testing.T) {
	tests := []struct {
		name    string
		sidc    SIDC
		wantErr error
	}{
		{
			name: "real Air fighter entity is accepted",
			sidc: SIDC{
				Version:   VersionD10,
				SymbolSet: SymbolSetAir,
				Entity:    EntityAir_MilitaryFixedWingFighter,
			},
		},
		{
			name: "real Land unit infantry entity is accepted",
			sidc: SIDC{
				Version:   VersionD10,
				SymbolSet: SymbolSetLandUnit,
				Entity:    121100,
			},
		},
		{
			name: "unknown entity in known symbol set is rejected",
			sidc: SIDC{
				Version:   VersionD10,
				SymbolSet: SymbolSetAir,
				Entity:    999999,
			},
			wantErr: ErrUnknownEntity,
		},
		{
			name: "zero entity is allowed",
			sidc: SIDC{
				Version:   VersionD10,
				SymbolSet: SymbolSetAir,
			},
		},
		{
			name: "unknown modifier1 is rejected",
			sidc: SIDC{
				Version:   VersionD10,
				SymbolSet: SymbolSetAir,
				Modifier1: 99,
			},
			wantErr: ErrUnknownModifier,
		},
		{
			name: "unknown modifier2 is rejected",
			sidc: SIDC{
				Version:   VersionD10,
				SymbolSet: SymbolSetAir,
				Modifier2: 99,
			},
			wantErr: ErrUnknownModifier,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.sidc.ValidateStrict()
			if tt.wantErr == nil {
				if err != nil {
					t.Fatalf("expected nil error, got %v", err)
				}
				return
			}
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("got %v, expected to wrap %v", err, tt.wantErr)
			}
		})
	}
}

func TestEntity_Name(t *testing.T) {
	tests := []struct {
		name   string
		set    SymbolSet
		entity Entity
		want   string
	}{
		{name: "fighter in air symbol set", set: SymbolSetAir, entity: EntityAir_MilitaryFixedWingFighter, want: "Military / Fixed Wing / Fighter"},
		{name: "unknown entity returns empty string", set: SymbolSetAir, entity: 999999, want: ""},
		{name: "same numeric entity code maps to different name in different symbol set", set: SymbolSetLandUnit, entity: 110000, want: "Command and Control"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.entity.Name(tt.set)
			if got != tt.want {
				t.Errorf("got %q, expected %q", got, tt.want)
			}
		})
	}
}

func TestModifierNames(t *testing.T) {
	if got := Modifier1(1).Name(SymbolSetAir); got == "" {
		t.Errorf("expected non-empty name for modifier1 1 in Air, got empty string")
	}
	if got := Modifier1(99).Name(SymbolSetAir); got != "" {
		t.Errorf("expected empty name for unknown modifier1 99 in Air, got %q", got)
	}
	if got := Modifier2(1).Name(SymbolSetAir); got == "" {
		t.Errorf("expected non-empty name for modifier2 1 in Air, got empty string")
	}
}

// with returns a copy of base with the mutator applied. Helper for table tests.
func with(base SIDC, mutate func(*SIDC)) SIDC {
	mutate(&base)
	return base
}
