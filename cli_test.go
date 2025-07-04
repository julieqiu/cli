package cli

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParseAndSetFlags(t *testing.T) {
	var (
		strFlag string
		intFlag int
	)

	cmd := &Command{
		Short:     "test is used for testing",
		Long:      "This is the long documentation for command test.",
		UsageLine: "foobar test [arguments]",
		Run: func(ctx context.Context, cmd *Command, args []string) error {
			return cmd.Flags.Parse(args)
		},
	}
	cmd.Flags.StringVar(&strFlag, "name", "default", "name flag")
	cmd.Flags.IntVar(&intFlag, "count", 0, "count flag")

	args := []string{"-name=foo", "-count=5"}
	if err := cmd.Run(t.Context(), cmd, args); err != nil {
		t.Fatalf("Parse() failed: %v", err)
	}

	if strFlag != "foo" {
		t.Errorf("expected name=foo, got %q", strFlag)
	}
	if intFlag != 5 {
		t.Errorf("expected count=5, got %d", intFlag)
	}
}

func TestLookup(t *testing.T) {
	commands := []*Command{
		{Short: "foo runs the foo command"},
		{Short: "bar runs the bar command"},
	}

	for _, test := range []struct {
		name    string
		wantErr bool
	}{
		{"foo", false},
		{"bar", false},
		{"baz", true}, // not found case
	} {
		t.Run(test.name, func(t *testing.T) {
			cmd := &Command{}
			cmd.Commands = commands
			sub, err := cmd.Lookup(test.name)
			if test.wantErr {
				if err == nil {
					t.Fatal(err)
				}
				return
			}

			if err != nil {
				t.Fatal(err)
			}
			if sub.Name() != test.name {
				t.Errorf("got = %q, want = %q", sub.Name(), test.name)
			}
		})
	}
}

func TestRun(t *testing.T) {
	executed := false
	cmd := &Command{
		Short: "run runs the command",
		Run: func(ctx context.Context, cmd *Command, args []string) error {
			executed = true
			return nil
		},
	}

	if err := cmd.Run(t.Context(), cmd, nil); err != nil {
		t.Fatal(err)
	}
	if !executed {
		t.Errorf("cmd.Run was not executed")
	}
}

func TestUsage(t *testing.T) {
	preamble := `Test prints test information.

Usage:
  test [flags]

`

	for _, test := range []struct {
		name  string
		flags []func(c *Command)
		want  string
	}{
		{
			name:  "no flags",
			flags: nil,
			want:  preamble,
		},
		{
			name: "with string flag",
			flags: []func(c *Command){
				func(c *Command) {
					c.Flags.String("name", "default", "name flag")
				},
			},
			want: fmt.Sprintf(`%sFlags:
  -name string
    	name flag (default "default")


`, preamble),
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			c := &Command{
				Short:     "test prints test information",
				UsageLine: "test [flags]",
				Long:      "Test prints test information.",
			}
			for _, fn := range test.flags {
				fn(c)
			}

			var buf bytes.Buffer
			c.usage(&buf)
			got := buf.String()
			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("mismatch(-want + got):\n%s", diff)
			}
		})
	}
}
