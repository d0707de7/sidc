// Command tsvgen reads vendored TSV tables under ../../tables/ and emits Go
// source files in the app6d package containing typed Entity and Modifier
// constants per symbol set, plus a composite name-lookup map.
//
// Run via go:generate from the module root.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"go/format"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

func main() {
	tablesDir := flag.String("tables", "tables", "directory holding the vendored TSV tables")
	outDir := flag.String("out", "app6d", "directory to write generated Go files into")
	flag.Parse()

	log := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo}))

	if err := run(log, *tablesDir, *outDir); err != nil {
		log.Error("generation failed", slog.Any("error", err))
		os.Exit(1)
	}
}

func run(log *slog.Logger, tablesDir, outDir string) (err error) {
	sets, err := collectSymbolSets(tablesDir)
	if err != nil {
		return fmt.Errorf("collecting symbol sets: %w", err)
	}

	log.Info("collected symbol sets", slog.Int("count", len(sets)))

	if err := writeEntityFiles(log, outDir, sets); err != nil {
		return fmt.Errorf("writing entity files: %w", err)
	}
	if err := writeModifierFiles(log, outDir, sets); err != nil {
		return fmt.Errorf("writing modifier files: %w", err)
	}
	if err := writeLookupFile(log, outDir, sets); err != nil {
		return fmt.Errorf("writing lookup file: %w", err)
	}
	return nil
}

// symbolSetTables groups the TSV tables that contribute to one symbol set.
type symbolSetTables struct {
	code        uint8  // 1, 10, 60, etc.
	name        string // human-readable
	goName      string // CamelCase, e.g. "LandUnit"
	mainTable   string // path to <name>.tsv, possibly empty
	modifier1   string // path to <name> sector 1.tsv, possibly empty
	modifier2   string // path to <name> sector 2.tsv, possibly empty
	versionFlag string // "D", "E", or "DE"
}

// symbolSetManifest maps the canonical TSV filename stem to a symbol set code
// and a Go identifier. The same stem can map to multiple codes (SIGINT shares
// one table across 50-54).
type symbolSetManifestEntry struct {
	codes  []uint8
	goName string // singular CamelCase used for the primary code; SIGINT variants append a suffix
}

var symbolSetManifest = map[string]symbolSetManifestEntry{
	"Air":                   {codes: []uint8{1}, goName: "Air"},
	"Air missile":           {codes: []uint8{2}, goName: "AirMissile"},
	"Space":                 {codes: []uint8{5}, goName: "Space"},
	"Space missile":         {codes: []uint8{6}, goName: "SpaceMissile"},
	"Land unit":             {codes: []uint8{10}, goName: "LandUnit"},
	"Land civilian":         {codes: []uint8{11}, goName: "LandCivilian"},
	"Land equipment":        {codes: []uint8{15}, goName: "LandEquipment"},
	"Land installation":     {codes: []uint8{20}, goName: "LandInstallation"},
	"Control Measures":      {codes: []uint8{25}, goName: "ControlMeasure"},
	"Dismounted individual": {codes: []uint8{27}, goName: "DismountedIndividual"},
	"Sea surface":           {codes: []uint8{30}, goName: "SeaSurface"},
	"Sea subsurface":        {codes: []uint8{35}, goName: "SeaSubsurface"},
	"Mine warfare":          {codes: []uint8{36}, goName: "MineWarfare"},
	"Activities":            {codes: []uint8{40}, goName: "Activities"},
	"Signals intelligence":  {codes: []uint8{50, 51, 52, 53, 54}, goName: "SIGINT"},
	"Cyberspace":            {codes: []uint8{60}, goName: "Cyberspace"},
}

// collectSymbolSets walks tablesDir/app6d and tablesDir/app6e, merging
// matching tables. E entries override D entries; the resulting versionFlag
// reflects which version(s) supplied data.
func collectSymbolSets(tablesDir string) ([]symbolSetTables, error) {
	// We expand the SIGINT manifest entry into 5 symbol-set rows but they
	// share a single set of files.
	type accumulator struct {
		main, mod1, mod2 map[string]string // version -> path
	}
	acc := map[string]*accumulator{} // keyed by stem

	for _, version := range []string{"app6d", "app6e"} {
		dir := filepath.Join(tablesDir, version)
		entries, err := os.ReadDir(dir)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("reading %s: %w", dir, err)
		}
		for _, ent := range entries {
			if ent.IsDir() || !strings.HasSuffix(ent.Name(), ".tsv") {
				continue
			}
			stem, kind := classifyTable(ent.Name())
			if stem == "" {
				continue
			}
			// "Common Modifiers" tables hold 3-digit modifier codes that only
			// apply within the extended 30-char SIDC layout. The canonical
			// 20-char SIDC has 2-digit modifiers per symbol set, so skip them.
			if stem == "Common Modifiers" {
				continue
			}
			a, ok := acc[stem]
			if !ok {
				a = &accumulator{
					main: map[string]string{},
					mod1: map[string]string{},
					mod2: map[string]string{},
				}
				acc[stem] = a
			}
			full := filepath.Join(dir, ent.Name())
			switch kind {
			case "main":
				a.main[version] = full
			case "mod1":
				a.mod1[version] = full
			case "mod2":
				a.mod2[version] = full
			}
		}
	}

	var out []symbolSetTables
	for stem, a := range acc {
		entry, ok := symbolSetManifest[stem]
		if !ok {
			return nil, fmt.Errorf("unknown table stem %q (add to symbolSetManifest)", stem)
		}
		for i, code := range entry.codes {
			st := symbolSetTables{
				code:        code,
				name:        stem,
				goName:      entry.goName,
				versionFlag: pickVersionFlag(a.main, a.mod1, a.mod2),
				mainTable:   pickPath(a.main),
				modifier1:   pickPath(a.mod1),
				modifier2:   pickPath(a.mod2),
			}
			if len(entry.codes) > 1 {
				st.goName = entry.goName + sigintSuffix(code)
			}
			_ = i
			out = append(out, st)
		}
	}

	slices.SortFunc(out, func(a, b symbolSetTables) int { return int(a.code) - int(b.code) })
	return out, nil
}

func sigintSuffix(code uint8) string {
	switch code {
	case 50:
		return "Space"
	case 51:
		return "Air"
	case 52:
		return "Land"
	case 53:
		return "SeaSurface"
	case 54:
		return "Subsurface"
	}
	return ""
}

// pickPath prefers the E table over the D table when both exist.
func pickPath(m map[string]string) string {
	if p, ok := m["app6e"]; ok {
		return p
	}
	return m["app6d"]
}

// pickVersionFlag returns "D", "E", or "DE" based on which versions
// contributed any of the three input maps.
func pickVersionFlag(maps ...map[string]string) string {
	hasD, hasE := false, false
	for _, m := range maps {
		if _, ok := m["app6d"]; ok {
			hasD = true
		}
		if _, ok := m["app6e"]; ok {
			hasE = true
		}
	}
	switch {
	case hasD && hasE:
		return "DE"
	case hasE:
		return "E"
	case hasD:
		return "D"
	}
	return ""
}

// classifyTable extracts the stem (e.g. "Land unit") and the kind ("main",
// "mod1", "mod2") from a TSV filename, or returns "" if it should be skipped.
func classifyTable(name string) (stem, kind string) {
	base := strings.TrimSuffix(name, ".tsv")
	switch {
	case strings.HasSuffix(base, " sector 1"):
		return strings.TrimSuffix(base, " sector 1"), "mod1"
	case strings.HasSuffix(base, " sector 2"):
		return strings.TrimSuffix(base, " sector 2"), "mod2"
	default:
		return base, "main"
	}
}

// row represents one parsed TSV row with column lookups by header name.
type row struct {
	cols map[string]string
}

// readTSV parses a TSV file. The first non-empty line is the header.
// Returns header columns in declared order plus rows.
func readTSV(path string) (headers []string, rows []row, err error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	for scanner.Scan() {
		line := strings.TrimRight(scanner.Text(), "\r")
		if line == "" {
			continue
		}
		fields := strings.Split(line, "\t")
		if headers == nil {
			for _, h := range fields {
				headers = append(headers, strings.TrimSpace(h))
			}
			continue
		}
		r := row{cols: map[string]string{}}
		for i, h := range headers {
			if i < len(fields) {
				r.cols[h] = strings.TrimSpace(fields[i])
			}
		}
		rows = append(rows, r)
	}
	if err := scanner.Err(); err != nil {
		return nil, nil, fmt.Errorf("scanning %s: %w", path, err)
	}
	return headers, rows, nil
}

// entityRow is a parsed row from a main-icon TSV with carried-forward hierarchy.
type entityRow struct {
	code    uint32
	name    string // "Entity / Entity Type / Entity Subtype" joined
	remarks string
}

func readEntities(path string) ([]entityRow, error) {
	_, rows, err := readTSV(path)
	if err != nil {
		return nil, err
	}
	var out []entityRow
	var lastEntity, lastEntityType string
	for _, r := range rows {
		ent := r.cols["Entity"]
		et := r.cols["Entity Type"]
		es := r.cols["Entity Subtype"]
		code := r.cols["Code"]
		if code == "" {
			continue
		}
		if ent != "" {
			lastEntity = ent
			lastEntityType = ""
		}
		if et != "" {
			lastEntityType = et
		}

		var nameParts []string
		switch {
		case es != "":
			if lastEntity != "" {
				nameParts = append(nameParts, lastEntity)
			}
			if lastEntityType != "" {
				nameParts = append(nameParts, lastEntityType)
			}
			nameParts = append(nameParts, es)
		case et != "":
			if lastEntity != "" {
				nameParts = append(nameParts, lastEntity)
			}
			nameParts = append(nameParts, et)
		case ent != "":
			nameParts = append(nameParts, ent)
		}

		n, err := strconv.ParseUint(strings.TrimSpace(code), 10, 32)
		if err != nil {
			continue
		}
		out = append(out, entityRow{
			code:    uint32(n),
			name:    strings.Join(nameParts, " / "),
			remarks: r.cols["Remarks"],
		})
	}
	return out, nil
}

// modifierRow is a parsed row from a sector-N TSV.
type modifierRow struct {
	code uint8
	name string
}

func readModifiers(path string) ([]modifierRow, error) {
	headers, rows, err := readTSV(path)
	if err != nil {
		return nil, err
	}
	nameCol := ""
	for _, h := range headers {
		if h == "First Modifier" || h == "Second Modifier" {
			nameCol = h
			break
		}
	}
	if nameCol == "" {
		return nil, fmt.Errorf("no modifier name column in %s (headers: %v)", path, headers)
	}
	var out []modifierRow
	for _, r := range rows {
		raw := strings.TrimSpace(r.cols["Code"])
		if raw == "" {
			continue
		}
		n, err := strconv.ParseUint(raw, 10, 8)
		if err != nil {
			continue
		}
		out = append(out, modifierRow{
			code: uint8(n),
			name: r.cols[nameCol],
		})
	}
	return out, nil
}

// goIdent sanitises a free-text name into a Go identifier suffix.
func goIdent(name string) string {
	var b strings.Builder
	upper := true
	for _, r := range name {
		switch {
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			if upper {
				b.WriteRune(unicode.ToUpper(r))
				upper = false
			} else {
				b.WriteRune(r)
			}
		case r == '/' || r == ' ' || r == '-' || r == '(' || r == ')' || r == ',' || r == '.' || r == '\'' || r == '&':
			upper = true
		case r == '+':
			b.WriteString("Plus")
			upper = false
		default:
			upper = true
		}
	}
	id := b.String()
	if id == "" {
		return "Unknown"
	}
	if unicode.IsDigit(rune(id[0])) {
		id = "N" + id
	}
	return id
}

// writeEntityFiles writes one Go file per symbol set with Entity constants.
func writeEntityFiles(log *slog.Logger, outDir string, sets []symbolSetTables) error {
	for _, set := range sets {
		if set.mainTable == "" {
			continue
		}
		entries, err := readEntities(set.mainTable)
		if err != nil {
			return fmt.Errorf("reading %s: %w", set.mainTable, err)
		}
		if len(entries) == 0 {
			continue
		}

		var src strings.Builder
		fmt.Fprintf(&src, "// Code generated by internal/tsvgen; DO NOT EDIT.\n")
		fmt.Fprintf(&src, "// Source: %s (APP-6 %s)\n\n", filepath.Base(set.mainTable), set.versionFlag)
		fmt.Fprintf(&src, "package app6d\n\n")
		fmt.Fprintf(&src, "// Entity constants for symbol set %02d (%s).\n", set.code, set.name)
		fmt.Fprintf(&src, "const (\n")

		seen := map[string]bool{}
		for i, e := range entries {
			suffix := goIdent(e.name)
			ident := fmt.Sprintf("Entity%s_%s", set.goName, suffix)
			if seen[ident] {
				ident = fmt.Sprintf("%s_%06d", ident, e.code)
			}
			seen[ident] = true
			if i > 0 {
				fmt.Fprintln(&src)
			}
			if e.name != "" {
				fmt.Fprintf(&src, "\t// %s is %s.\n", ident, escapeComment(e.name))
			}
			fmt.Fprintf(&src, "\t%s Entity = %d\n", ident, e.code)
		}
		fmt.Fprintf(&src, ")\n")

		outFile := filepath.Join(outDir, fmt.Sprintf("entities_%s.go", strings.ToLower(setFilenameStem(set))))
		if err := writeFormatted(outFile, src.String()); err != nil {
			return fmt.Errorf("writing %s: %w", outFile, err)
		}
		log.Info("wrote entities", slog.String("file", outFile), slog.Int("count", len(entries)))
	}
	return nil
}

// writeModifierFiles writes one Go file per symbol set per modifier (1 or 2).
func writeModifierFiles(log *slog.Logger, outDir string, sets []symbolSetTables) error {
	for _, set := range sets {
		for sector, path := range map[int]string{1: set.modifier1, 2: set.modifier2} {
			if path == "" {
				continue
			}
			mods, err := readModifiers(path)
			if err != nil {
				return fmt.Errorf("reading %s: %w", path, err)
			}
			if len(mods) == 0 {
				continue
			}
			var src strings.Builder
			fmt.Fprintf(&src, "// Code generated by internal/tsvgen; DO NOT EDIT.\n")
			fmt.Fprintf(&src, "// Source: %s\n\n", filepath.Base(path))
			fmt.Fprintf(&src, "package app6d\n\n")
			fmt.Fprintf(&src, "// Modifier %d constants for symbol set %02d (%s).\n", sector, set.code, set.name)
			fmt.Fprintf(&src, "const (\n")
			seen := map[string]bool{}
			modType := fmt.Sprintf("Modifier%d", sector)
			for i, m := range mods {
				ident := fmt.Sprintf("%s%s_%s", modType, set.goName, goIdent(m.name))
				if seen[ident] {
					ident = fmt.Sprintf("%s_%02d", ident, m.code)
				}
				seen[ident] = true
				if i > 0 {
					fmt.Fprintln(&src)
				}
				if m.name != "" {
					fmt.Fprintf(&src, "\t// %s is %s.\n", ident, escapeComment(m.name))
				}
				fmt.Fprintf(&src, "\t%s %s = %d\n", ident, modType, m.code)
			}
			fmt.Fprintf(&src, ")\n")

			outFile := filepath.Join(outDir, fmt.Sprintf("modifier%d_%s.go", sector, strings.ToLower(setFilenameStem(set))))
			if err := writeFormatted(outFile, src.String()); err != nil {
				return fmt.Errorf("writing %s: %w", outFile, err)
			}
			log.Info("wrote modifiers", slog.String("file", outFile), slog.Int("sector", sector), slog.Int("count", len(mods)))
		}
	}
	return nil
}

// writeLookupFile writes a name-lookup function: given (SymbolSet, Entity)
// return the entity's hierarchical name.
func writeLookupFile(log *slog.Logger, outDir string, sets []symbolSetTables) error {
	type entry struct {
		set  uint8
		code uint32
		name string
	}
	type modEntry struct {
		set  uint8
		code uint8
		name string
	}
	var entityEntries []entry
	var mod1Entries, mod2Entries []modEntry

	for _, set := range sets {
		if set.mainTable != "" {
			es, err := readEntities(set.mainTable)
			if err != nil {
				return err
			}
			for _, e := range es {
				if e.name == "" {
					continue
				}
				entityEntries = append(entityEntries, entry{set.code, e.code, e.name})
			}
		}
		if set.modifier1 != "" {
			ms, err := readModifiers(set.modifier1)
			if err != nil {
				return err
			}
			for _, m := range ms {
				if m.name == "" {
					continue
				}
				mod1Entries = append(mod1Entries, modEntry{set.code, m.code, m.name})
			}
		}
		if set.modifier2 != "" {
			ms, err := readModifiers(set.modifier2)
			if err != nil {
				return err
			}
			for _, m := range ms {
				if m.name == "" {
					continue
				}
				mod2Entries = append(mod2Entries, modEntry{set.code, m.code, m.name})
			}
		}
	}

	sort.Slice(entityEntries, func(i, j int) bool {
		if entityEntries[i].set != entityEntries[j].set {
			return entityEntries[i].set < entityEntries[j].set
		}
		return entityEntries[i].code < entityEntries[j].code
	})
	sort.Slice(mod1Entries, func(i, j int) bool {
		if mod1Entries[i].set != mod1Entries[j].set {
			return mod1Entries[i].set < mod1Entries[j].set
		}
		return mod1Entries[i].code < mod1Entries[j].code
	})
	sort.Slice(mod2Entries, func(i, j int) bool {
		if mod2Entries[i].set != mod2Entries[j].set {
			return mod2Entries[i].set < mod2Entries[j].set
		}
		return mod2Entries[i].code < mod2Entries[j].code
	})

	var src strings.Builder
	fmt.Fprintf(&src, "// Code generated by internal/tsvgen; DO NOT EDIT.\n\n")
	fmt.Fprintf(&src, "package app6d\n\n")

	fmt.Fprintf(&src, "// entityKey composes a SymbolSet and Entity into one comparable key.\n")
	fmt.Fprintf(&src, "type entityKey struct { Set SymbolSet; E Entity }\n\n")
	fmt.Fprintf(&src, "// modifier1Key composes a SymbolSet and Modifier1 into one comparable key.\n")
	fmt.Fprintf(&src, "type modifier1Key struct { Set SymbolSet; M Modifier1 }\n\n")
	fmt.Fprintf(&src, "// modifier2Key composes a SymbolSet and Modifier2 into one comparable key.\n")
	fmt.Fprintf(&src, "type modifier2Key struct { Set SymbolSet; M Modifier2 }\n\n")

	fmt.Fprintf(&src, "var entityNames = map[entityKey]string{\n")
	for _, e := range entityEntries {
		fmt.Fprintf(&src, "\t{Set: %d, E: %d}: %q,\n", e.set, e.code, e.name)
	}
	fmt.Fprintf(&src, "}\n\n")

	fmt.Fprintf(&src, "var modifier1Names = map[modifier1Key]string{\n")
	for _, m := range mod1Entries {
		fmt.Fprintf(&src, "\t{Set: %d, M: %d}: %q,\n", m.set, m.code, m.name)
	}
	fmt.Fprintf(&src, "}\n\n")

	fmt.Fprintf(&src, "var modifier2Names = map[modifier2Key]string{\n")
	for _, m := range mod2Entries {
		fmt.Fprintf(&src, "\t{Set: %d, M: %d}: %q,\n", m.set, m.code, m.name)
	}
	fmt.Fprintf(&src, "}\n")

	outFile := filepath.Join(outDir, "tables_gen.go")
	if err := writeFormatted(outFile, src.String()); err != nil {
		return fmt.Errorf("writing %s: %w", outFile, err)
	}
	log.Info("wrote lookup tables", slog.String("file", outFile),
		slog.Int("entities", len(entityEntries)),
		slog.Int("modifier1", len(mod1Entries)),
		slog.Int("modifier2", len(mod2Entries)))
	return nil
}

func setFilenameStem(set symbolSetTables) string {
	return strings.ReplaceAll(strings.ToLower(set.goName), " ", "_")
}

func escapeComment(s string) string {
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\r", "")
	return s
}

func writeFormatted(path, src string) (err error) {
	formatted, err := format.Source([]byte(src))
	if err != nil {
		return fmt.Errorf("formatting %s: %w\nsource was:\n%s", path, err, src)
	}
	return os.WriteFile(path, formatted, 0o644)
}
