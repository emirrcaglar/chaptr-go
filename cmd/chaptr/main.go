// SPDX-License-Identifier: GPL-3.0-or-later
// Go port and modifications of chaptr 1.5.0, 2026-10-27, by emirrcaglar.
// Upstream: https://github.com/johnthreekay/chaptr. See NOTICE for attribution.

// Command chaptr prints parsed release filename metadata as JSON.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/emirrcaglar/chaptr-go/chaptr"
)

func main() {
	novel := flag.Bool("novel", false, "parse as a light novel instead of manga")
	flag.Usage = func() {
		fmt.Fprintln(flag.CommandLine.Output(), "Usage: chaptr [-novel] \"release filename\"")
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
	}
	var parsed any
	if *novel {
		parsed = chaptr.ParseNovel(flag.Arg(0))
	} else {
		parsed = chaptr.ParseManga(flag.Arg(0))
	}
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(parsed); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
