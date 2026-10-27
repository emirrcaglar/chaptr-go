// SPDX-License-Identifier: GPL-3.0-or-later
// Go port and modifications of chaptr 1.5.0, 2026-10-27, by emirrcaglar.
// Upstream: https://github.com/johnthreekay/chaptr. See NOTICE for attribution.

package chaptr

import "testing"

func TestLookupTables(t *testing.T) {
	for _, entry := range languageTags {
		if got := LookupLanguage(entry.tag); got != entry.language {
			t.Errorf("%q: got %s, want %s", entry.tag, got, entry.language)
		}
	}
	if LookupLanguage("FRANÇAIS") != French || LookupLanguage("РУССКИЙ") != Russian || LookupLanguage("Raw") != "" || LookupLanguage("engineering") != "" {
		t.Fatal("language matching regression")
	}
	if LookupPublisher("J-NOVEL CLUB") != "J-Novel Club" || LookupScanner("LUCAZ") != "LuCaZ" || LookupScanner("nobody") != "" || LookupPublisher("Some Press") != "" {
		t.Fatal("canonical credit lookup regression")
	}
}
