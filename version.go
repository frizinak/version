package version

import (
	"fmt"
	"runtime/debug"
	"time"
)

// set with e.g.: -ldflags "-X github.com/frizinak/version.version=$(git describe --always)"
var version string

const VersionUnknown = "unknown"

func Get() string {
	if version != "" {
		return version
	}

	version = VersionUnknown
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return version
	}

	var commit string
	var stamp time.Time
	var modified bool
	for _, v := range bi.Settings {
		switch v.Key {
		case "vcs.revision":
			commit = v.Value
		case "vcs.time":
			s, err := time.Parse(time.RFC3339, v.Value)
			if err == nil {
				stamp = s
			}
		case "vcs.modified":
			modified = v.Value != "false"
		}
	}

	if commit == "" && stamp == (time.Time{}) {
		return version
	}

	if len(commit) == 40 {
		commit = commit[:7]
	} else if commit == "" {
		commit = "??????"
	}

	mod := ""
	if modified {
		mod = "*"
	}

	version = fmt.Sprintf(
		"%s%s-%s",
		mod,
		commit,
		stamp.UTC().Format("2006.02.01.15.03.04"),
	)

	return version
}
