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
	shell.Register("login", doLogin)

	// handle "greet".
	shell.Register("greet", func(c *ishell.Context) {
		name := "Stranger"
		if len(c.Args) > 0 {
			name = strings.Join(c.Args, " ")
		}
		c.Println("Hello", name)
	})

	// read multiple lines with "multi" command
	shell.Register("multi", func(c *ishell.Context) {
		c.Println("Input multiple lines and end with semicolon ';'.")
		lines := c.ReadMultiLines(";")
		c.Println("Done reading. You wrote:")
		c.Println(lines)
	})

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
