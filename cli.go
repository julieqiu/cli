package cli

import (
	"context"
	"flag"
	"fmt"
	"io"
	"strings"
)

// Command represents a single command that can be executed by the application.
type Command struct {
	// Short is a concise one-line description of the command.
	Short string

	// UsageLine is the one line usage.
	UsageLine string

	// Long is the full description of the command.
	Long string

	// Run executes the command.
	Run func(ctx context.Context, args []string) error

	// Commands are the child commands.
	Commands []*Command

	// Flag is a set of flags specific to this command.
	Flags flag.FlagSet
}

// Parse parses the provided command-line arguments using the command's flag
// set.
func (c *Command) Parse(args []string) error {
	return c.Flags.Parse(args)
}

// Name is the command name. Command.Short is always expected to begin with
// this name.
func (c *Command) Name() string {
	if c.Short == "" {
		panic("command is missing documentation")
	}
	parts := strings.Fields(c.Short)
	return parts[0]
}

// Lookup returns the subcommand with the given name.
func (c *Command) Lookup(name string) (*Command, error) {
	for _, sub := range c.Commands {
		if sub.Name() == name {
			return sub, nil
		}
	}
	return nil, fmt.Errorf("invalid command: %q", name)
}

func (c *Command) usage(w io.Writer) {
	if c.Short == "" || c.UsageLine == "" || c.Long == "" {
		panic(fmt.Sprintf("command %q is missing documentation", c.Name()))
	}

	fmt.Fprintf(w, "%s\n\n", c.Long)
	fmt.Fprintf(w, "Usage:\n  %s", c.UsageLine)
	if len(c.Commands) > 0 {
		fmt.Fprint(w, "\n\nCommands:\n")
		for _, c := range c.Commands {
			parts := strings.Fields(c.Short)
			short := strings.Join(parts[1:], " ")
			fmt.Fprintf(w, "\n  %-25s  %s", c.Name(), short)
		}
	}
	if hasFlags(c.Flags) {
		fmt.Fprint(w, "\n\nFlags:\n")
	}
	c.Flags.SetOutput(w)
	c.Flags.PrintDefaults()
	fmt.Fprintf(w, "\n\n")
}

func (c *Command) Usage() {
	c.Flags.Usage = func() {
		c.usage(c.Flags.Output())
	}
	c.Flags.Usage()
}

func hasFlags(fs *flag.FlagSet) bool {
	visited := false
	fs.VisitAll(func(f *flag.Flag) {
		visited = true
	})
	return visited
}
