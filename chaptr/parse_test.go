// SPDX-License-Identifier: GPL-3.0-or-later
// Go port and modifications of chaptr 1.5.0, 2026-10-27, by emirrcaglar.
// Upstream: https://github.com/johnthreekay/chaptr. See NOTICE for attribution.

package chaptr

import "testing"

var benchmarkManga ParsedManga
var benchmarkNovel ParsedNovel

func BenchmarkManga(b *testing.B) {
	for _, tc := range []struct{ name, input string }{
		{"simple", "BTOOOM! v01 (2013) (Digital) (Shadowcat-Empire)"},
		{"full", "[Yen Press] Sword Art Online v10 c042.5 (Digital).cbz"},
		{"cjk", "幽游白书完全版 第03卷 天下"},
		{"range", "Historys Strongest Disciple Kenichi_v11_c90-98.zip"},
	} {
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				benchmarkManga = ParseManga(tc.input)
			}
		})
	}
}
func BenchmarkNovel(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchmarkNovel = ParseNovel("[Unpaid Ferryman] Youjo Senki v01-23 (2018-2024) (Digital) (LuCaZ)")
	}
}
