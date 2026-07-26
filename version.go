package main

import (
	"runtime/debug"
	"strings"
)

// Version is the build version. It defaults to the module version embedded by
// the Go toolchain and can be overridden at build time via -ldflags.
var Version = "devel"

// Commit is the VCS revision, optionally set via -ldflags at release time.
var Commit = ""

func init() {
	if info, ok := debug.ReadBuildInfo(); ok {
		if v := info.Main.Version; v != "" && v != "(devel)" && Version == "devel" {
			Version = strings.TrimPrefix(v, "v")
		}
		for _, s := range info.Settings {
			if s.Key == "vcs.revision" && Commit == "" {
				Commit = s.Value
			}
		}
	}
	Version = strings.TrimPrefix(Version, "v")
}
