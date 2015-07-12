package main

import (
	"strings"

	"github.com/abiosoft/ishell"
)

func main() {
	shell := ishell.NewShell()

	// display info
	shell.Println("Sample Interactive Shell")

	// handle exit
	shell.Register("login", func(cmd string, args []string) (string, error) {
		doLogin(shell)
		return "", nil
	})

	// register a function for "greet" command.
	shell.Register("greet", func(cmd string, args []string) (string, error) {
		name := "Stranger"
		if len(args) > 0 {
			name = strings.Join(args, " ")
		}
		return "Hello " + name, nil
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
	username, _ := shell.ReadLine()
	shell.Print("Password: ")
	password := shell.ReadPassword(false)

	// do something with username and password
	shell.Println("Your inputs were", username, "and", password+".")

}
