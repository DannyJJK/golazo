package worldcup

import (
	"sort"
	"testing"
)

// alternateFlagCodes lists flagEmojis keys that intentionally have no
// matching wcNameToCode entry: they exist purely so FotMob payloads that
// surface alternate codes (e.g. "IRI" for Iran, "HOL" for Netherlands)
// still render with a flag. The canonical name → code mapping lives in
// wcNameToCode under the FIFA-standard code.
var alternateFlagCodes = map[string]bool{
	"HOL": true, // Netherlands alternate
	"IRI": true, // Iran alternate
	"GBR": true, // Great Britain (no corresponding national team in WC)
}

// unflaggedNameCodes lists wcNameToCode values that intentionally have no
// flagEmojis entry because the country has no representable Unicode flag
// emoji. Northern Ireland is the canonical example: there is no Unicode
// regional indicator or subdivision tag for it (only ENG, SCT, WLS exist).
// Teams in this list still get their FIFA code rendered, just without an
// emoji prefix.
var unflaggedNameCodes = map[string]bool{
	"NIR": true, // Northern Ireland — no Unicode flag exists
}

// TestFlagCoverage_NameOverridesHaveFlags asserts every code that
// wcNameToCode maps to is also a key in flagEmojis. Without this, adding
// a new country to wcNameToCode without its flag would render as
// "   XYZ" silently, with no compile-time or runtime warning.
func TestFlagCoverage_NameOverridesHaveFlags(t *testing.T) {
	var missing []string
	for name, code := range wcNameToCode {
		if unflaggedNameCodes[code] {
			continue
		}
		if _, ok := flagEmojis[code]; !ok {
			missing = append(missing, name+" → "+code)
		}
	}
	if len(missing) > 0 {
		sort.Strings(missing)
		t.Errorf("wcNameToCode entries missing a matching flagEmojis flag (%d):\n  %s\n"+
			"add the flag to flagEmojis, or add the code to unflaggedNameCodes if no Unicode flag exists",
			len(missing), joinNewline(missing))
	}
}

// TestFlagCoverage_NoOrphanedFlagEntries asserts every flagEmojis key
// (except the documented alternate codes) is reachable from at least one
// wcNameToCode value. This prevents flag entries from drifting out of
// sync — e.g. a flag added without the corresponding name override would
// never render in any World Cup view.
func TestFlagCoverage_NoOrphanedFlagEntries(t *testing.T) {
	reachable := make(map[string]bool, len(wcNameToCode))
	for _, code := range wcNameToCode {
		reachable[code] = true
	}

	var orphans []string
	for code := range flagEmojis {
		if alternateFlagCodes[code] {
			continue
		}
		if !reachable[code] {
			orphans = append(orphans, code)
		}
	}
	if len(orphans) > 0 {
		sort.Strings(orphans)
		t.Errorf("flagEmojis entries with no wcNameToCode mapping (%d):\n  %s\n"+
			"add the country name to wcNameToCode, or add the code to alternateFlagCodes if intentional",
			len(orphans), joinNewline(orphans))
	}
}

func joinNewline(items []string) string {
	out := ""
	for i, s := range items {
		if i > 0 {
			out += "\n  "
		}
		out += s
	}
	return out
}
