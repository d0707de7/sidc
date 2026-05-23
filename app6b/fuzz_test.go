package app6b

import (
	"strings"
	"testing"
	"unicode/utf8"
)

// FuzzParse asserts Parse never panics regardless of input, and that any
// 15-character input it accepts round-trips cleanly through String.
func FuzzParse(f *testing.F) {
	f.Add("")
	f.Add(strings.Repeat("-", SIDCLength))
	f.Add("SFAPMFF--------")
	f.Add("SHGPUCD-----US-")
	f.Add("\xff\xfe\xfd\xfc")                         // non-UTF-8 short input
	f.Add("\U0001f600\U0001f600\U0001f600\U0001f600") // emoji, wrong byte count
	f.Add(strings.Repeat("\x00", SIDCLength))         // 15 null bytes
	// Whitespace characters at the correct length must still be rejected.
	f.Add("\t" + strings.Repeat("-", SIDCLength-1))
	f.Add("\n" + strings.Repeat("-", SIDCLength-1))
	f.Add("\r" + strings.Repeat("-", SIDCLength-1))
	f.Add("SFAPMFF\t-------")                   // tab inside FunctionID
	f.Add("SFAPMFF--------\n")                  // 16 bytes with trailing newline
	f.Add("SFAPMFF--------\r\n")                // 17 bytes with CRLF
	f.Add(strings.Repeat("\t", SIDCLength))     // 15 tabs
	f.Add("SFAPMFF\xc2\xa0------")              // U+00A0 NBSP encoded as 2 bytes (also wrong length)

	f.Fuzz(func(t *testing.T, input string) {
		sidc, err := Parse(input)
		if err != nil {
			return
		}

		// Parse only accepts inputs of length SIDCLength with printable
		// ASCII, so the internal render must return the same string.
		encoded := sidc.render()
		if encoded != input {
			t.Fatalf("Parse accepted %q but render produced %q", input, encoded)
		}
		if !utf8.ValidString(encoded) {
			t.Fatalf("render produced invalid UTF-8 from accepted input %q: %q", input, encoded)
		}

		reparsed, err := Parse(encoded)
		if err != nil {
			t.Fatalf("re-parsing accepted input failed: %q, error: %v", encoded, err)
		}
		if reparsed != sidc {
			t.Fatalf("Parse round trip mismatch: input %q, parsed %+v, reparsed %+v",
				input, sidc, reparsed)
		}

		// When the parsed SIDC is also semantically valid, Value must yield
		// exactly the original input.
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
