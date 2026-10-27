// SPDX-License-Identifier: GPL-3.0-or-later
// Go port and modifications of chaptr 1.5.0, 2026-10-27, by emirrcaglar.
// Upstream: https://github.com/johnthreekay/chaptr. See NOTICE for attribution.

package chaptr

import "strings"

var languageTags = []struct {
	tag      string
	language Language
}{
	{"en", English},
	{"eng", English},
	{"english", English},
	{"jp", Japanese},
	{"ja", Japanese},
	{"jpn", Japanese},
	{"japanese", Japanese},
	{"cn", SimplifiedChinese},
	{"chinese", SimplifiedChinese},
	{"sc", SimplifiedChinese},
	{"zh", SimplifiedChinese},
	{"zh-cn", SimplifiedChinese},
	{"zh_cn", SimplifiedChinese},
	{"zhcn", SimplifiedChinese},
	{"simplified chinese", SimplifiedChinese},
	{"chinese (simplified)", SimplifiedChinese},
	{"简体中文", SimplifiedChinese},
	{"tc", TraditionalChinese},
	{"tw", TraditionalChinese},
	{"hk", TraditionalChinese},
	{"zh-tw", TraditionalChinese},
	{"zh_tw", TraditionalChinese},
	{"zhtw", TraditionalChinese},
	{"zh-hk", TraditionalChinese},
	{"traditional chinese", TraditionalChinese},
	{"chinese (traditional)", TraditionalChinese},
	{"繁體中文", TraditionalChinese},
	{"繁体中文", TraditionalChinese},
	{"kr", Korean},
	{"ko", Korean},
	{"kor", Korean},
	{"korean", Korean},
	{"한국어", Korean},
	{"es", Spanish},
	{"esp", Spanish},
	{"spa", Spanish},
	{"spanish", Spanish},
	{"español", Spanish},
	{"espanol", Spanish},
	{"fr", French},
	{"fra", French},
	{"fre", French},
	{"french", French},
	{"français", French},
	{"francais", French},
	{"de", German},
	{"ger", German},
	{"deu", German},
	{"german", German},
	{"deutsch", German},
	{"it", Italian},
	{"ita", Italian},
	{"italian", Italian},
	{"italiano", Italian},
	{"pt", Portuguese},
	{"por", Portuguese},
	{"portuguese", Portuguese},
	{"pt-br", Portuguese},
	{"pt_br", Portuguese},
	{"ptbr", Portuguese},
	{"português", Portuguese},
	{"portugues", Portuguese},
	{"ru", Russian},
	{"rus", Russian},
	{"russian", Russian},
	{"русский", Russian},
	{"id", Indonesian},
	{"idn", Indonesian},
	{"ind", Indonesian},
	{"indonesian", Indonesian},
	{"bahasa", Indonesian},
	{"bahasa indonesia", Indonesian},
	{"vi", Vietnamese},
	{"vn", Vietnamese},
	{"vie", Vietnamese},
	{"vietnamese", Vietnamese},
	{"tiếng việt", Vietnamese},
	{"th", Thai},
	{"tha", Thai},
	{"thai", Thai},
	{"ไทย", Thai},
}

// LookupLanguage resolves a whole language tag; Raw is deliberately unrecognized.
func LookupLanguage(content string) Language {
	content = strings.ToLower(content)
	for _, p := range languageTags {
		if content == p.tag {
			return p.language
		}
	}
	return ""
}

// DetectLanguage returns the first recognized bracket or parenthesis language tag.
func DetectLanguage(tokens []Token) Language {
	for _, t := range tokens {
		if bracket(t) {
			if l := LookupLanguage(t.Text); l != "" {
				return l
			}
		}
	}
	return ""
}
