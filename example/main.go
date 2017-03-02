package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/abiosoft/ishell"
	"github.com/briandowns/spinner"
)

func main() {
	shell := ishell.New()

	// display info.
	shell.Println("Sample Interactive Shell")

	// handle login.
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
			s := spinner.New(spinner.CharSets[0], 100*time.Millisecond) // Build our new spinner
			s.Suffix = "Some suffix"
			s.Start() // Start the spinner

			time.Sleep(4 * time.Second) // Run for some time to simulate work
			s.Stop()
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
			// clear := func() { c.Print("\033[2K") }
			clear := func(n int) {
				for i := 0; i < n; i++ {
					c.Print("\b")
				}
			}
			for i := 0; i < 100; i++ {
				c.Print(i+1, "%")
				time.Sleep(time.Millisecond * 100)
				clear(len(strconv.Itoa(i+1)) + 1)
			}
			c.Println()
		},
	})

	// subcommands and custom autocomplete.
	{
		var words []string
		autoCmd := &ishell.Cmd{
			Name: "suggest",
			Help: "try auto complete",
			LongHelp: `Try dynamic autocomplete by adding and removing words.
Then view the autocomplete by tabbing after "words" subcommand.

This is an example of a long help.`,
		}
		autoCmd.AddCmd(&ishell.Cmd{
			Name: "add",
			Help: "add words to autocomplete",
			Func: func(c *ishell.Context) {
				if len(c.Args) == 0 {
					c.Err(errors.New("missing word(s)"))
					return
				}
				words = append(words, c.Args...)
			},
		})

		autoCmd.AddCmd(&ishell.Cmd{
			Name: "clear",
			Help: "clear words in autocomplete",
			Func: func(c *ishell.Context) {
				words = nil
			},
		})

		autoCmd.AddCmd(&ishell.Cmd{
			Name: "words",
			Help: "add words with 'suggest add', then tab after typing 'suggest words '",
			Completer: func([]string) []string {
				return words
			},
		})

		shell.AddCmd(autoCmd)
	}

	shell.AddCmd(&ishell.Cmd{
		Name: "paged",
		Help: "show paged text",
		Func: func(c *ishell.Context) {
			lines := ""
			line := `%d. This is a paged text input.
This is another line of it.

`
			for i := 0; i < 100; i++ {
				lines += fmt.Sprintf(line, i+1)
			}
			c.ShowPaged(lines)
		},
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
