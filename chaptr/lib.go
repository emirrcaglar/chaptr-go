// SPDX-License-Identifier: GPL-3.0-or-later
// Go port and modifications of chaptr 1.5.0, 2026-10-27, by emirrcaglar.
// Upstream: https://github.com/johnthreekay/chaptr. See NOTICE for attribution.

package chaptr

// ChapterNumber stores a number without floating point conversion. Decimal is
// optional; nil sorts before a present zero. Decimal digits are stored as an
// integer, exactly as in the Rust model (42.05 and 42.5 both store decimal 5).
type ChapterNumber struct {
	Whole   uint32  `json:"whole"`
	Decimal *uint16 `json:"decimal"`
}

// NewChapterNumber constructs a number without a decimal component.
func NewChapterNumber(whole uint32) ChapterNumber { return ChapterNumber{Whole: whole} }

// WithDecimal constructs a number with an explicit decimal component.
func WithDecimal(whole uint32, decimal uint16) ChapterNumber {
	return ChapterNumber{Whole: whole, Decimal: &decimal}
}

// Compare orders numbers lexicographically by whole part, then optional decimal.
func (n ChapterNumber) Compare(other ChapterNumber) int {
	if n.Whole < other.Whole {
		return -1
	}
	if n.Whole > other.Whole {
		return 1
	}
	if n.Decimal == nil {
		if other.Decimal == nil {
			return 0
		}
		return -1
	}
	if other.Decimal == nil {
		return 1
	}
	if *n.Decimal < *other.Decimal {
		return -1
	}
	if *n.Decimal > *other.Decimal {
		return 1
	}
	return 0
}

// NumberRange is a single number or an inclusive range.
type NumberRange struct {
	Start ChapterNumber  `json:"start"`
	End   *ChapterNumber `json:"end"`
}

// Single constructs a range containing just n.
func Single(n ChapterNumber) NumberRange { return NumberRange{Start: n} }

// Range constructs an inclusive range.
func Range(start, end ChapterNumber) NumberRange { return NumberRange{Start: start, End: &end} }

// IsRange reports whether an end number is present.
func (r NumberRange) IsRange() bool { return r.End != nil }

// Compare orders ranges by start, then optional end, using ChapterNumber.Compare.
func (r NumberRange) Compare(other NumberRange) int {
	if c := r.Start.Compare(other.Start); c != 0 {
		return c
	}
	if r.End == nil {
		if other.End == nil {
			return 0
		}
		return -1
	}
	if other.End == nil {
		return 1
	}
	return r.End.Compare(*other.End)
}

// Language identifies an explicitly declared language tag. Empty means absent.
type Language string

const (
	English            Language = "English"
	Japanese           Language = "Japanese"
	SimplifiedChinese  Language = "SimplifiedChinese"
	TraditionalChinese Language = "TraditionalChinese"
	Korean             Language = "Korean"
	Spanish            Language = "Spanish"
	French             Language = "French"
	German             Language = "German"
	Italian            Language = "Italian"
	Portuguese         Language = "Portuguese"
	Russian            Language = "Russian"
	Indonesian         Language = "Indonesian"
	Vietnamese         Language = "Vietnamese"
	Thai               Language = "Thai"
)

// MangaSource identifies the release's source. Empty means absent.
type MangaSource string

const (
	Digital   MangaSource = "Digital"
	MangaPlus MangaSource = "MangaPlus"
	Viz       MangaSource = "Viz"
	Kodansha  MangaSource = "Kodansha"
	Scan      MangaSource = "Scan"
	Lezhin    MangaSource = "Lezhin"
	Naver     MangaSource = "Naver"
	Kakao     MangaSource = "Kakao"
)

// ParsedManga contains the recognized fields. Empty strings mean absent;
// pointers preserve optional numbers, including explicit zero values.
type ParsedManga struct {
	Title     string       `json:"title"`
	Volume    *NumberRange `json:"volume"`
	Chapter   *NumberRange `json:"chapter"`
	Group     string       `json:"group"`
	Source    MangaSource  `json:"source"`
	Language  Language     `json:"language"`
	Revision  *uint8       `json:"revision"`
	IsOneshot bool         `json:"is_oneshot"`
	Edition   string       `json:"edition"`
	Extension string       `json:"extension"`
}

// ParsedNovel contains the recognized light novel release fields.
type ParsedNovel struct {
	Title     string       `json:"title"`
	Volume    *NumberRange `json:"volume"`
	Group     string       `json:"group"`
	Publisher string       `json:"publisher"`
	Scanner   string       `json:"scanner"`
	Language  Language     `json:"language"`
	Year      *uint16      `json:"year"`
	IsPremium bool         `json:"is_premium"`
	IsDigital bool         `json:"is_digital"`
	Revision  *uint8       `json:"revision"`
	Extension string       `json:"extension"`
}
