# ishell
ishell is an interactive shell library for creating interactive cli applications.

[![Documentation](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/abiosoft/ishell)

## Usage

```go
import "strings"
import "github.com/abiosoft/ishell"

func main(){
    // create new shell.
    // by default, new shell includes 'exit', 'help' and 'clear' commands.
    shell := ishell.New()

	// display welcome info.
	shell.Println("Sample Interactive Shell")

	// register a function for "greet" command.
    shell.Register("greet", func(args ...string) (string, error) {
        name := "Stranger"
        if len(args) > 0 {
            name = strings.Join(args, " ")
        }
		return "Hello "+name, nil
	})

	// start shell
	shell.Start()
}
```
Execution
```
Sample Interactive Shell
>>> help
Commands:
exit help greet
>>> greet Someone Somewhere
Hello Someone Somewhere
>>> exit
$
```

### Reading input.
```go
// simulate an authentication
shell.Register("login", func(args ...string) (string, error) {
	// disable the '>>>' for cleaner same line input.
	shell.ShowPrompt(false)
	defer shell.ShowPrompt(true) // yes, revert after login.

    // get username
	shell.Print("Username: ")
	username := shell.ReadLine()

    // get password.
	shell.Print("Password: ")
	password := shell.ReadPassword()

	... // do something with username and password

    return "Authentication Successful.", nil
})
```
Execution
```
>>> login
Username: someusername
Password:
Authentication Successful.
```

### Multiline input.
Builtin support for multiple lines.
```
>>> This is \
... multi line

>>> Cool that << EOF
... everything here goes
... as a single argument. 
... EOF
```
User defined
```go
shell.Register("multi", func(args ...string) (string, error) {
	shell.Println("Input some lines:")
	// read until a semicolon ';' is found
	// use shell.ReadMultiLinesFunc for more control.
	lines := shell.ReadMultiLines(";")
	shell.Println("You wrote:")
	return lines, nil
})
```
Execution
```
>>> multi
Input some lines:
>>> this is user defined 
... multiline input;
You wrote:
this is user defined
multiline input;
```
### Keyboard interrupt.
Builtin interrupt handler.
```
>>> ^C
Input Ctrl-C once more to exit
>>> ^C
Interrupted
exit status 1
```
Custom
```go
shell.RegisterInterrupt(func(args ...string) (string, error) { ... })
```

### Durable history.
```go
// Read and write history to $HOME/.ishell_history
shell.SetHomeHistoryPath(".ishell_history")
```

Check example code for more.

## Supported Platforms
* [x] Linux
* [x] OSX
* [x] Windows

## Note
ishell is in active development and can still change significantly.

## Roadmap (in no particular order)
* [x] Support multiline inputs.
* [x] Command history.
* [x] Tab completion.
* [x] Handle ^C interrupts.
* [ ] Subcommands and help texts.
* [ ] Coloured outputs.
* [ ] Testing, testing, testing.

## Contribution
1. Create an issue to discuss it.
2. Send in Pull Request.

## License
MIT

## Credits
Library | Use
------- | -----
[github.com/flynn/go-shlex](http://github.com/flynn/go-shlex) | splitting input into command and args.
[gopkg.in/readline.v1](http://gopkg.in/readline.v1) | history, tab completion and reading passwords.
