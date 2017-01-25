package ishell

import (
	"os"
)

// Func represents a command function that is called after an input to the shell.
type Func func(c *Context)

func exitFunc(c *Context) {
	c.Stop()
}

func helpFunc(c *Context) {
	c.PrintCommands()
}

func clearFunc(c *Context) {
	err := c.ClearScreen()
	if err != nil {
		c.Err(err)
	}
}

func addDefaultFuncs(s *Shell) {
	s.Register("exit", exitFunc)
	s.Register("help", helpFunc)
	s.Register("clear", clearFunc)
	s.RegisterInterrupt(interruptFunc(s))
}

func interruptFunc(s *Shell) Func {
	return func(c *Context) {
		s.interruptCount++
		if s.interruptCount >= 2 {
			c.Println("Interrupted")
			os.Exit(1)
		}
		c.Println("Input Ctrl-C once more to exit")
	}
}
