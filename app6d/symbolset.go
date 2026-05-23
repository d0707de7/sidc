package app6d

import "fmt"

// SymbolSet is the 2-digit symbol-set field at positions 4-5. It identifies
// which appendix of the standard governs the meaning of the Entity, Modifier1
// and Modifier2 fields.
type SymbolSet uint8

const (
	SymbolSetUnknown              SymbolSet = 0
	SymbolSetAir                  SymbolSet = 1  // 01
	SymbolSetAirMissile           SymbolSet = 2  // 02 (E only)
	SymbolSetSpace                SymbolSet = 5  // 05
	SymbolSetSpaceMissile         SymbolSet = 6  // 06 (E only)
	SymbolSetLandUnit             SymbolSet = 10 // 10
	SymbolSetLandCivilian         SymbolSet = 11 // 11
	SymbolSetLandEquipment        SymbolSet = 15 // 15
	SymbolSetLandInstallation     SymbolSet = 20 // 20
	SymbolSetControlMeasure       SymbolSet = 25 // 25
	SymbolSetDismountedIndividual SymbolSet = 27 // 27 (E only)
	SymbolSetSeaSurface           SymbolSet = 30 // 30
	SymbolSetSeaSubsurface        SymbolSet = 35 // 35
	SymbolSetMineWarfare          SymbolSet = 36 // 36
	SymbolSetActivities           SymbolSet = 40 // 40
	SymbolSetSIGINTSpace          SymbolSet = 50 // 50 (E only)
	SymbolSetSIGINTAir            SymbolSet = 51 // 51 (E only)
	SymbolSetSIGINTLand           SymbolSet = 52 // 52 (E only)
	SymbolSetSIGINTSeaSurface     SymbolSet = 53 // 53 (E only)
	SymbolSetSIGINTSubsurface     SymbolSet = 54 // 54 (E only)
	SymbolSetCyberspace           SymbolSet = 60 // 60 (E only)
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
