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

	// Action is the action to execute when the command is run.
	Action func(ctx context.Context, cmd *Command) error

	// Commands are the child commands.
	Commands []*Command

	// Flag is a set of flags specific to this command.
	Flags flag.FlagSet
}

// Run parses the arguments and executes Command.Action.
func (c *Command) Run(ctx context.Context, args []string) error {
	if err := c.Flags.Parse(args); err != nil {
		return err
	}
	c.Flags.Usage = func() {
		c.usage(c.Flags.Output())
	}
	if c.Action != nil {
		return c.Action(ctx, c)
	}
	if len(args) == 0 {
		c.Flags.Usage()
		return fmt.Errorf("no arguments provided")
	}
	sub, err := c.lookup(args[0])
	if err != nil {
		return err
	}
	return sub.Run(ctx, args[1:])
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

// lookup returns the subcommand with the given name.
func (c *Command) lookup(name string) (*Command, error) {
	for _, sub := range c.Commands {
		if sub.Name() == name {
			return sub, nil
		}
	}
	return nil, fmt.Errorf("invalid command: %q", name)
}

func hasFlags(fs flag.FlagSet) bool {
	visited := false
	fs.VisitAll(func(f *flag.Flag) {
		visited = true
	})
	return visited
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
