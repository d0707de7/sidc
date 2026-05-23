// Package app6d parses, builds, and validates 20-digit Symbol Identification
// Codes as defined by STANAG APP-6 D and E, and MIL-STD-2525 D and E.
//
// The SIDC is a fixed-width 20-character numeric string. The layout is:
//
//	pos  0-1   Version    (10-12 = APP-6 D, 13-14 = APP-6 E)
//	pos  2     Context    (Reality / Exercise / Simulation)
//	pos  3     Affiliation
//	pos  4-5   SymbolSet  (01 Air, 10 Land unit, 60 Cyberspace, etc.)
//	pos  6     Status     (Present / Planned / Damaged / etc.)
//	pos  7     HQTFD      (bitfield: feint/dummy, HQ, task force)
//	pos  8-9   Amplifier  (echelon, mobility, or leadership)
//	pos 10-15  Entity     (6-digit code, meaning depends on SymbolSet)
//	pos 16-17  Modifier1  (sector 1 modifier, meaning depends on SymbolSet)
//	pos 18-19  Modifier2  (sector 2 modifier, meaning depends on SymbolSet)
//
// The zero value of SIDC is the syntactically valid string
// "10000000000000000000" (APP-6 D, reality, pending affiliation, unknown
// symbol set, present, no amplifier, no entity, no modifiers).
package app6d
