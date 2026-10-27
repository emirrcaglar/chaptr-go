// SPDX-License-Identifier: GPL-3.0-or-later
// Go port and modifications of chaptr 1.5.0, 2026-10-27, by emirrcaglar.
// Upstream: https://github.com/johnthreekay/chaptr. See NOTICE for attribution.

package chaptr

import "testing"

// Go strings can contain arbitrary bytes, unlike Rust &str. Preserve those
// bytes in spans and tolerate invalid UTF-8 rather than panicking.
func FuzzParse(f *testing.F) {
	for _, s := range []string{"", "[Group] Title v03 c042.5 (Digital).cbz", "第03卷", "시즌34삽화2", "[unclosed", "\xff\x00Title 42", "Title c4294967296.65536v256"} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, s string) {
		ts := Tokenize(s)
		cursor := 0
		for _, tok := range ts {
			if tok.Start != cursor || tok.End <= tok.Start || tok.End > len(s) {
				t.Fatalf("invalid span: %+v in %q", tok, s)
			}
			if tok.Kind == Word || tok.Kind == Delimiter {
				if tok.Text != s[tok.Start:tok.End] {
					t.Fatal("token text lost input bytes")
				}
			}
			cursor = tok.End
		}
		if cursor != len(s) {
			t.Fatal("tokens do not cover the input")
		}
		m, n := ParseManga(s), ParseNovel(s)
		for _, r := range []*NumberRange{m.Volume, m.Chapter, n.Volume} {
			if r != nil && r.End != nil && r.Start.Compare(*r.End) > 0 {
				t.Fatal("backward range")
			}
		}
	})
}
