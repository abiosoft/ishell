package ishell

import (
	"bytes"
	"fmt"
	"text/tabwriter"
)

// Cmd is a shell command handler.
type Cmd struct {
	// Command name.
	Name string
	// Function to execute for the command.
	Func func(c *Context)
	// One liner help message for the command.
	Help string
	// More descriptive help message for the command.
	LongHelp string

	// subcommands.
	children map[string]*Cmd
}

// AddCmd adds cmd as a subcommand.
func (c *Cmd) AddCmd(cmd *Cmd) {
	if c.children == nil {
		c.children = make(map[string]*Cmd)
	}
	c.children[cmd.Name] = cmd
}

// DeleteCmd deletes cmd from subcommands.
func (c *Cmd) DeleteCmd(name string) {
	delete(c.children, name)
}

// Children returns the subcommands of c.
func (c *Cmd) Children() map[string]*Cmd {
	return c.children
}

func (c *Cmd) hasSubcommand() bool {
	if len(c.children) > 1 {
		return true
	}
	if _, ok := c.children["help"]; !ok {
		return len(c.children) > 0
	}
	return false
}

// HelpText returns the computed help of the command and its subcommands.
func (c Cmd) HelpText() string {
	var b bytes.Buffer
	p := func(s ...interface{}) {
		fmt.Fprintln(&b)
		if len(s) > 0 {
			fmt.Fprintln(&b, s...)
		}
	}
	if c.LongHelp != "" {
		p(c.LongHelp)
	} else if c.Help != "" {
		p(c.Help)
	} else if c.Name != "" {
		p(c.Name, "has no help")
	}
	if c.hasSubcommand() {
		p("Commands:")
		w := tabwriter.NewWriter(&b, 0, 4, 2, ' ', 0)
		for _, child := range c.children {
			fmt.Fprintf(w, "\t%s\t\t\t%s\n", child.Name, child.Help)
		}
		w.Flush()
		p()
	}
	return b.String()
}

// FindCmd finds the matching Cmd for args.
// It returns the Cmd and the remaining args.
func (c Cmd) FindCmd(args []string) (*Cmd, []string) {
	var cmd *Cmd
	for i, arg := range args {
		if cmd1, ok := c.children[arg]; ok {
			cmd = cmd1
			c = *cmd
			continue
		}
		return cmd, args[i:]
	}
	return cmd, nil
}

// GetSubcommand gets the corresponding Cmd for subcommand in args.
func (c Cmd) GetSubcommand(args []string) *Cmd {
	for i, arg := range args {
		if _, ok := c.children[arg]; !ok {
			break
		}
		if i == len(args)-1 {
			return c.children[arg]
		}
	}
	return nil
}
