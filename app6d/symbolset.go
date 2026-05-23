package app6d

import "fmt"

// SymbolSet is the 2-digit symbol-set field at positions 4-5. It identifies
// which appendix of the standard governs the meaning of the Entity, Modifier1
// and Modifier2 fields.
type SymbolSet uint8

const (
	SymbolSetUnknown SymbolSet = 0

	SymbolSetAir SymbolSet = 1

	// SymbolSetAirMissile is only valid in APP-6 E.
	SymbolSetAirMissile SymbolSet = 2

	SymbolSetSpace SymbolSet = 5

	// SymbolSetSpaceMissile is only valid in APP-6 E.
	SymbolSetSpaceMissile SymbolSet = 6

	SymbolSetLandUnit         SymbolSet = 10
	SymbolSetLandCivilian     SymbolSet = 11
	SymbolSetLandEquipment    SymbolSet = 15
	SymbolSetLandInstallation SymbolSet = 20
	SymbolSetControlMeasure   SymbolSet = 25

	// SymbolSetDismountedIndividual is only valid in APP-6 E.
	SymbolSetDismountedIndividual SymbolSet = 27

	SymbolSetSeaSurface    SymbolSet = 30
	SymbolSetSeaSubsurface SymbolSet = 35
	SymbolSetMineWarfare   SymbolSet = 36
	SymbolSetActivities    SymbolSet = 40

	// SymbolSetSIGINTSpace is only valid in APP-6 E.
	SymbolSetSIGINTSpace SymbolSet = 50

	// SymbolSetSIGINTAir is only valid in APP-6 E.
	SymbolSetSIGINTAir SymbolSet = 51

	// SymbolSetSIGINTLand is only valid in APP-6 E.
	SymbolSetSIGINTLand SymbolSet = 52

	// SymbolSetSIGINTSeaSurface is only valid in APP-6 E.
	SymbolSetSIGINTSeaSurface SymbolSet = 53

	// SymbolSetSIGINTSubsurface is only valid in APP-6 E.
	SymbolSetSIGINTSubsurface SymbolSet = 54

	// SymbolSetCyberspace is only valid in APP-6 E.
	SymbolSetCyberspace SymbolSet = 60
)

// eOnly lists symbol sets that are only valid in APP-6 E.
var eOnlySymbolSets = map[SymbolSet]bool{
	SymbolSetAirMissile:           true,
	SymbolSetSpaceMissile:         true,
	SymbolSetDismountedIndividual: true,
	SymbolSetSIGINTSpace:          true,
	SymbolSetSIGINTAir:            true,
	SymbolSetSIGINTLand:           true,
	SymbolSetSIGINTSeaSurface:     true,
	SymbolSetSIGINTSubsurface:     true,
	SymbolSetCyberspace:           true,
}

// IsEOnly reports whether this symbol set is only valid in APP-6 E.
func (s SymbolSet) IsEOnly() bool { return eOnlySymbolSets[s] }

func (s SymbolSet) String() string {
	switch s {
	case SymbolSetUnknown:
		return "Unknown"
	case SymbolSetAir:
		return "Air"
	case SymbolSetAirMissile:
		return "Air missile"
	case SymbolSetSpace:
		return "Space"
	case SymbolSetSpaceMissile:
		return "Space missile"
	case SymbolSetLandUnit:
		return "Land unit"
	case SymbolSetLandCivilian:
		return "Land civilian unit/Organization"
	case SymbolSetLandEquipment:
		return "Land equipment"
	case SymbolSetLandInstallation:
		return "Land installations"
	case SymbolSetControlMeasure:
		return "Control measure"
	case SymbolSetDismountedIndividual:
		return "Dismounted individual"
	case SymbolSetSeaSurface:
		return "Sea surface"
	case SymbolSetSeaSubsurface:
		return "Sea subsurface"
	case SymbolSetMineWarfare:
		return "Mine warfare"
	case SymbolSetActivities:
		return "Activity/Event"
	case SymbolSetSIGINTSpace:
		return "Signals intelligence (space)"
	case SymbolSetSIGINTAir:
		return "Signals intelligence (air)"
	case SymbolSetSIGINTLand:
		return "Signals intelligence (land)"
	case SymbolSetSIGINTSeaSurface:
		return "Signals intelligence (sea surface)"
	case SymbolSetSIGINTSubsurface:
		return "Signals intelligence (subsurface)"
	case SymbolSetCyberspace:
		return "Cyberspace"
	default:
		return fmt.Sprintf("SymbolSet(%02d)", uint8(s))
	}
}
