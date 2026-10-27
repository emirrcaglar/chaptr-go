// SPDX-License-Identifier: GPL-3.0-or-later
// Go port and modifications of chaptr 1.5.0, 2026-10-27, by emirrcaglar.
// Upstream: https://github.com/johnthreekay/chaptr. See NOTICE for attribution.

package chaptr

import "testing"

func TestNumberOrdering(t *testing.T) {
	numbers := []ChapterNumber{NewChapterNumber(42), WithDecimal(42, 0), WithDecimal(42, 5), NewChapterNumber(43)}
	for i, a := range numbers {
		for j, b := range numbers {
			want := 0
			if i < j {
				want = -1
			}
			if i > j {
				want = 1
			}
			if got := a.Compare(b); got != want {
				t.Fatalf("Compare(%+v,%+v)=%d, want %d", a, b, got, want)
			}
		}
	}
	a, b := WithDecimal(42, 5), WithDecimal(42, 5)
	if a.Compare(b) != 0 {
		t.Fatal("decimal equality depends on pointer identity")
	}
	single, ranged := Single(a), Range(a, NewChapterNumber(43))
	if single.IsRange() || !ranged.IsRange() || single.Compare(ranged) >= 0 || ranged.Compare(single) <= 0 || ranged.Compare(Range(b, NewChapterNumber(43))) != 0 {
		t.Fatal("incorrect range ordering")
	}
}
