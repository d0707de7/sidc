package app6d

import "fmt"

// ValidationError reports the first invalid field encountered when validating
// a SIDC. It identifies the field by name and includes a human-readable reason.
type ValidationError struct {
	Field  string
	Reason string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("sidc: invalid %s: %s", e.Field, e.Reason)
}

// Validate checks that every field carries a meaningful value:
//
//   - Version is a recognised D or E value (or Unspecified).
//   - Context, Affiliation, Status are within their enum ranges.
//   - HQTFD is 0-7.
//   - Amplifier, if non-zero, names a defined echelon, mobility, or leadership code.
//   - SymbolSet is a recognised value.
//   - If SymbolSet is E-only, Version is in the E family.
//   - The Entity, if non-zero, is defined for the SymbolSet.
//   - Modifier1 and Modifier2, if non-zero, are defined for the SymbolSet.
//
// Validate stops at the first failure and returns a *ValidationError naming
// the field. Parse does not call Validate — it only checks the structure of
// the input — so callers must invoke Validate explicitly when they want to
// know that the contents are meaningful.
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
	if s.Entity != 0 {
		if _, ok := entityNames[entityKey{Set: s.SymbolSet, E: s.Entity}]; !ok {
			return &ValidationError{Field: "Entity", Reason: fmt.Sprintf("%s is not defined in symbol set %s", s.Entity, s.SymbolSet)}
		}
	}
	if s.Modifier1 != 0 {
		if _, ok := modifier1Names[modifier1Key{Set: s.SymbolSet, M: s.Modifier1}]; !ok {
			return &ValidationError{Field: "Modifier1", Reason: fmt.Sprintf("%s is not defined in symbol set %s", s.Modifier1, s.SymbolSet)}
		}
	}
	if s.Modifier2 != 0 {
		if _, ok := modifier2Names[modifier2Key{Set: s.SymbolSet, M: s.Modifier2}]; !ok {
			return &ValidationError{Field: "Modifier2", Reason: fmt.Sprintf("%s is not defined in symbol set %s", s.Modifier2, s.SymbolSet)}
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
