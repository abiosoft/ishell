# ishell
ishell is an interactive shell library for creating interactive cli applications.

[![Documentation](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/abiosoft/ishell)
[![Go Report Card](https://goreportcard.com/badge/github.com/abiosoft/ishell)](https://goreportcard.com/report/github.com/abiosoft/ishell)

## Older version
The current master is not backward compatible with older version. Kindly change your import path to `gopkg.in/abiosoft/ishell.v1`.

Older version of this library is still available at [https://gopkg.in/abiosoft/ishell.v1](https://gopkg.in/abiosoft/ishell.v1).

However, you are advised to upgrade to v2 [https://gopkg.in/abiosoft/ishell.v2](https://gopkg.in/abiosoft/ishell.v2).

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
    shell.AddCmd(&ishell.Cmd{
        Name: "greet",
        Help: "greet user",
        Func: func(c *ishell.Context) {
            c.Println("Hello", strings.Join(c.Args, " "))
        },
    })

    // run shell
    shell.Run()
}
```
Execution
```
Sample Interactive Shell
>>> help

Commands:
  clear      clear the screen
  greet      greet user
  exit       exit the program
  help       display help

>>> greet Someone Somewhere
Hello Someone Somewhere
>>> exit
$
```

### Reading input
```go
// simulate an authentication
shell.AddCmd(&ishell.Cmd{
    Name: "login",
    Help: "simulate a login",
    Func: func(c *ishell.Context) {
        // disable the '>>>' for cleaner same line input.
        c.ShowPrompt(false)
        defer c.ShowPrompt(true) // yes, revert after login.

        // get username
        c.Print("Username: ")
        username := c.ReadLine()

        // get password.
        c.Print("Password: ")
        password := c.ReadPassword()

        ... // do something with username and password

        c.Println("Authentication Successful.")
    },
})
```
Execution
```
>>> login
Username: someusername
Password:
Authentication Successful.
```

### Multiline input
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
```
Execution
```
>>> multi
Input multiple lines and end with semicolon ';'.
>>> this is user defined
... multiline input;
You wrote:
this is user defined
multiline input;
```
### Keyboard interrupt
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
shell.Interrupt(func(count int, c *ishell.Context) { ... })
```

### Multiple Choice

```go
func(c *ishell.Context) {
    choice := c.MultiChoice([]string{
        "Golangers",
        "Go programmers",
        "Gophers",
        "Goers",
    }, "What are Go programmers called ?")
    if choice == 2 {
        c.Println("You got it!")
    } else {
        c.Println("Sorry, you're wrong.")
    }
},
```
Output
```
What are Go programmers called ?
  Golangers
  Go programmers
> Gophers
  Goers

You got it!
```
### Checklist
```go
func(c *ishell.Context) {
    languages := []string{"Python", "Go", "Haskell", "Rust"}
    choices := c.Checklist(languages,
        "What are your favourite programming languages ?", nil)
    out := func() []string { ... } // convert index to language
    c.Println("Your choices are", strings.Join(out(), ", "))
}
```
Output
```
What are your favourite programming languages ?
    Python
  ✓ Go
    Haskell
 >✓ Rust

Your choices are Go, Rust
```

### Progress Bar
Determinate
```go
func(c *ishell.Context) {
    c.ProgressBar().Start()
    for i := 0; i < 101; i++ {
        c.ProgressBar().Suffix(fmt.Sprint(" ", i, "%"))
        c.ProgressBar().Progress(i)
        ... // some background computation
    }
    c.ProgressBar().Stop()
}
```
Output
```
[==========>         ] 50%
```

Indeterminate
```go

func(c *ishell.Context) {
    c.ProgressBar().Indeterminate(true)
    c.ProgressBar().Start()
    ... // some background computation
    c.ProgressBar().Stop()
}
```
Output
```
[ ====               ]
```

Custom display using [briandowns/spinner](https://github.com/briandowns/spinner).
```go
display := ishell.ProgressDisplayCharSet(spinner.CharSets[11])
func(c *Context) { c.ProgressBar().Display(display) ... }

// or set it globally
ishell.ProgressBar().Display(display)
```

### Durable history
```go
// Read and write history to $HOME/.ishell_history
shell.SetHomeHistoryPath(".ishell_history")
```


### Non-interactive execution
In some situations it is desired to exit the program directly after executing a single command.

```go
// when started with "exit" as first argument, assume non-interactive execution
if len(os.Args) > 1 && os.Args[1] == "exit" {
    shell.Process(os.Args[2:]...)
} else {
    // start shell
    shell.Run()
}
```

```bash
# Run normally - interactive mode:
$ go run main.go
>>> |

# Run non-interactivelly
$ go run main.go exit greet Someusername
Hello Someusername
```


### Output with Color
You can use [fatih/color](https://github.com/fatih/color).

```go
func(c *ishell.Context) {
    yellow := color.New(color.FgYellow).SprintFunc()
    c.Println(yellow("This line is yellow"))
}
```
Execution
```sh
>>> color
This line is yellow
```


### Example
Available [here](https://github.com/abiosoft/ishell/blob/master/example/main.go).
```sh
go run example/main.go
```

## Supported Platforms
* [x] Linux
* [x] OSX
* [x] Windows [Not tested but should work]

## Note
ishell is in active development and can still change significantly.

## Roadmap (in no particular order)
* [x] Multiline inputs
* [x] Command history
* [x] Customizable tab completion
* [x] Handle ^C interrupts
* [x] Subcommands and help texts
* [x] Scrollable paged output
* [x] Progress bar
* [x] Multiple choice prompt
* [x] Checklist prompt
* [x] Support for command aliases
* [ ] Multiple line progress bars
* [ ] Testing, testing, testing

## Contribution
1. Create an issue to discuss it.
2. Send in Pull Request.

## License
MIT

## Credits
Library | Use
------- | -----
[github.com/flynn-archive/go-shlex](https://github.com/flynn-archive/go-shlex) | splitting input into command and args.
[github.com/chzyer/readline](https://github.com/chzyer/readline) | readline capabilities.


## Donate
```
bitcoin: 1GTHYEDiy2C7RzXn5nY4wVRaEN2GvLjwZN
paypal: a@abiosoft.com
```
