package app6d

import (
	"errors"
	"fmt"
)

// ValidationError reports the first invalid field encountered when validating
// a SIDC. It identifies the field by name and includes a human-readable reason.
type ValidationError struct {
	Field  string
	Reason string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("sidc: invalid %s: %s", e.Field, e.Reason)
}

// ErrUnknownEntity indicates the entity code is not defined for the symbol set.
var ErrUnknownEntity = errors.New("sidc: unknown entity for symbol set")

// ErrUnknownModifier indicates a modifier code is not defined for the symbol set.
var ErrUnknownModifier = errors.New("sidc: unknown modifier for symbol set")

// Validate performs cross-field structural and table-lookup checks. It does
// not require the SIDC to identify a known entity (zero is treated as "not
// specified"); pass StrictEntity to require entity membership in the symbol
// set's table.
//
// Checks performed:
//   - Version must be a recognised D or E value (or Unspecified).
//   - Context, Affiliation, Status must be in their enum ranges.
//   - HQTFD must be 0-7.
//   - Amplifier, if non-zero, must fall in one of the echelon, mobility,
//     or leadership ranges (Category() != AmplifierCategoryNone).
//   - SymbolSet must be a recognised value.
//   - If SymbolSet is E-only, Version must be in the E family.
func (s SIDC) Validate() error {
	if s.Version != VersionUnspecified && !s.Version.IsD() && !s.Version.IsE() {
		return &ValidationError{Field: "Version", Reason: fmt.Sprintf("unrecognised value %d", uint8(s.Version))}
	}
	if s.Context > ContextSimulation {
		return &ValidationError{Field: "Context", Reason: fmt.Sprintf("out of range value %d", uint8(s.Context))}
	}
	if s.Affiliation > AffiliationHostile {
		return &ValidationError{Field: "Affiliation", Reason: fmt.Sprintf("out of range value %d", uint8(s.Affiliation))}
	}
	if s.Status > StatusFullToCapacity {
		return &ValidationError{Field: "Status", Reason: fmt.Sprintf("out of range value %d", uint8(s.Status))}
	}
	if s.HQTFD > HQTFDFeintDummyHeadquartersTaskForce {
		return &ValidationError{Field: "HQTFD", Reason: fmt.Sprintf("out of range value %d", uint8(s.HQTFD))}
	}
	if s.Amplifier != AmplifierNone && s.Amplifier.Category() == AmplifierCategoryNone {
		return &ValidationError{Field: "Amplifier", Reason: fmt.Sprintf("value %d is not a defined echelon, mobility, or leadership code", uint8(s.Amplifier))}
	}
	if !isKnownSymbolSet(s.SymbolSet) {
		return &ValidationError{Field: "SymbolSet", Reason: fmt.Sprintf("unrecognised value %d", uint8(s.SymbolSet))}
	}
	if s.SymbolSet.IsEOnly() && s.Version != VersionUnspecified && !s.Version.IsE() {
		return &ValidationError{Field: "SymbolSet", Reason: fmt.Sprintf("%s is only valid in APP-6 E, got version %d", s.SymbolSet, uint8(s.Version))}
	}
	return nil
}

// ValidateStrict performs all the checks of Validate plus table lookups: the
// Entity must exist in the symbol set, and any non-zero modifiers must too.
func (s SIDC) ValidateStrict() error {
	if err := s.Validate(); err != nil {
		return err
	}
	if s.Entity != 0 {
		if _, ok := entityNames[entityKey{Set: s.SymbolSet, E: s.Entity}]; !ok {
			return fmt.Errorf("%w: entity %s in symbol set %s", ErrUnknownEntity, s.Entity, s.SymbolSet)
		}
	}
	if s.Modifier1 != 0 {
		if _, ok := modifier1Names[modifier1Key{Set: s.SymbolSet, M: s.Modifier1}]; !ok {
			return fmt.Errorf("%w: modifier1 %s in symbol set %s", ErrUnknownModifier, s.Modifier1, s.SymbolSet)
		}
	}
	if s.Modifier2 != 0 {
		if _, ok := modifier2Names[modifier2Key{Set: s.SymbolSet, M: s.Modifier2}]; !ok {
			return fmt.Errorf("%w: modifier2 %s in symbol set %s", ErrUnknownModifier, s.Modifier2, s.SymbolSet)
		}
	}
	return nil
}

func isKnownSymbolSet(s SymbolSet) bool {
	switch s {
	case SymbolSetUnknown,
		SymbolSetAir, SymbolSetAirMissile,
		SymbolSetSpace, SymbolSetSpaceMissile,
		SymbolSetLandUnit, SymbolSetLandCivilian, SymbolSetLandEquipment, SymbolSetLandInstallation,
		SymbolSetControlMeasure, SymbolSetDismountedIndividual,
		SymbolSetSeaSurface, SymbolSetSeaSubsurface, SymbolSetMineWarfare,
		SymbolSetActivities,
		SymbolSetSIGINTSpace, SymbolSetSIGINTAir, SymbolSetSIGINTLand,
		SymbolSetSIGINTSeaSurface, SymbolSetSIGINTSubsurface,
		SymbolSetCyberspace:
		return true
	}
	return false
}

// Name returns the hierarchical name of this entity in the given symbol set,
// e.g. "Military / Fixed Wing / Fighter". Returns the empty string if not found.
func (e Entity) Name(set SymbolSet) string {
	return entityNames[entityKey{Set: set, E: e}]
}

// Name returns the modifier 1 name for the given symbol set, or "" if not found.
func (m Modifier1) Name(set SymbolSet) string {
	return modifier1Names[modifier1Key{Set: set, M: m}]
}

// Name returns the modifier 2 name for the given symbol set, or "" if not found.
func (m Modifier2) Name(set SymbolSet) string {
	return modifier2Names[modifier2Key{Set: set, M: m}]
}
