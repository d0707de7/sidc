package app6d

import "fmt"

// Entity is the 6-digit field at positions 10-15. Its meaning depends on the
// containing SymbolSet; the same numeric value can name different things in
// different sets.
//
// Generated constants in this package are named with their containing symbol
// set as a prefix, e.g. EntityLandUnit_Infantry, EntityAir_FixedWing.
type Entity uint32

func (e Entity) String() string {
	return fmt.Sprintf("%06d", uint32(e))
}

// Modifier1 is the 2-digit sector-1 modifier at positions 16-17. Meaning depends on SymbolSet.
type Modifier1 uint8

func (m Modifier1) String() string {
	return fmt.Sprintf("%02d", uint8(m))
}

// Modifier2 is the 2-digit sector-2 modifier at positions 18-19. Meaning depends on SymbolSet.
type Modifier2 uint8

func (m Modifier2) String() string {
	return fmt.Sprintf("%02d", uint8(m))
}
