package main

import (
	"fmt"
	"github.com/alecthomas/kong"
	"io"
)

// COMMIT is set when build
var COMMIT = ""
var GOVERSION = ""

type version bool

// BeforeApply impl a kong hook
func (v version) BeforeApply(ctx *kong.Context) error {
	if GOVERSION == "" {
		GOVERSION = "Unknown"
	}
	if COMMIT == "" {
		COMMIT = "Unknown"
	}
	io.WriteString(ctx.Stdout, fmt.Sprintf("COMMIT: %v\nGO VERSION: %v\n", COMMIT, GOVERSION))
	ctx.Exit(0)
	return nil
}
