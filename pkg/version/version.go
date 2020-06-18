/*
 * Copyright 2017 caicloud authors. All rights reserved.
 */

package version

import (
	"fmt"
	"runtime"
)

// Following values should be substituted with a real value during build.
var (
	// VERSION is the app-global version string.
	VERSION = "UNKNOWN"

	// COMMIT is the app-global git sha string.
	COMMIT = "UNKNOWN"

	// REPOROOT is the app-global repository path string.
	REPOROOT = "UNKNOWN"
)

// PrintVersion prints versions from the array returned by Info().
func PrintVersion() {
	for _, i := range Info() {
		fmt.Printf("%v\n", i)
	}
}

// Info returns an array of various service versions
func Info() []string {
	return []string{
		fmt.Sprintf("Version: %s", VERSION),
		fmt.Sprintf("Git SHA: %s", COMMIT),
		fmt.Sprintf("Repo Root: %s", REPOROOT),
		fmt.Sprintf("Go Version: %s", runtime.Version()),
		fmt.Sprintf("Go OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH),
	}
}
