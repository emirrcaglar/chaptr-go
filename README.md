# chaptr for Go

A dependency-free Go port of [Rust `chaptr` 1.5.0](https://github.com/johnthreekay/chaptr). Parses
manga, manhwa, manhua, and light novel release filenames without I/O.
Requires Go 1.22 or newer.

After the source is pushed to GitHub, install the module in your project:

```sh
go get github.com/emirrcaglar/chaptr-go/chaptr@latest
```

```go
import "github.com/emirrcaglar/chaptr-go/chaptr"

m := chaptr.ParseManga("[Group] Title v03 c042.5 (Digital).cbz")
// m.Title == "Title"
// m.Volume.Start.Whole == 3
// m.Chapter.Start.Whole == 42, *m.Chapter.Start.Decimal == 5
// m.Group == "Group", m.Source == chaptr.Digital

n := chaptr.ParseNovel("Some Light Novel v05 (Yen Press) (Digital) [LuCaZ].epub")
// n.Publisher == "Yen Press", n.Scanner == "LuCaZ", n.IsDigital == true
```

The module path is `github.com/emirrcaglar/chaptr-go`. The public library
package is `github.com/emirrcaglar/chaptr-go/chaptr`, named `chaptr` in Go code.

Publish the contents of this directory at the root of
[emirrcaglar/chaptr-go](https://github.com/emirrcaglar/chaptr-go), with `go.mod`,
`LICENSE`, and `NOTICE` at the repository root. Keep `chaptr/`, `cmd/`,
`LICENSES/`, and `tools/`; `.build/` is local build output and is
gitignored. The Rust reference source can remain outside this module.

## Project layout

Library source and its tests live in `chaptr/`. The CLI lives in
`cmd/chaptr/`, while fixture-generation tools live in `tools/`. The repository
root contains module metadata, license files, and documentation.

```text
chaptr-go/
├── go.mod
├── chaptr/                 # Importable library package
│   ├── doc.go              # Package documentation
│   ├── lib.go              # Parsed results, numbers, ranges, and enums
│   ├── lexer.go            # Shared tokenizer
│   ├── common.go           # Shared volume, chapter, and title detection
│   ├── manga.go            # Manga parser
│   ├── novel.go            # Light novel parser
│   ├── languages.go        # Language lookups
│   ├── ln_publishers.go    # Light novel publisher lookups
│   ├── ln_scanners.go      # Light novel scanner lookups
│   ├── *_test.go           # Unit tests, examples, fuzzing, and benchmarks
│   └── testdata/           # Rust comparison fixtures and attribution
├── cmd/chaptr/main.go      # Filename CLI
├── tools/                  # Fixture generator and Rust oracle
├── LICENSES/               # Additional upstream license texts
├── LICENSE
├── NOTICE
└── README.md
```

Executable examples in `chaptr/example_test.go` use the public package import,
while `chaptr/parse_test.go` and `chaptr/fuzz_test.go` keep performance and robustness checks
separate from the unit tests. The test files remain next to the library code,
and fixture files stay under `chaptr/testdata/`.

## Try a filename

From this directory, pass a quoted filename to the command to print its metadata:

```sh
go run ./cmd/chaptr "Blue Lock v39 (2026) (Digital) (Stick).cbz"
go run ./cmd/chaptr -novel "Some Light Novel v05 (Yen Press) (Digital) [LuCaZ].epub"
```

`go test -run` filters test names; use `go run` above to parse your own filename.

## API and Rust mapping

| Rust | Go |
|---|---|
| `manga::parse` | `chaptr.ParseManga` |
| `novel::parse` | `chaptr.ParseNovel` |
| `lexer::tokenize` | `chaptr.Tokenize` |
| `detect_volume`, `detect_chapter` | `chaptr.DetectVolume`, `chaptr.DetectChapter` |
| `detect_chapter_revision`, `detect_language` | `chaptr.DetectChapterRevision`, `chaptr.DetectLanguage` |
| `tables::*::lookup` | `chaptr.LookupPublisher`, `chaptr.LookupScanner`, `chaptr.LookupLanguage` |
| `ChapterNumber::new`, `with_decimal` | `chaptr.NewChapterNumber`, `chaptr.WithDecimal` |
| `NumberRange::single`, `range`, `is_range` | `chaptr.Single`, `chaptr.Range`, `NumberRange.IsRange` |

Missing text and enum fields use empty strings. Optional numbers and ranges
use pointers, so absent values remain distinguishable from explicit zero.
Use `ChapterNumber.Compare` and `NumberRange.Compare` for value equality and
sorting; Go's `==` compares the pointer fields by address.

Number components retain Rust's lexicographic ordering and integer decimal
representation: absent decimal sorts before a present zero; `.05` and `.5`
both store decimal `5`. Backward ranges collapse to their starting number.

Tokens include byte offsets and preserve input slices. Titles allocate when
underscore normalization is needed. Unclosed brackets consume the remaining
input. The lexer also tolerates invalid UTF-8 in Go strings.

The source behavior is preserved, including conservative bare chapter
heuristics, top-level chapter revision detection, the existing manga group
skip list, the first year of novel year ranges, and the absence of generic
`(Scan)` source detection. Thai volume markers and X-suffix range syntax retain
the upstream limitations. Unicode character classification uses Go's Unicode
tables, so characters added in a newer Rust Unicode version may differ.

## Validation

Run from this directory:

```sh
go test ./...
go vet ./...
go test ./chaptr -fuzz=FuzzParse -fuzztime=10s
go test ./chaptr -run='^$' -bench=. -benchmem
```

`chaptr/testdata/rust_parity.json` records lexer output and every manga/novel field
from Rust for all unique filenames in the four JSON corpora and smoke corpus,
plus upstream unit-test string literals and numeric boundary cases. The Go
test compares all fields and token spans against these saved outputs, including
upstream scope gaps. Normal Go builds and tests do not require Rust or Python.
To regenerate the fixtures with Python 3 and `rustc` installed:

```sh
python tools/generate_parity.py
```

Regeneration requires the upstream Rust source and corpus in the sibling
`../rust/` directory. Keep that layout when checking out the complete workspace;
if publishing just the Go module, obtain the matching upstream source separately
before running the generator. Normal library use and Go tests need only this module.

## License and attribution

GPL-3.0-or-later, matching the Rust source. See [LICENSE](LICENSE) and
[NOTICE](NOTICE). Corpus-derived fixtures retain their upstream attribution;
the Mihon-derived cases also carry Apache-2.0 licensing information in
`LICENSES/Apache-2.0.txt`.

Port and modifications dated 2026-10-27 are described in `NOTICE`. Corpus
provenance is retained in `chaptr/testdata/NOTICE`.
