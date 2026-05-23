package app6b

import "fmt"

// CodingScheme is the 1-character coding scheme at position 0.
type CodingScheme byte

const (
	CodingSchemeWarfighting      CodingScheme = 'S'
	CodingSchemeTacticalGraphics CodingScheme = 'G'
	CodingSchemeWeather          CodingScheme = 'W'
	CodingSchemeIntelligence     CodingScheme = 'I'
	CodingSchemeMOOTW            CodingScheme = 'O' // military operations other than war / stability
	CodingSchemeEmergencyMgmt    CodingScheme = 'E'
)

func (c CodingScheme) String() string {
	switch c {
	case CodingSchemeWarfighting:
		return "Warfighting"
	case CodingSchemeTacticalGraphics:
		return "Tactical graphics"
	case CodingSchemeWeather:
		return "Weather"
	case CodingSchemeIntelligence:
		return "Intelligence"
	case CodingSchemeMOOTW:
		return "Stability operations"
	case CodingSchemeEmergencyMgmt:
		return "Emergency management"
	case 0, '-':
		return "Unspecified"
	default:
		return fmt.Sprintf("CodingScheme(%q)", byte(c))
	}
}

// Affiliation is the 1-character affiliation at position 1.
// Values map roughly into four families per the JS metadata logic:
//
//	Hostile family:  H, S, J, K
//	Friend family:   F, A, D, M
//	Neutral family:  N, L
//	Unknown family:  P, U, G, W, O
type Affiliation byte

const (
	AffiliationPending               Affiliation = 'P'
	AffiliationUnknown               Affiliation = 'U'
	AffiliationAssumedFriend         Affiliation = 'A'
	AffiliationFriend                Affiliation = 'F'
	AffiliationNeutral               Affiliation = 'N'
	AffiliationSuspect               Affiliation = 'S'
	AffiliationHostile               Affiliation = 'H'
	AffiliationJoker                 Affiliation = 'J'
	AffiliationFaker                 Affiliation = 'K'
	AffiliationExercisePending       Affiliation = 'G'
	AffiliationExerciseUnknown       Affiliation = 'W'
	AffiliationExerciseFriend        Affiliation = 'D'
	AffiliationExerciseNeutral       Affiliation = 'L'
	AffiliationExerciseAssumedFriend Affiliation = 'M'
	AffiliationNoneSpecified         Affiliation = 'O'
)

// Family describes which of the four broad colour groups this affiliation maps to.
type AffiliationFamily uint8

const (
	AffiliationFamilyUnknown AffiliationFamily = iota
	AffiliationFamilyFriend
	AffiliationFamilyNeutral
	AffiliationFamilyHostile
)

func (a Affiliation) Family() AffiliationFamily {
	switch a {
	case AffiliationHostile, AffiliationSuspect, AffiliationJoker, AffiliationFaker:
		return AffiliationFamilyHostile
	case AffiliationFriend, AffiliationAssumedFriend, AffiliationExerciseFriend, AffiliationExerciseAssumedFriend:
		return AffiliationFamilyFriend
	case AffiliationNeutral, AffiliationExerciseNeutral:
		return AffiliationFamilyNeutral
	case AffiliationPending, AffiliationUnknown, AffiliationExercisePending, AffiliationExerciseUnknown, AffiliationNoneSpecified:
		return AffiliationFamilyUnknown
	}
	return AffiliationFamilyUnknown
}

// IsExercise reports whether this affiliation is one of the exercise variants.
func (a Affiliation) IsExercise() bool {
	switch a {
	case AffiliationExercisePending, AffiliationExerciseUnknown, AffiliationExerciseFriend, AffiliationExerciseNeutral, AffiliationExerciseAssumedFriend, AffiliationJoker, AffiliationFaker:
		return true
	}
	return false
}

func (a Affiliation) String() string {
	switch a {
	case AffiliationPending:
		return "Pending"
	case AffiliationUnknown:
		return "Unknown"
	case AffiliationAssumedFriend:
		return "Assumed Friend"
	case AffiliationFriend:
		return "Friend"
	case AffiliationNeutral:
		return "Neutral"
	case AffiliationSuspect:
		return "Suspect"
	case AffiliationHostile:
		return "Hostile"
	case AffiliationJoker:
		return "Joker (exercise hostile)"
	case AffiliationFaker:
		return "Faker (exercise hostile)"
	case AffiliationExercisePending:
		return "Exercise pending"
	case AffiliationExerciseUnknown:
		return "Exercise unknown"
	case AffiliationExerciseFriend:
		return "Exercise friend"
	case AffiliationExerciseNeutral:
		return "Exercise neutral"
	case AffiliationExerciseAssumedFriend:
		return "Exercise assumed friend"
	case AffiliationNoneSpecified:
		return "None specified"
	case 0, '-':
		return "Unspecified"
	default:
		return fmt.Sprintf("Affiliation(%q)", byte(a))
	}
}

// BattleDimension is the 1-character battle dimension at position 2.
type BattleDimension byte

const (
	BattleDimensionSpace           BattleDimension = 'P'
	BattleDimensionAir             BattleDimension = 'A'
	BattleDimensionGround          BattleDimension = 'G'
	BattleDimensionGroundEquipment BattleDimension = 'Z'
	BattleDimensionSOF             BattleDimension = 'F'
	BattleDimensionTactical        BattleDimension = 'X'
	BattleDimensionSeaSurface      BattleDimension = 'S'
	BattleDimensionSubsurface      BattleDimension = 'U'
	BattleDimensionTacticalGraphic BattleDimension = 'T'
	BattleDimensionOther           BattleDimension = 'O'
	BattleDimensionVehicleEvent    BattleDimension = 'V'
	BattleDimensionRiot            BattleDimension = 'R'
)

func (b BattleDimension) String() string {
	switch b {
	case BattleDimensionSpace:
		return "Space"
	case BattleDimensionAir:
		return "Air"
	case BattleDimensionGround:
		return "Ground"
	case BattleDimensionGroundEquipment:
		return "Ground equipment"
	case BattleDimensionSOF:
		return "Special operations forces"
	case BattleDimensionTactical:
		return "Tactical"
	case BattleDimensionSeaSurface:
		return "Sea surface"
	case BattleDimensionSubsurface:
		return "Subsurface"
	case BattleDimensionTacticalGraphic:
		return "Tactical graphic"
	case BattleDimensionOther:
		return "Other"
	case BattleDimensionVehicleEvent:
		return "Vehicle/event"
	case BattleDimensionRiot:
		return "Civil disturbance"
	case 0, '-':
		return "Unspecified"
	default:
		return fmt.Sprintf("BattleDimension(%q)", byte(b))
	}
}

// Status is the 1-character status at position 3. Letters cover both planning
// status (A=anticipated, P=present) and condition (C=fully capable, D=damaged,
// X=destroyed, F=full to capacity).
type Status byte

const (
	StatusPresent        Status = 'P'
	StatusAnticipated    Status = 'A'
	StatusFullyCapable   Status = 'C'
	StatusDamaged        Status = 'D'
	StatusDestroyed      Status = 'X'
	StatusFullToCapacity Status = 'F'
)

func (s Status) String() string {
	switch s {
	case StatusPresent:
		return "Present"
	case StatusAnticipated:
		return "Anticipated/Planned"
	case StatusFullyCapable:
		return "Fully Capable"
	case StatusDamaged:
		return "Damaged"
	case StatusDestroyed:
		return "Destroyed"
	case StatusFullToCapacity:
		return "Full to Capacity"
	case 0, '-':
		return "Unspecified"
	default:
		return fmt.Sprintf("Status(%q)", byte(s))
	}
}
