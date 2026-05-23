# sidc

Parse, build, and validate Symbol Identification Codes (SIDC) in Go.

Supports the standards in use today:

- **APP-6 B and C** (and MIL-STD-2525 B/C) — 15-character letter-based encoding, in [`app6b/`](./app6b/).
- **APP-6 D and E** (and MIL-STD-2525 D/E) — 20-character number-based encoding, in [`app6d/`](./app6d/).

This module is a Go port of the SIDC-handling parts of Måns Beckman's
JavaScript libraries, [`milsymbol`](https://github.com/spatialillusions/milsymbol),
[`stanag-app6`](https://github.com/spatialillusions/stanag-app6), and
[`milstandard-e`](https://github.com/spatialillusions/milstandard-e). The
field layouts, symbol-set mappings, version-detection rules, and lookup
tables all come from those projects; this module restates them as typed
Go code with table-driven tests and fuzzing.

This module only handles the SIDC string. Rendering symbols to SVG, canvas,
or icons is out of scope; pair it with [`milsymbol`](https://github.com/spatialillusions/milsymbol)
on the JavaScript side if you need pictures.

## Install

```
go get github.com/d0707de7/sidc
```

## What's a SIDC?

A SIDC is a fixed-width string that identifies a military symbol: its
affiliation (friend, hostile, neutral), what it represents (a fighter
aircraft, an infantry company, a control measure), its status, and so on.
Different versions of the standard use different encodings; this module
covers all the encodings currently in field use.

For the 20-character APP-6 D/E encoding, the layout is:

| Pos   | Width | Field       |
|-------|-------|-------------|
| 0-1   | 2     | Version     |
| 2     | 1     | Context     |
| 3     | 1     | Affiliation |
| 4-5   | 2     | Symbol set  |
| 6     | 1     | Status      |
| 7     | 1     | HQ/TF/Dummy |
| 8-9   | 2     | Amplifier   |
| 10-15 | 6     | Entity      |
| 16-17 | 2     | Modifier 1  |
| 18-19 | 2     | Modifier 2  |

Each field has a typed Go enum, and the entity/modifier values for each
symbol set are generated from the official tables — so anything that
compiles is a known code, and the IDE can complete `EntityLandUnit_` to
show every land-unit entity.

## Quick start

### Build a SIDC from typed constants

`SIDC` does not implement `fmt.Stringer`. To get the wire-format string,
call `Value`, which returns `(string, error)` — it validates first and
refuses to render anything that isn't a real symbol. This means you
literally cannot ship an invalid SIDC.

```go
package main

import (
	"fmt"
	"log"

	"github.com/d0707de7/sidc/app6d"
)

func main() {
	s := app6d.SIDC{
		Version:     app6d.VersionE13,
		Context:     app6d.ContextReality,
		Affiliation: app6d.AffiliationFriend,
		SymbolSet:   app6d.SymbolSetAir,
		Status:      app6d.StatusPresent,
		Entity:      app6d.EntityAir_MilitaryFixedWingFighter,
		Modifier1:   app6d.Modifier1Air_Fighter,
		Modifier2:   app6d.Modifier2Air_BoomOnly,
	}

	v, err := s.Value()
	if err != nil {
		log.Fatalf("constructed SIDC is not valid: %v", err)
	}
	fmt.Println(v)
	// 13030100001101040404

	fmt.Println(s.Entity.Name(s.SymbolSet))
	// Military / Fixed Wing / Fighter
}
```

If you want to inspect what's wrong without rendering, call `Validate`
directly; `Value` is just a thin wrapper around it.

### Parse a SIDC

```go
s, err := app6d.Parse("13030100001101040404")
if err != nil {
    return err
}

fmt.Println(s.Affiliation)              // Friend
fmt.Println(s.SymbolSet)                // Air
fmt.Println(s.Entity.Name(s.SymbolSet)) // Military / Fixed Wing / Fighter
```

`Parse` only checks the **structure** of the input: that it is the
expected length and that every byte is in the expected character range.
It does not check that the resulting field values are meaningful. To
verify the SIDC names a real symbol, call `Validate` on the result (or
call `Value` to do both at once when you want the canonical string back):

```go
s, err := app6d.Parse(input)
if err != nil {
    return fmt.Errorf("malformed SIDC: %w", err)
}
if err := s.Validate(); err != nil {
    return fmt.Errorf("SIDC is well-formed but not meaningful: %w", err)
}
```

The split is deliberate. Parse always gives you a `SIDC` you can inspect
and repair when it is structurally valid; Validate is the gate for using
the value as a real symbol.

### JSON, XML, and other text encodings

`SIDC` implements `encoding.TextMarshaler` and `encoding.TextUnmarshaler`,
which means JSON, XML, and any other text-based encoder serialise it as a
string (not a struct dump). Encoding fails if the SIDC is not valid, so
marshalling can't ship a broken value either.

```go
type Track struct {
    Name string      `json:"name"`
    SIDC app6d.SIDC  `json:"sidc"`
}

raw, err := json.Marshal(Track{
    Name: "alpha-1",
    SIDC: s,
})
// raw is `{"name":"alpha-1","sidc":"13030100001101040404"}`
```

### Detect which standard a SIDC belongs to

When you receive a SIDC from an unknown source, use the top-level `Detect`
to dispatch to the right package:

```go
import (
    "github.com/d0707de7/sidc"
    "github.com/d0707de7/sidc/app6b"
    "github.com/d0707de7/sidc/app6d"
)

func handle(input string) error {
    version, ok := sidc.Detect(input)
    if !ok {
        return fmt.Errorf("input %q is not a recognised SIDC layout", input)
    }
    switch version {
    case sidc.VersionAPP6B:
        s, err := app6b.Parse(input)
        // ...
    case sidc.VersionAPP6D, sidc.VersionAPP6E:
        s, err := app6d.Parse(input)
        // ...
    }
    return nil
}
```

`Detect` is a cheap structural check (length + leading characters); it
does not parse the body of the SIDC.

### APP-6 B and C (letter-based)

```go
s := app6b.SIDC{
    CodingScheme:    app6b.CodingSchemeWarfighting,
    Affiliation:     app6b.AffiliationFriend,
    BattleDimension: app6b.BattleDimensionAir,
    Status:          app6b.StatusPresent,
    FunctionID:      app6b.FunctionID{'M', 'F', 'F', '-', '-', '-'},
    CountryCode:     [2]byte{'U', 'S'},
}
v, err := s.Value()
// v == "SFAPMFF-----US-"
```

Same `Value`/`Validate`/`MarshalText` shape as `app6d`. The letter-based
encoding pre-dates the symbol-set table format, so the package has fewer
auto-generated constants — you supply the function ID directly as a 6-byte
array, and the validator only checks that every byte is printable ASCII.

## How the entity / modifier constants are organised

Each symbol set has its own generated file under `app6d/`, with constants
prefixed by the set name:

```
EntityAir_MilitaryFixedWingFighter
EntityLandUnit_MovementAndManeuverInfantry
EntityCyberspace_CyberspaceUnitCyberspaceUnitNonSpecified

Modifier1Air_Fighter
Modifier2Air_BoomOnly
```

The same numeric entity code can mean different things in different symbol
sets — `110000` is "Military" in Air but "Command and Control" in Land
unit — so entities are not portable between sets. The `Entity.Name(set)`
method takes the symbol set as context.

## Regenerating from the source tables

The entity and modifier constants are generated from the official TSV
tables vendored under [`tables/`](./tables/). To regenerate after updating
the tables:

```
go generate ./...
```

The generator lives in [`internal/tsvgen`](./internal/tsvgen/).

## Parse, Validate, Value

The three operations split along three concerns:

| Function | What it does | When to call |
|----------|--------------|--------------|
| `Parse` | Checks length and character range only. Returns a `SIDC` for any structurally well-formed input, even if the field values don't name a real symbol. | When you have a SIDC string and want a typed value. |
| `Validate` | Checks every field carries a meaningful value: enums in range, known symbol set, E-only sets only in E versions, amplifier in a defined range, entity and any non-zero modifiers defined for the symbol set. Returns an error naming the first invalid field. | When you want to know whether a SIDC names a real symbol, without rendering it. |
| `Value` | Calls `Validate`, then returns the canonical 20-digit (or 15-character) string. Refuses to render anything Validate would reject. | When you want the wire-format string to hand to another system. |

`Parse` and `Validate` are independent so a caller can inspect or repair a
structurally valid but semantically broken SIDC. `Value` is the only public
way to get the string, which means the API surface itself prevents shipping
a broken SIDC to a downstream system.

`SIDC` does not implement `fmt.Stringer`. `fmt.Println(s)` prints the
struct dump rather than the SIDC string — that's deliberate, since
rendering can fail and `Stringer.String()` cannot signal an error.

## Stability

The module follows semantic versioning. Generated entity and modifier
constants are part of the public API and will not change name across
non-major releases; if upstream renames a code, the old name stays as an
alias until the next major version.

## Upstream credits

The design of this module is a direct port of Måns Beckman's
MIT-licensed JavaScript libraries; without those, this project would not
exist. All three are works to look at if you want to understand the
standards in detail:

- [`milsymbol`](https://github.com/spatialillusions/milsymbol) — the
  reference JS library for rendering MIL-STD-2525 / STANAG APP-6 symbols.
  This module reproduces (in Go) the SIDC parsing logic from
  [`src/numbersidc/metadata.js`](https://github.com/spatialillusions/milsymbol/blob/master/src/numbersidc/metadata.js)
  and [`src/lettersidc/metadata.js`](https://github.com/spatialillusions/milsymbol/blob/master/src/lettersidc/metadata.js):
  the field positions, the version → standard mapping, the symbol-set →
  dimension/affiliation rules, and the HQ/task-force/feint-dummy bitfield
  semantics.
- [`stanag-app6`](https://github.com/spatialillusions/stanag-app6) — the
  TSV tables for APP-6 B and APP-6 D. The files under
  [`tables/app6b/`](./tables/app6b/) and [`tables/app6d/`](./tables/app6d/)
  are copied verbatim from this project's `tsv-tables/` directory.
- [`milstandard-e`](https://github.com/spatialillusions/milstandard-e) —
  the TSV tables for APP-6 E. The files under
  [`tables/app6e/`](./tables/app6e/) are copied verbatim from this
  project's `tsv-tables/` directory.

Both upstream projects, and this module, are MIT-licensed. If you spot a
bug in the field semantics or the table data, the fix probably belongs
upstream first.

## Licence

MIT.
