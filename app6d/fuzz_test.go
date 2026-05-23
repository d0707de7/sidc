package app6d

import (
	"strings"
	"testing"
	"unicode/utf8"
)

// FuzzParse asserts three properties for arbitrary input:
//
//  1. Parse never panics.
//  2. Whenever Parse accepts the input, the parsed SIDC re-renders (via the
//     internal render function reached through Value when Validate passes,
//     or via Parse-then-Parse symmetry) to the same byte sequence.
//  3. When Validate passes too, Value yields exactly the input Parse received.
func FuzzParse(f *testing.F) {
	// Seed with a mix of valid and invalid inputs to give the fuzzer
	// realistic starting points.
	f.Add("")
	f.Add(strings.Repeat("0", SIDCLength))
	f.Add(landUnitInfantryPlatoon.render())
	f.Add(airFighterWithModifiers.render())
	f.Add("1003100014121100000a") // structurally invalid (letter)
	f.Add("99031000141211000000") // semantically invalid (bad version)
	f.Add("\xff\xfe\xfd\xfc")     // non-UTF-8 short input
	f.Add("\U0001f600 - emoji")   // valid UTF-8, wrong length
	// Whitespace characters at the correct length must still be rejected.
	f.Add("\t" + strings.Repeat("0", SIDCLength-1))
	f.Add("\n" + strings.Repeat("0", SIDCLength-1))
	f.Add("\r" + strings.Repeat("0", SIDCLength-1))
	f.Add(strings.Repeat("0", SIDCLength-1) + "\n")
	f.Add("10031000\r\n1211000000") // embedded CRLF inside an otherwise-valid SIDC

	f.Fuzz(func(t *testing.T, input string) {
		// Property 1: Parse must never panic.
		sidc, err := Parse(input)
		if err != nil {
			// Errors are fine; the only requirement is that Parse returns.
			return
		}

		// Anything Parse accepted is 20 ASCII digits (it checked that), so the
		// internal render must produce the same string we received.
		encoded := sidc.render()
		if encoded != input {
			t.Fatalf("Parse accepted %q but re-rendering yielded %q", input, encoded)
		}
		if !utf8.ValidString(encoded) {
			t.Fatalf("render produced invalid UTF-8 from accepted input %q: %q", input, encoded)
		}

		// Property 2: round-trip through Parse again must yield an equal SIDC.
		reparsed, err := Parse(encoded)
		if err != nil {
			t.Fatalf("re-parsing accepted input failed: %q, error: %v", encoded, err)
		}
		if reparsed != sidc {
			t.Fatalf("Parse round trip mismatch: original %q, parsed %+v, reparsed %+v",
				input, sidc, reparsed)
		}

		// Property 3: when the parsed SIDC is also semantically valid, Value
		// must yield exactly the original input.
		if err := sidc.Validate(); err == nil {
			v, err := sidc.Value()
			if err != nil {
				t.Fatalf("Validate accepted %+v but Value rejected it: %v", sidc, err)
			}
			if v != input {
				t.Fatalf("Value mismatch: input %q, Value() %q", input, v)
			}
		}
	})
}
