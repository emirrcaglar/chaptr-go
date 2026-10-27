// SPDX-License-Identifier: GPL-3.0-or-later
// Go port and modifications of chaptr 1.5.0, 2026-10-27, by emirrcaglar.
// Upstream: https://github.com/johnthreekay/chaptr. See NOTICE for attribution.

package chaptr

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
)

// The fixture records actual Rust outputs, including its intentional scope gaps.
// Regenerate with python tools/generate_parity.py from this module directory.
func TestRustParity(t *testing.T) {
	data, err := os.ReadFile("testdata/rust_parity.json")
	if err != nil {
		t.Fatal(err)
	}
	var cases []struct {
		Input  string      `json:"input"`
		Manga  ParsedManga `json:"manga"`
		Novel  ParsedNovel `json:"novel"`
		Tokens []Token     `json:"tokens"`
	}
	if err = json.Unmarshal(data, &cases); err != nil {
		t.Fatal(err)
	}
	if len(cases) < 1000 {
		t.Fatalf("unexpectedly small parity corpus: %d", len(cases))
	}
	for _, c := range cases {
		t.Run(c.Input, func(t *testing.T) {
			if got := Tokenize(c.Input); !reflect.DeepEqual(got, c.Tokens) {
				t.Errorf("tokens: got %+v; want %+v", got, c.Tokens)
			}
			if got := ParseManga(c.Input); !reflect.DeepEqual(got, c.Manga) {
				a, _ := json.Marshal(got)
				b, _ := json.Marshal(c.Manga)
				t.Errorf("manga: got %s; want %s", a, b)
			}
			if got := ParseNovel(c.Input); !reflect.DeepEqual(got, c.Novel) {
				a, _ := json.Marshal(got)
				b, _ := json.Marshal(c.Novel)
				t.Errorf("novel: got %s; want %s", a, b)
			}
		})
	}
}
