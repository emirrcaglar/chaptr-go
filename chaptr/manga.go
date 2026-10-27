// SPDX-License-Identifier: GPL-3.0-or-later
// Go port and modifications of chaptr 1.5.0, 2026-10-27, by emirrcaglar.
// Upstream: https://github.com/johnthreekay/chaptr. See NOTICE for attribution.

package chaptr

import "strings"

var mangaExtensions = []string{"cbz", "cbr", "zip", "7z", "rar", "pdf"}
var formatTags = []string{"full color", "colored", "color", "b&w", "bw", "hd", "hi-res", "uhd", "raw", "raws", "manga", "audiobook", "english", "japanese", "chinese", "korean", "french", "spanish", "german", "italian"}
var editionKeywords = []string{"omnibus", "uncensored"}
var sourcePatterns = []struct {
	pattern string
	source  MangaSource
}{
	{"manga plus", MangaPlus}, {"mangaplus", MangaPlus}, {"viz media", Viz}, {"viz", Viz},
	{"kodansha usa", Kodansha}, {"kodansha", Kodansha}, {"lezhin comics", Lezhin}, {"lezhin", Lezhin},
	{"naver webtoon", Naver}, {"naver series", Naver}, {"naver", Naver}, {"kakao page", Kakao}, {"kakao", Kakao}, {"digital", Digital},
}

// ParseManga parses a manga, manhwa, or manhua filename. It never performs I/O.
func ParseManga(filename string) ParsedManga {
	ts := Tokenize(filename)
	v, c := DetectVolume(ts), DetectChapter(ts)
	return ParsedManga{
		Title: detectTitle(filename, ts, mangaExtensions), Volume: v, Chapter: c,
		Group: mangaGroup(ts), Source: detectSource(ts), Language: DetectLanguage(ts),
		Revision: DetectChapterRevision(ts), IsOneshot: detectOneshot(ts, v != nil, c != nil),
		Edition: detectEdition(ts), Extension: detectExtension(filename, mangaExtensions),
	}
}
func sourceFromContent(s string) MangaSource {
	s = lowerASCII(s)
	for _, p := range sourcePatterns {
		if s == p.pattern || (strings.HasPrefix(s, p.pattern) && len(s) > len(p.pattern) && (s[len(p.pattern)] == ' ' || s[len(p.pattern)] == '-')) {
			return p.source
		}
	}
	return ""
}
func mangaGroup(ts []Token) string {
	for _, t := range ts {
		if t.Kind != Bracketed || t.Text == "" {
			continue
		}
		s := lowerASCII(t.Text)
		if !containsVolumeKeyword(t.Text) && !has(editionKeywords, s) && !has(formatTags, s) && sourceFromContent(t.Text) == "" {
			return t.Text
		}
	}
	return ""
}
func detectSource(ts []Token) MangaSource {
	for _, t := range ts {
		if bracket(t) {
			if s := sourceFromContent(t.Text); s != "" {
				return s
			}
		}
	}
	return ""
}
func nextWord(ts []Token, i int) string {
	for ; i < len(ts); i++ {
		if ts[i].Kind == Delimiter {
			continue
		}
		if ts[i].Kind == Word {
			return ts[i].Text
		}
		return ""
	}
	return ""
}
func detectEdition(ts []Token) string {
	for i, t := range ts {
		if t.Kind != Word || !has(editionKeywords, lowerASCII(t.Text)) {
			continue
		}
		if next := nextWord(ts, i+1); equalASCII(next, "edition") {
			return t.Text + " " + next
		}
		return t.Text
	}
	return ""
}
func detectOneshot(ts []Token, volume, chapter bool) bool {
	for _, t := range ts {
		if t.Kind == Delimiter {
			continue
		}
		if t.Kind == Bracketed && len(t.Text) >= 5 && allDigits(t.Text) {
			return true
		}
		break
	}
	if volume || chapter {
		return false
	}
	for i, t := range ts {
		if t.Kind != Word {
			continue
		}
		w := lowerASCII(t.Text)
		if w == "special" || w == "specials" {
			return true
		}
		next := lowerASCII(nextWord(ts, i+1))
		if w == "art" && next == "collection" || w == "all" && next == "stars" {
			return true
		}
	}
	return false
}
