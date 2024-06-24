package internal

import (
	"runtime/debug"
)

// GetVersion extracts the commit information from the build info, falling back
// to the string "unknown". If any non-empty "override" variables are passed in,
// they will override the returned version. This still allows for stamping
// releases, while also providing a fallback.
func GetVersion(overrides ...string) string {
	for _, override := range overrides {
		if override != "" {
			return override
		}
	}

	var rev string
	var dirty string

	// If there's embedded build info we use that for the embedded Version.
	if buildInfo, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range buildInfo.Settings {
			switch setting.Key {
			case "vcs.revision":
				rev = setting.Value
			case "vcs.modified":
				dirty = "-dirty"
			}
		}

		if rev != "" {
			return rev + dirty
		}

		if buildInfo.Main.Version == "(devel)" {
			return "dev"
		}
	}

	return "unknown"
}
