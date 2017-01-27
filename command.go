package ishell

import (
	"fmt"
	"os"
	"text/tabwriter"
)

// Cmd is a shell command handler.
type Cmd struct {
	Name     string
	Func     func(c *Context)
	Help     string
	children map[string]*Cmd
	parent   *Cmd
}

func addHelpCmd(c *Cmd) {
	c.AddCmd(&Cmd{
		Name: "help",
		Help: "displays help",
		Func: func(*Context) {
			c.PrintHelp()
		},
	})
}

// AddCmd adds cmd as a subcommand.
func (c *Cmd) AddCmd(cmd *Cmd) {
	if c.children == nil {
		c.children = make(map[string]*Cmd)
	}
	c.children[cmd.Name] = cmd
	if _, ok := c.children["help"]; !ok {
		addHelpCmd(cmd)
	}
}

// DeleteCmd deletes cmd from subcommands.
func (c *Cmd) DeleteCmd(name string) {
	delete(c.children, name)
}

// Children returns the subcommands of c.
func (c *Cmd) Children() map[string]*Cmd {
	return c.children
}

// Parent returns the parent command to c.
func (c *Cmd) Parent() *Cmd {
	return c.parent
}

// PrintHelp prints the help of the command.
func (c Cmd) PrintHelp() {
	fmt.Println("Commands:")
	w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
	for _, child := range c.children {
		fmt.Fprintf(w, "\t%s\t\t\t%s\n", child.Name, child.Help)
	}
	w.Flush()
}

func (c Cmd) findFunc(args []string) (Func, []string) {
	var i int
	var arg string
	found := false
	for i, arg = range args {
		if cmd, ok := c.children[arg]; ok {
			c = *cmd
			found = true
			continue
		}
		found = false
		break
	}
	if found {
		if len(args) > i {
			return c.Func, args[i+1:]
		}
		return c.Func, []string{}
	}
	// not found
	if i < 0 {
		// no top level match
		return nil, nil
	}
	if c.Func == nil {
		return nil, nil
	}
	return c.Func, args[i:]
}
