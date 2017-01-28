package main

import (
	"strings"

	"github.com/abiosoft/ishell"
)

func main() {
	shell := ishell.New()

	// display info
	shell.Println("Sample Interactive Shell")

	// handle login
	shell.AddCmd(&ishell.Cmd{
		Name: "login",
		Func: doLogin,
		Help: "simulate a login",
	})

	// handle "greet".
	shell.AddCmd(&ishell.Cmd{
		Name: "greet",
		Help: "greet user",
		Func: func(c *ishell.Context) {
			name := "Stranger"
			if len(c.Args) > 0 {
				name = strings.Join(c.Args, " ")
			}
			c.Println("Hello", name)
		},
	})

	// read multiple lines with "multi" command
	shell.AddCmd(&ishell.Cmd{
		Name: "multi",
		Help: "input in multiple lines",
		Func: func(c *ishell.Context) {
			c.Println("Input multiple lines and end with semicolon ';'.")
			lines := c.ReadMultiLines(";")
			c.Println("Done reading. You wrote:")
			c.Println(lines)
		},
	})

	cmd := &ishell.Cmd{
		Name: "test",
		Help: "test subcommand",
		LongHelp: `Test subcommands

This test how subcommand works using sub1 and sub2 subcommands.
It also shows how long help is used, if set.`,
		Func: func(c *ishell.Context) {
			c.Println("parent command works if Func is not nil. Has args", c.Args)
		},
	}
	cmd.AddCmd(&ishell.Cmd{
		Name: "sub1",
		Help: "test sub 1",
		Func: func(c *ishell.Context) {
			c.Println("this is sub1 with args", c.Args)
		},
	})
	cmd.AddCmd(&ishell.Cmd{
		Name: "sub2",
		Help: "test sub 2",
		Func: func(c *ishell.Context) {
			c.Println("this is sub2 with args", c.Args)
		},
	})
	shell.AddCmd(cmd)

	// start shell
	shell.Start()
}

func doLogin(c *ishell.Context) {
	c.ShowPrompt(false)
	defer c.ShowPrompt(true)

	c.Println("Let's simulate login")

	// prompt for input
	c.Print("Username: ")
	username := c.ReadLine()
	c.Print("Password: ")
	password := c.ReadPassword()

	// do something with username and password
	c.Println("Your inputs were", username, "and", password+".")

}
