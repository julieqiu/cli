package cli

import (
	"runtime/debug"
	"testing"
)

func TestVersion(t *testing.T) {
	for _, test := range []struct {
		name      string
		want      string
		buildinfo *debug.BuildInfo
	}{
		{
			name: "tagged version",
			want: "1.2.3",
			buildinfo: &debug.BuildInfo{
				Main: debug.Module{
					Version: "1.2.3",
				},
			},
		},
		{
			name: "pseudoversion",
			want: "0.0.0-123456789000-20230125195754",
			buildinfo: &debug.BuildInfo{
				Settings: []debug.BuildSetting{
					{Key: "vcs.revision", Value: "1234567890001234"},
					{Key: "vcs.time", Value: "2023-01-25T19:57:54Z"},
				},
			},
		},
		{
			name:      "local development",
			want:      "not available",
			buildinfo: &debug.BuildInfo{},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			if got := version(test.buildinfo); got != test.want {
				t.Errorf("got %s; want %s", got, test.want)
			}
		})
	}
}
