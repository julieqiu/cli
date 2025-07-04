package cli

import (
	"runtime/debug"
	"strings"
	"time"
)

// Version return the version information for the binary, which is constructed
// following https://go.dev/ref/mod#versions.
func Version() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return ""
	}
	return version(info)
}

func version(info *debug.BuildInfo) string {
	if info.Main.Version != "" && info.Main.Version != "(devel)" {
		return info.Main.Version
	}

	var revision, at string
	for _, s := range info.Settings {
		if s.Key == "vcs.revision" {
			revision = s.Value
		}
		if s.Key == "vcs.time" {
			at = s.Value
		}
	}

	if revision == "" && at == "" {
		return "not available"
	}

	// Construct the pseudo-version string per
	// https://go.dev/ref/mod#pseudo-versions.
	var buf strings.Builder
	buf.WriteString("0.0.0")
	if revision != "" {
		buf.WriteString("-")
		// Per https://go.dev/ref/mod#pseudo-versions, only use the first 12
		// letters of the commit hash.
		buf.WriteString(revision[:12])
	}
	if at != "" {
		// commit time is of the form 2023-01-25T19:57:54Z
		p, err := time.Parse(time.RFC3339, at)
		if err == nil {
			buf.WriteString("-")
			buf.WriteString(p.Format("20060102150405"))
		}
	}
	return buf.String()
}
