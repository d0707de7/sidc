// Package app6b parses, builds, and validates 15-character letter-based
// Symbol Identification Codes as defined by STANAG APP-6 B and C, and
// MIL-STD-2525 B and C.
//
// The SIDC layout is:
//
//	pos  0     CodingScheme    (S=warfighting, G=tactical graphics, O=stability, E=EMS)
//	pos  1     Affiliation     (P/U/A/F/N/H/J/K plus exercise variants)
//	pos  2     BattleDimension (P=space, A=air, G=ground, S=sea, U=subsurface, F=SOF, T=tactical, X=other)
//	pos  3     Status          (A=anticipated/planned, P=present, C/D/X/F=condition)
//	pos  4-9   FunctionID      (6-letter hierarchical code, e.g. "MFF---" for fighter)
//	pos 10     SymbolModifier1 (single letter, set-specific)
//	pos 11     SymbolModifier2 (single letter, set-specific)
//	pos 12-13  CountryCode     (two-letter, e.g. "US", "GB", "--")
//	pos 14     OrderOfBattle   (single letter)
//
// Unspecified positions are represented by '-' in the canonical form, or '*'
// in the source tables to indicate "any value". The String method emits '-'
// for any zero (unset) byte. Parse accepts both '-' and the actual values.
package app6b
