// SPDX-License-Identifier: GPL-3.0-or-later
// Go port and modifications of chaptr 1.5.0, 2026-10-27, by emirrcaglar.
// Upstream: https://github.com/johnthreekay/chaptr. See NOTICE for attribution.

package chaptr

import (
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

var volumeKeywords = []string{"vol", "volume", "tome", "том", "тома"}
var chapterKeywords = []string{"ch", "chp", "chapter", "chapters", "episode", "глава", "главы"}

const volumeMarkers = "巻卷册권장"
const chapterMarkers = "话話章回회화"

func ptr[T any](v T) *T { return &v }
func lowerASCII(s string) string {
	for i := 0; i < len(s); i++ {
		if s[i] >= 'A' && s[i] <= 'Z' {
			b := []byte(s)
			for j := i; j < len(b); j++ {
				if b[j] >= 'A' && b[j] <= 'Z' {
					b[j] += 32
				}
			}
			return string(b)
		}
	}
	return s
}
func equalASCII(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		x, y := a[i], b[i]
		if x >= 'A' && x <= 'Z' {
			x += 32
		}
		if y >= 'A' && y <= 'Z' {
			y += 32
		}
		if x != y {
			return false
		}
	}
	return true
}
func has(ss []string, s string) bool {
	for _, v := range ss {
		if v == s {
			return true
		}
	}
	return false
}
func digitsEnd(s string) int {
	i := 0
	for i < len(s) && s[i] >= '0' && s[i] <= '9' {
		i++
	}
	return i
}
func allDigits(s string) bool { return len(s) > 0 && digitsEnd(s) == len(s) }
func parseUint(s string, bits int) (uint32, bool) {
	// Rust accepts a leading plus in integer parsing.
	n, err := strconv.ParseUint(strings.TrimPrefix(s, "+"), 10, bits)
	return uint32(n), err == nil
}
func wordAt(ts []Token, i int) (string, bool) {
	if i >= 0 && i < len(ts) && ts[i].Kind == Word {
		return ts[i].Text, true
	}
	return "", false
}
func delimAt(ts []Token, i int, s string) bool {
	return i >= 0 && i < len(ts) && ts[i].Kind == Delimiter && ts[i].Text == s
}
func bracket(t Token) bool { return t.Kind == Bracketed || t.Kind == Parenthesized }
func keywordMatch(w string, keywords []string) bool {
	for i := 0; i < len(w); i++ {
		if w[i] >= utf8.RuneSelf {
			return has(keywords, strings.ToLower(w))
		}
	}
	for _, k := range keywords {
		if equalASCII(w, k) {
			return true
		}
	}
	return false
}
func isVolumeKeyword(w string) bool  { return keywordMatch(w, volumeKeywords) }
func isChapterKeyword(w string) bool { return keywordMatch(w, chapterKeywords) }
func looksLikeYear(n uint32) bool    { return n >= 1900 && n <= 2099 }
func containsVolumeKeyword(s string) bool {
	s = lowerASCII(s)
	return strings.Contains(s, "volume") || strings.HasPrefix(s, "vol ") || strings.HasPrefix(s, "vol.")
}
func stripVolumePrefix(w string) (string, bool) {
	for _, p := range []string{"volume", "vol", "v", "sp", "s", "t"} {
		if len(w) >= len(p) && equalASCII(w[:len(p)], p) && allDigits(w[len(p):]) {
			return w[len(p):], true
		}
	}
	return "", false
}
func cjkPrefixVolume(w string) (uint32, bool) {
	if !strings.HasPrefix(w, "시즌") {
		return 0, false
	}
	s := strings.TrimPrefix(w, "시즌")
	return parseUint(s[:digitsEnd(s)], 32)
}
func chapterWord(w string) (uint32, bool) {
	for _, p := range []string{"chapter", "chp", "ch", "c"} {
		if len(w) < len(p) || !equalASCII(w[:len(p)], p) {
			continue
		}
		s := w[len(p):]
		i := digitsEnd(s)
		if i == 0 {
			continue
		}
		n, ok := parseUint(s[:i], 32)
		if !ok {
			continue
		}
		after := s[i:]
		if after == "" || (len(after) > 1 && equalASCII(after[:1], "v") && allDigits(after[1:])) {
			return n, true
		}
	}
	return 0, false
}

// DetectChapterRevision extracts a top-level chapter suffix such as c042v2.
func DetectChapterRevision(ts []Token) *uint8 {
	for _, t := range ts {
		if t.Kind != Word {
			continue
		}
		w := t.Text
		for _, p := range []string{"chapter", "chp", "ch", "c"} {
			if len(w) < len(p) || !equalASCII(w[:len(p)], p) {
				continue
			}
			s := w[len(p):]
			i := digitsEnd(s)
			if i == 0 {
				continue
			}
			s = s[i:]
			if len(s) > 1 && equalASCII(s[:1], "v") && allDigits(s[1:]) {
				if n, ok := parseUint(s[1:], 8); ok {
					return ptr(uint8(n))
				}
			}
		}
	}
	return nil
}

func extendedNumber(w string) *ChapterNumber {
	i := digitsEnd(w)
	if i == 0 {
		return nil
	}
	n, ok := parseUint(w[:i], 32)
	if !ok {
		return nil
	}
	if len(w)-i == 1 && w[i] >= 'a' && w[i] <= 'z' {
		return ptr(WithDecimal(n, 5))
	}
	return ptr(NewChapterNumber(n))
}
func decimalSuffix(ts []Token, i int) (uint16, int, bool) {
	if !delimAt(ts, i, ".") {
		return 0, i, false
	}
	w, ok := wordAt(ts, i+1)
	if !ok {
		return 0, i, false
	}
	n, ok := parseUint(w, 16)
	return uint16(n), i + 2, ok
}
func buildRange(ts []Token, i int, n uint32, chapter bool) *NumberRange {
	start := NewChapterNumber(n)
	if d, next, ok := decimalSuffix(ts, i); ok {
		start = WithDecimal(n, d)
		i = next
	}
	r := Single(start)
	if !delimAt(ts, i, "-") {
		return &r
	}
	w, ok := wordAt(ts, i+1)
	if !ok {
		return &r
	}
	end := extendedNumber(w)
	if end == nil && chapter {
		if n, ok := chapterWord(w); ok {
			end = ptr(NewChapterNumber(n))
		}
	}
	if end == nil {
		return &r
	}
	if end.Decimal == nil {
		if d, _, ok := decimalSuffix(ts, i+2); ok {
			*end = WithDecimal(end.Whole, d)
		}
	}
	if end.Compare(start) >= 0 {
		r.End = end
	}
	return &r
}
func numberAfterMarker(ts []Token, i int, chapter bool) *NumberRange {
	for skipped := 0; i < len(ts) && skipped < 3; {
		t := ts[i]
		if t.Kind == Delimiter {
			skipped++
			i++
			continue
		}
		if t.Kind != Word {
			return nil
		}
		n, ok := parseUint(t.Text, 32)
		if !ok {
			return nil
		}
		return buildRange(ts, i+1, n, chapter)
	}
	return nil
}
func volumeInSeq(ts []Token) *NumberRange {
	for i, t := range ts {
		if t.Kind != Word {
			continue
		}
		w := t.Text
		if s, ok := stripVolumePrefix(w); ok {
			if n, ok := parseUint(s, 32); ok {
				return buildRange(ts, i+1, n, false)
			}
		}
		if n, ok := cjkPrefixVolume(w); ok {
			return buildRange(ts, i+1, n, false)
		}
		if !isVolumeKeyword(w) {
			continue
		}
		if r := numberAfterMarker(ts, i+1, false); r != nil {
			return r
		}
		if lower := strings.ToLower(w); lower == "том" || lower == "тома" {
			skipped := 0
			for j := i - 1; j >= 0; j-- {
				if ts[j].Kind == Delimiter {
					skipped++
					if skipped > 3 {
						break
					}
					continue
				}
				if ts[j].Kind == Word {
					if n, ok := parseUint(ts[j].Text, 32); ok && !looksLikeYear(n) {
						return ptr(Single(NewChapterNumber(n)))
					}
				}
				break
			}
		}
	}
	return nil
}
func chapterInSeq(ts []Token) *NumberRange {
	for i, t := range ts {
		if t.Kind != Word {
			continue
		}
		if n, ok := chapterWord(t.Text); ok {
			return buildRange(ts, i+1, n, true)
		}
		if isChapterKeyword(t.Text) {
			if r := numberAfterMarker(ts, i+1, true); r != nil {
				return r
			}
		}
	}
	return nil
}

// DetectVolume recognizes explicit volume markers, including bracketed markers.
func DetectVolume(ts []Token) *NumberRange {
	if r := volumeInSeq(ts); r != nil {
		return r
	}
	for _, t := range ts {
		if bracket(t) {
			if r := volumeInSeq(Tokenize(t.Text)); r != nil {
				return r
			}
		}
	}
	return cjkSeq(ts, volumeMarkers)
}

// DetectChapter recognizes explicit chapters, CJK markers, and conservative bare numbers.
func DetectChapter(ts []Token) *NumberRange {
	if r := hashChapter(ts); r != nil {
		return r
	}
	if r := chapterInSeq(ts); r != nil {
		return r
	}
	for _, t := range ts {
		if bracket(t) {
			if r := chapterInSeq(Tokenize(t.Text)); r != nil {
				return r
			}
		}
	}
	if r := cjkSeq(ts, chapterMarkers); r != nil {
		return r
	}
	return bareChapter(ts)
}
func hashChapter(ts []Token) *NumberRange {
	hash := false
	for _, t := range ts {
		if t.Kind == Delimiter && t.Text == "#" {
			hash = true
			break
		}
	}
	if !hash {
		return nil
	}
	for _, t := range ts {
		if t.Kind == Word {
			if _, ok := cjkNumber(t.Text, chapterMarkers); ok {
				return nil
			}
		}
	}
	for i := range ts {
		if !delimAt(ts, i, "#") {
			continue
		}
		w, ok := wordAt(ts, i+1)
		if !ok {
			continue
		}
		n, ok := parseUint(w, 32)
		if !ok {
			continue
		}
		end := true
		for _, t := range ts[i+2:] {
			if t.Kind != Delimiter && (t.Kind != Word || !looksLikeExtension(t.Text)) {
				end = false
				break
			}
		}
		if end {
			return ptr(Single(NewChapterNumber(n)))
		}
	}
	return nil
}
func cjkNumber(w, markers string) (uint32, bool) {
	for i, c := range w {
		if !strings.ContainsRune(markers, c) {
			continue
		}
		j := i
		for j > 0 && w[j-1] >= '0' && w[j-1] <= '9' {
			j--
		}
		if j < i {
			if n, ok := parseUint(w[j:i], 32); ok {
				return n, true
			}
		}
		j = i + utf8.RuneLen(c)
		k := j + digitsEnd(w[j:])
		if k > j {
			if n, ok := parseUint(w[j:k], 32); ok {
				return n, true
			}
		}
	}
	return 0, false
}
func cjkSeq(ts []Token, markers string) *NumberRange {
	for i, t := range ts {
		if t.Kind != Word {
			continue
		}
		n, ok := cjkNumber(t.Text, markers)
		if !ok {
			continue
		}
		if i >= 2 {
			if w, ok := wordAt(ts, i-2); ok {
				if before, ok := parseUint(w, 32); ok {
					if delimAt(ts, i-1, "-") {
						if before <= n {
							return ptr(Range(NewChapterNumber(before), NewChapterNumber(n)))
						}
						return ptr(Single(NewChapterNumber(before)))
					}
					if delimAt(ts, i-1, ".") && n <= 65535 {
						return ptr(Single(WithDecimal(before, uint16(n))))
					}
				}
			}
		}
		return ptr(Single(NewChapterNumber(n)))
	}
	return nil
}
func volumeMarker(w string) bool {
	if isVolumeKeyword(w) {
		return true
	}
	if _, ok := stripVolumePrefix(w); ok {
		return true
	}
	if _, ok := cjkPrefixVolume(w); ok {
		return true
	}
	_, ok := cjkNumber(w, volumeMarkers)
	return ok
}
func structuredMarker(w string) bool {
	if volumeMarker(w) || isChapterKeyword(w) {
		return true
	}
	if _, ok := chapterWord(w); ok {
		return true
	}
	_, ok := cjkNumber(w, chapterMarkers)
	return ok
}
func bareChapter(ts []Token) *NumberRange {
	min := 0
	for i, t := range ts {
		if t.Kind == Word && volumeMarker(t.Text) {
			min = i + 1
			break
		}
	}
	saw, prev := false, false
	for i, t := range ts {
		if t.Kind == Delimiter {
			continue
		}
		if t.Kind != Word {
			prev = false
			continue
		}
		w := t.Text
		digits := allDigits(w)
		if i >= min && saw && !prev && followedByMetadata(ts, i+1) {
			if digits {
				if n, ok := parseUint(w, 32); ok && !looksLikeYear(n) {
					return buildRange(ts, i+1, n, false)
				}
			} else {
				if n := extendedNumber(w); n != nil && n.Decimal != nil && !looksLikeYear(n.Whole) {
					return ptr(Single(*n))
				}
			}
		}
		if !digits {
			saw = true
		}
		prev = isVolumeKeyword(w)
	}
	return nil
}
func looksLikeExtension(w string) bool {
	if len(w) < 2 || len(w) > 5 {
		return false
	}
	letter := false
	for _, b := range []byte(w) {
		if b >= 'a' && b <= 'z' || b >= 'A' && b <= 'Z' {
			letter = true
		} else if b < '0' || b > '9' {
			return false
		}
	}
	return letter
}
func groupCode(w string) bool {
	letter := false
	for _, b := range []byte(w) {
		if b >= 'A' && b <= 'Z' {
			letter = true
		} else if b < '0' || b > '9' {
			return false
		}
	}
	return letter
}
func followedByMetadata(ts []Token, i int) bool {
	if w, ok := wordAt(ts, i+1); ok {
		if delimAt(ts, i, ".") && (allDigits(w) || looksLikeExtension(w)) {
			return true
		}
		if delimAt(ts, i, "-") && extendedNumber(w) != nil {
			return true
		}
	}
	for ; i < len(ts); i++ {
		t := ts[i]
		if t.Kind == Delimiter {
			continue
		}
		if t.Kind != Word {
			return true
		}
		return groupCode(t.Text) || isChapterKeyword(t.Text) || isVolumeKeyword(t.Text)
	}
	return true
}

func detectExtension(s string, known []string) string {
	i := strings.LastIndexByte(s, '.')
	if i < 0 {
		return ""
	}
	ext := s[i+1:]
	if len(ext) > 0 && len(ext) <= 4 && has(known, lowerASCII(ext)) {
		return ext
	}
	return ""
}
func detectTitle(s string, ts []Token, known []string) string {
	end := len(s)
	seen, saw, prev, later := false, false, false, false
	for _, t := range ts {
		if t.Kind == Word && structuredMarker(t.Text) {
			later = true
			break
		}
	}
	found := false
	for i, t := range ts {
		if t.Kind == Word {
			seen = true
			digits := allDigits(t.Text)
			stop := structuredMarker(t.Text)
			if !later && digits && saw && !prev {
				if n, ok := parseUint(t.Text, 32); ok && !looksLikeYear(n) && followedByMetadata(ts, i+1) {
					stop = true
				}
			}
			if stop {
				end = t.Start
				found = true
				break
			}
			if !digits {
				saw = true
			}
			prev = isVolumeKeyword(t.Text)
		} else if t.Kind != Delimiter && seen {
			end = t.Start
			found = true
			break
		}
	}
	if !found && detectExtension(s, known) != "" {
		end = strings.LastIndexByte(s, '.')
	}
	start := 0
	for _, t := range ts {
		if t.Kind == Delimiter {
			continue
		}
		if t.Kind == Word {
			break
		}
		if t.Kind == Bracketed && containsVolumeKeyword(t.Text) {
			start = 0
			break
		}
		start = t.Start + 2 + len(t.Text)
	}
	if start >= end || start < 0 || end > len(s) {
		return ""
	}
	title := strings.TrimFunc(s[start:end], func(r rune) bool { return unicode.IsSpace(r) || strings.ContainsRune("-_,#:", r) })
	if strings.ContainsRune(title, '_') {
		title = strings.TrimSpace(strings.ReplaceAll(title, "_", " "))
	}
	return title
}
