package ishell

import (
  "github.com/abiosoft/readline"
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
  s.Interrupt(interruptFunc)
  s.FilterInput(filterInput)
}

func filterInput(r rune) (rune, bool){
  switch r {
  // block CtrlZ feature
  case readline.CharCtrlZ:
  return r, false
  }
  return r, true
}

func interruptFunc(c *Context, count int, line string) {
	if count >= 2 {
		c.Println("Interrupted")
		os.Exit(1)
	}
	c.Println("Input Ctrl-c once more to exit")
}
