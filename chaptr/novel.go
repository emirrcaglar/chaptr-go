// SPDX-License-Identifier: GPL-3.0-or-later
// Go port and modifications of chaptr 1.5.0, 2026-10-27, by emirrcaglar.
// Upstream: https://github.com/johnthreekay/chaptr. See NOTICE for attribution.

package chaptr

var novelExtensions = []string{"epub", "pdf", "azw3", "mobi", "txt"}

// ParseNovel parses a light novel filename. Year ranges retain their first year.
func ParseNovel(filename string) ParsedNovel {
	ts := Tokenize(filename)
	n := ParsedNovel{Title: detectTitle(filename, ts, novelExtensions), Volume: DetectVolume(ts), Group: novelGroup(ts), Language: DetectLanguage(ts), Extension: detectExtension(filename, novelExtensions)}
	for _, t := range ts {
		if bracket(t) {
			if n.Publisher == "" {
				n.Publisher = LookupPublisher(t.Text)
			}
			if n.Scanner == "" {
				n.Scanner = LookupScanner(t.Text)
			}
			if equalASCII(t.Text, "digital") {
				n.IsDigital = true
			}
			if equalASCII(t.Text, "premium") {
				n.IsPremium = true
			}
			if n.Year == nil {
				s := t.Text
				if y, ok := parseUint(s[:digitsEnd(s)], 32); ok && looksLikeYear(y) {
					n.Year = ptr(uint16(y))
				}
			}
		}
		if t.Kind == Curly && n.Revision == nil && len(t.Text) > 1 && equalASCII(t.Text[:1], "r") && allDigits(t.Text[1:]) {
			if r, ok := parseUint(t.Text[1:], 8); ok {
				n.Revision = ptr(uint8(r))
			}
		}
	}
	return n
}
func novelGroup(ts []Token) string {
	for _, t := range ts {
		if t.Kind == Delimiter {
			continue
		}
		if t.Kind != Bracketed {
			return ""
		}
		if t.Text == "" {
			continue
		}
		if LookupPublisher(t.Text) != "" || LookupScanner(t.Text) != "" || containsVolumeKeyword(t.Text) {
			return ""
		}
		return t.Text
	}
	return ""
}
