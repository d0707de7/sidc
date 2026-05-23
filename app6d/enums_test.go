package app6d

import "testing"

func TestVersion_Family(t *testing.T) {
	tests := []struct {
		name    string
		version Version
		isD     bool
		isE     bool
	}{
		{name: "VersionD10 is D family", version: VersionD10, isD: true, isE: false},
		{name: "VersionD11 is D family", version: VersionD11, isD: true, isE: false},
		{name: "VersionD12 is D family", version: VersionD12, isD: true, isE: false},
		{name: "VersionE13 is E family", version: VersionE13, isD: false, isE: true},
		{name: "VersionE14 is E family", version: VersionE14, isD: false, isE: true},
		{name: "VersionUnspecified is neither family", version: VersionUnspecified, isD: false, isE: false},
		{name: "unknown version 99 is neither family", version: Version(99), isD: false, isE: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.version.IsD(); got != tt.isD {
				t.Errorf("IsD() = %v, expected %v", got, tt.isD)
			}
			if got := tt.version.IsE(); got != tt.isE {
				t.Errorf("IsE() = %v, expected %v", got, tt.isE)
			}
		})
	}
}

func TestHQTFD_Bits(t *testing.T) {
	tests := []struct {
		name         string
		value        HQTFD
		feintDummy   bool
		headquarters bool
		taskForce    bool
		display      string
	}{
		{name: "none has no bits", value: HQTFDNone, display: "None"},
		{name: "feint dummy bit alone", value: HQTFDFeintDummy, feintDummy: true, display: "FeintDummy"},
		{name: "headquarters bit alone", value: HQTFDHeadquarters, headquarters: true, display: "Headquarters"},
		{name: "feint dummy combined with headquarters", value: HQTFDFeintDummyHeadquarters, feintDummy: true, headquarters: true, display: "FeintDummy+Headquarters"},
		{name: "task force bit alone", value: HQTFDTaskForce, taskForce: true, display: "TaskForce"},
		{name: "feint dummy combined with task force", value: HQTFDFeintDummyTaskForce, feintDummy: true, taskForce: true, display: "FeintDummy+TaskForce"},
		{name: "headquarters combined with task force", value: HQTFDHeadquartersTaskForce, headquarters: true, taskForce: true, display: "Headquarters+TaskForce"},
		{name: "all three bits set", value: HQTFDFeintDummyHeadquartersTaskForce, feintDummy: true, headquarters: true, taskForce: true, display: "FeintDummy+Headquarters+TaskForce"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.value.FeintDummy(); got != tt.feintDummy {
				t.Errorf("FeintDummy() = %v, expected %v", got, tt.feintDummy)
			}
			if got := tt.value.Headquarters(); got != tt.headquarters {
				t.Errorf("Headquarters() = %v, expected %v", got, tt.headquarters)
			}
			if got := tt.value.TaskForce(); got != tt.taskForce {
				t.Errorf("TaskForce() = %v, expected %v", got, tt.taskForce)
			}
			if got := tt.value.String(); got != tt.display {
				t.Errorf("String() = %q, expected %q", got, tt.display)
			}
		})
	}
}

func TestAmplifier_Category(t *testing.T) {
	tests := []struct {
		name      string
		amplifier Amplifier
		category  AmplifierCategory
	}{
		{name: "none has no category", amplifier: AmplifierNone, category: AmplifierCategoryNone},
		{name: "team or crew is an echelon", amplifier: AmplifierTeamCrew, category: AmplifierCategoryEchelon},
		{name: "command is an echelon", amplifier: AmplifierCommand, category: AmplifierCategoryEchelon},
		{name: "tracked is a mobility", amplifier: AmplifierTracked, category: AmplifierCategoryMobility},
		{name: "long towed array is a mobility", amplifier: AmplifierLongTowedArray, category: AmplifierCategoryMobility},
		{name: "leader individual is leadership", amplifier: AmplifierLeaderIndividual, category: AmplifierCategoryLeadership},
		{name: "deputy individual is leadership", amplifier: AmplifierDeputyIndividual, category: AmplifierCategoryLeadership},
		{name: "value above leadership range has no category", amplifier: Amplifier(99), category: AmplifierCategoryNone},
		{name: "value below smallest echelon has no category", amplifier: Amplifier(10), category: AmplifierCategoryNone},
		{name: "value between mobility and leadership has no category", amplifier: Amplifier(65), category: AmplifierCategoryNone},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.amplifier.Category(); got != tt.category {
				t.Errorf("Category() = %v, expected %v", got, tt.category)
			}
		})
	}
}

func TestSymbolSet_IsEOnly(t *testing.T) {
	tests := []struct {
		name      string
		symbolSet SymbolSet
		eOnly     bool
	}{
		{name: "air is valid in D and E", symbolSet: SymbolSetAir, eOnly: false},
		{name: "land unit is valid in D and E", symbolSet: SymbolSetLandUnit, eOnly: false},
		{name: "cyberspace is E only", symbolSet: SymbolSetCyberspace, eOnly: true},
		{name: "dismounted individual is E only", symbolSet: SymbolSetDismountedIndividual, eOnly: true},
		{name: "air missile is E only", symbolSet: SymbolSetAirMissile, eOnly: true},
		{name: "SIGINT land is E only", symbolSet: SymbolSetSIGINTLand, eOnly: true},
		{name: "mine warfare is valid in D and E", symbolSet: SymbolSetMineWarfare, eOnly: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.symbolSet.IsEOnly(); got != tt.eOnly {
				t.Errorf("IsEOnly() = %v, expected %v", got, tt.eOnly)
			}
		})
	}
}
