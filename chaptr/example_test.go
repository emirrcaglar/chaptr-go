// SPDX-License-Identifier: GPL-3.0-or-later
// Go port and modifications of chaptr 1.5.0, 2026-10-27, by emirrcaglar.
// Upstream: https://github.com/johnthreekay/chaptr. See NOTICE for attribution.

package chaptr_test

import (
	"fmt"

	"github.com/emirrcaglar/chaptr-go/chaptr"
)

func ExampleParseManga() {
	m := chaptr.ParseManga("[Group] Title v03 c042.5 (Digital).cbz")
	fmt.Println(m.Title, m.Volume.Start.Whole, m.Chapter.Start.Whole, *m.Chapter.Start.Decimal, m.Source)
	// Output: Title 3 42 5 Digital
}

func ExampleParseNovel() {
	n := chaptr.ParseNovel("[Unpaid Ferryman] Youjo Senki v01-23 (2018-2024) (Digital) (LuCaZ)")
	fmt.Println(n.Title, n.Volume.Start.Whole, n.Volume.End.Whole, *n.Year, n.Scanner)
	// Output: Youjo Senki 1 23 2018 LuCaZ
}
