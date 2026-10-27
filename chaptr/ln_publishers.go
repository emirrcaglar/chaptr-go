// SPDX-License-Identifier: GPL-3.0-or-later
// Go port and modifications of chaptr 1.5.0, 2026-10-27, by emirrcaglar.
// Upstream: https://github.com/johnthreekay/chaptr. See NOTICE for attribution.

package chaptr

// LookupPublisher returns the canonical name, or an empty string if unknown.
func LookupPublisher(name string) string {
	for _, n := range []string{"Yen Press", "J-Novel Club", "Seven Seas", "Seven Seas Siren", "Vertical", "Kodansha", "One Peace Books", "Cross Infinite World", "Tentai Books", "Hanashi Media"} {
		if equalASCII(name, n) {
			return n
		}
	}
	return ""
}
