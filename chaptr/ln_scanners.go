// SPDX-License-Identifier: GPL-3.0-or-later
// Go port and modifications of chaptr 1.5.0, 2026-10-27, by emirrcaglar.
// Upstream: https://github.com/johnthreekay/chaptr. See NOTICE for attribution.

package chaptr

// LookupScanner returns the canonical name, or an empty string if unknown.
func LookupScanner(name string) string {
	for _, n := range []string{"LuCaZ", "Stick", "Ushi", "Oak", "nao", "CleanBookGuy", "Kobo", "faratnis", "Antithetical", "DigitalMangaFan"} {
		if equalASCII(name, n) {
			return n
		}
	}
	return ""
}
