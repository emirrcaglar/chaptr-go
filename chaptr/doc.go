// SPDX-License-Identifier: GPL-3.0-or-later
// Go port and modifications of chaptr 1.5.0, 2026-10-27, by emirrcaglar.
// Upstream: https://github.com/johnthreekay/chaptr. See NOTICE for attribution.

// Package chaptr parses manga, manhwa, manhua, and light novel release filenames.
// Parsing is pure and forgiving: unrecognized fields are left absent.
//
// ParseManga and ParseNovel are the main entry points. Tokenize and the Detect
// functions expose the shared lexical and detection layers for custom parsers.
package chaptr
