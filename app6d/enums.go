package app6d

import "fmt"

// Version is the 2-digit version field at positions 0-1.
type Version uint8

const (
	VersionUnspecified Version = 0
	VersionD10         Version = 10
	VersionD11         Version = 11
	VersionD12         Version = 12
	VersionE13         Version = 13
	VersionE14         Version = 14
)

// IsE reports whether the version belongs to the APP-6 E family.
func (v Version) IsE() bool { return v == VersionE13 || v == VersionE14 }

// IsD reports whether the version belongs to the APP-6 D family.
func (v Version) IsD() bool { return v == VersionD10 || v == VersionD11 || v == VersionD12 }

func (v Version) String() string {
	switch v {
	case VersionUnspecified:
		return "Unspecified"
	case VersionD10, VersionD11, VersionD12:
		return fmt.Sprintf("APP-6 D (v%02d)", uint8(v))
	case VersionE13, VersionE14:
		return fmt.Sprintf("APP-6 E (v%02d)", uint8(v))
	default:
		return fmt.Sprintf("Version(%02d)", uint8(v))
	}
}

// Context is the 1-digit context field at position 2 (also called Standard Identity 1).
type Context uint8

const (
	ContextReality    Context = 0
	ContextExercise   Context = 1
	ContextSimulation Context = 2
)

func (c Context) String() string {
	switch c {
	case ContextReality:
		return "Reality"
	case ContextExercise:
		return "Exercise"
	case ContextSimulation:
		return "Simulation"
	default:
		return fmt.Sprintf("Context(%d)", uint8(c))
	}
}

// Affiliation is the 1-digit affiliation field at position 3 (also called Standard Identity 2).
type Affiliation uint8

const (
	AffiliationPending       Affiliation = 0
	AffiliationUnknown       Affiliation = 1
	AffiliationAssumedFriend Affiliation = 2
	AffiliationFriend        Affiliation = 3
	AffiliationNeutral       Affiliation = 4
	AffiliationSuspect       Affiliation = 5
	AffiliationHostile       Affiliation = 6
)

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
	default:
		return fmt.Sprintf("Affiliation(%d)", uint8(a))
	}
}

// Status is the 1-digit status field at position 6.
type Status uint8

const (
	StatusPresent        Status = 0
	StatusPlanned        Status = 1
	StatusFullyCapable   Status = 2
	StatusDamaged        Status = 3
	StatusDestroyed      Status = 4
	StatusFullToCapacity Status = 5
)

func (s Status) String() string {
	switch s {
	case StatusPresent:
		return "Present"
	case StatusPlanned:
		return "Planned/Anticipated"
	case StatusFullyCapable:
		return "Fully Capable"
	case StatusDamaged:
		return "Damaged"
	case StatusDestroyed:
		return "Destroyed"
	case StatusFullToCapacity:
		return "Full to Capacity"
	default:
		return fmt.Sprintf("Status(%d)", uint8(s))
	}
}

// HQTFD is the 1-digit headquarters/task force/feint-dummy bitfield at position 7.
// Bit 0 (value 1) = feint or dummy, bit 1 (value 2) = headquarters, bit 2 (value 4) = task force.
type HQTFD uint8

const (
	HQTFDNone                            HQTFD = 0
	HQTFDFeintDummy                      HQTFD = 1
	HQTFDHeadquarters                    HQTFD = 2
	HQTFDFeintDummyHeadquarters          HQTFD = 3
	HQTFDTaskForce                       HQTFD = 4
	HQTFDFeintDummyTaskForce             HQTFD = 5
	HQTFDHeadquartersTaskForce           HQTFD = 6
	HQTFDFeintDummyHeadquartersTaskForce HQTFD = 7
)

// FeintDummy reports whether the feint-or-dummy bit is set.
func (h HQTFD) FeintDummy() bool { return h&1 != 0 }

// Headquarters reports whether the headquarters bit is set.
func (h HQTFD) Headquarters() bool { return h&2 != 0 }

// TaskForce reports whether the task force bit is set.
func (h HQTFD) TaskForce() bool { return h&4 != 0 }

func (h HQTFD) String() string {
	if h == HQTFDNone {
		return "None"
	}
	if h > 7 {
		return fmt.Sprintf("HQTFD(%d)", uint8(h))
	}
	var parts string
	if h.FeintDummy() {
		parts = "FeintDummy"
	}
	if h.Headquarters() {
		if parts != "" {
			parts += "+"
		}
		parts += "Headquarters"
	}
	if h.TaskForce() {
		if parts != "" {
			parts += "+"
		}
		parts += "TaskForce"
	}
	return parts
}

// Amplifier is the 2-digit echelon/mobility/leadership field at positions 8-9.
// Values 11-26 are echelons, 31-37 and 41-42 and 51-52 and 61-62 are mobilities,
// 71-72 are leadership indicators.
type Amplifier uint8

const (
	AmplifierNone Amplifier = 0

	AmplifierTeamCrew            Amplifier = 11
	AmplifierSquad               Amplifier = 12
	AmplifierSection             Amplifier = 13
	AmplifierPlatoonDetachment   Amplifier = 14
	AmplifierCompanyBatteryTroop Amplifier = 15
	AmplifierBattalionSquadron   Amplifier = 16
	AmplifierRegimentGroup       Amplifier = 17
	AmplifierBrigade             Amplifier = 18
	AmplifierDivision            Amplifier = 21
	AmplifierCorpsMEF            Amplifier = 22
	AmplifierArmy                Amplifier = 23
	AmplifierArmyGroupFront      Amplifier = 24
	AmplifierRegionTheater       Amplifier = 25
	AmplifierCommand             Amplifier = 26

	AmplifierWheeledLimitedCrossCountry Amplifier = 31
	AmplifierWheeledCrossCountry        Amplifier = 32
	AmplifierTracked                    Amplifier = 33
	AmplifierWheeledAndTracked          Amplifier = 34
	AmplifierTowed                      Amplifier = 35
	AmplifierRail                       Amplifier = 36
	AmplifierPackAnimals                Amplifier = 37
	AmplifierOverSnow                   Amplifier = 41
	AmplifierSled                       Amplifier = 42
	AmplifierBarge                      Amplifier = 51
	AmplifierAmphibious                 Amplifier = 52
	AmplifierShortTowedArray            Amplifier = 61
	AmplifierLongTowedArray             Amplifier = 62

	AmplifierLeaderIndividual Amplifier = 71
	AmplifierDeputyIndividual Amplifier = 72
)

// Category identifies which family of amplifier a value belongs to.
type AmplifierCategory uint8

const (
	AmplifierCategoryNone AmplifierCategory = iota
	AmplifierCategoryEchelon
	AmplifierCategoryMobility
	AmplifierCategoryLeadership
)

// Category returns the family this amplifier value belongs to.
func (a Amplifier) Category() AmplifierCategory {
	switch {
	case a == AmplifierNone:
		return AmplifierCategoryNone
	case a >= 11 && a <= 26:
		return AmplifierCategoryEchelon
	case a >= 31 && a <= 62:
		return AmplifierCategoryMobility
	case a >= 71 && a <= 72:
		return AmplifierCategoryLeadership
	default:
		return AmplifierCategoryNone
	}
}

func (a Amplifier) String() string {
	switch a {
	case AmplifierNone:
		return "None"
	case AmplifierTeamCrew:
		return "Team/Crew"
	case AmplifierSquad:
		return "Squad"
	case AmplifierSection:
		return "Section"
	case AmplifierPlatoonDetachment:
		return "Platoon/Detachment"
	case AmplifierCompanyBatteryTroop:
		return "Company/Battery/Troop"
	case AmplifierBattalionSquadron:
		return "Battalion/Squadron"
	case AmplifierRegimentGroup:
		return "Regiment/Group"
	case AmplifierBrigade:
		return "Brigade"
	case AmplifierDivision:
		return "Division"
	case AmplifierCorpsMEF:
		return "Corps/MEF"
	case AmplifierArmy:
		return "Army"
	case AmplifierArmyGroupFront:
		return "Army Group/Front"
	case AmplifierRegionTheater:
		return "Region/Theater"
	case AmplifierCommand:
		return "Command"
	case AmplifierWheeledLimitedCrossCountry:
		return "Wheeled, limited cross-country"
	case AmplifierWheeledCrossCountry:
		return "Wheeled, cross-country"
	case AmplifierTracked:
		return "Tracked"
	case AmplifierWheeledAndTracked:
		return "Wheeled and tracked combination"
	case AmplifierTowed:
		return "Towed"
	case AmplifierRail:
		return "Rail"
	case AmplifierPackAnimals:
		return "Pack animals"
	case AmplifierOverSnow:
		return "Over snow (prime mover)"
	case AmplifierSled:
		return "Sled"
	case AmplifierBarge:
		return "Barge"
	case AmplifierAmphibious:
		return "Amphibious"
	case AmplifierShortTowedArray:
		return "Short towed array"
	case AmplifierLongTowedArray:
		return "Long towed array"
	case AmplifierLeaderIndividual:
		return "Leader individual"
	case AmplifierDeputyIndividual:
		return "Deputy individual"
	default:
		return fmt.Sprintf("Amplifier(%02d)", uint8(a))
	}
}
