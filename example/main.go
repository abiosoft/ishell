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
	shell.Register("login", func(args ...string) (string, error) {
		doLogin(shell)
		return "", nil
	})

	// handle "greet".
	shell.Register("greet", func(args ...string) (string, error) {
		name := "Stranger"
		if len(args) > 0 {
			name = strings.Join(args, " ")
		}
		return "Hello " + name, nil
	})

	// read multiple lines with "multi" command
	shell.Register("multi", func(args ...string) (string, error) {
		shell.Println("Input multiple lines and end with semicolon ';'.")
		lines := shell.ReadMultiLines(";")
		shell.Println("Done reading. You wrote:")
		return lines, nil
	})

	// start shell
	shell.Start()
}

func doLogin(shell *ishell.Shell) {
	shell.ShowPrompt(false)
	defer shell.ShowPrompt(true)

	shell.Println("Let's simulate login")

	// prompt for input
	shell.Print("Username: ")
	username := shell.ReadLine()
	shell.Print("Password: ")
	password := shell.ReadPassword()

	// do something with username and password
	shell.Println("Your inputs were", username, "and", password+".")

}
