package ishell

import (
	"os"
)

func exitFunc(c *Context) {
	c.Stop()
}

func helpFunc(c *Context) {
	c.Println(c.HelpText())
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

func interruptFunc(s *Shell) func(int, *Context) {
	return func(count int, c *Context) {
		if count >= 2 {
			c.Println("Interrupted")
			os.Exit(1)
		}
		c.Println("Input Ctrl-c once more to exit")
	}
}
