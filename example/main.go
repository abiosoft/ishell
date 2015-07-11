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
	shell.Register("exit", func(cmd string, args []string) (string, error) {
		shell.Println("Do you want to do more ? y/n:")
		line, _ := shell.ReadLine()
		if strings.ToLower(line) == "y" {
			doLogin(shell)
		}
		shell.Stop()
		return "bye!", nil
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
	shell.Println("Let's simulate login")

	// prompt for input
	shell.Println("Username:")
	username, _ := shell.ReadLine()
	shell.Println("Password:")
	password := shell.ReadPassword(false)

	// do something with username and password
	shell.Println("Your inputs were", username, "and", password+".")

}
