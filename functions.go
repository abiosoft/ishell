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
	c.PrintHelp()
}

func clearFunc(c *Context) {
	err := c.ClearScreen()
	if err != nil {
		c.Err(err)
	}
}

func addDefaultFuncs(s *Shell) {
	s.AddCmd(&Cmd{
		Name: "exit",
		Help: "exit the program",
		Func: exitFunc,
	})
	s.AddCmd(&Cmd{
		Name: "help",
		Help: "display help",
		Func: helpFunc,
	})
	s.AddCmd(&Cmd{
		Name: "clear",
		Help: "clear the screen",
		Func: clearFunc,
	})
	s.Interrupt(interruptFunc(s))
}

func interruptFunc(s *Shell) Func {
	return func(c *Context) {
		s.interruptCount++
		if s.interruptCount >= 2 {
			c.Println("Interrupted")
			os.Exit(1)
		}
		c.Println("Input Ctrl-c once more to exit")
	}
}
