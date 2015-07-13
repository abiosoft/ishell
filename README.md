# ishell
ishell is an interactive shell library for creating interactive cli applications.

[![Documentation](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/abiosoft/ishell)

### Usage

```go
import "github.com/abiosoft/ishell"

func main(){
    // create new shell.
    // by default, new shell includes 'exit', 'help' and 'clear' commands.
    shell := ishell.NewShell()

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
>> help
Commands:
exit help greet
>> greet Someone Somewhere
Hello Someone Somewhere
>> exit
$
```

##### Reading input.
```go
// simulate an authentication
shell.Register("login", func(args ...string) (string, error) {
	// disable the '>>' for cleaner same line input.
	shell.ShowPrompt(false)
	defer shell.ShowPrompt(true) // yes, revert after login.

    // get username
	shell.Print("Username: ")
	username := shell.ReadLine()

    // get password. Does not echo characters.
	shell.Print("Password: ")
	password := shell.ReadPassword(false)

	... // do something with username and password

    return "Authentication Successful.", nil
})
```
Execution
```
>> login
Username: someusername
Password:
Authentication Successful.
```

##### How about multiline input.
```go
shell.Register("multi", func(args ...string) (string, error) {
	shell.Println("Input some lines:")
	// read until a semicolon ';' is found
	lines := shell.ReadMultiLines(";")
	shell.Println("You wrote:")
	return lines, nil
})
```
Execution
```
>> multi
Input some lines:
>> this is a sample 
>> of multiline input;
You wrote:
this is a sample
of multiline input;
```

Check example code for more.

### Note
ishell is in active development and can still change significantly.

### Roadmap (in no particular order)
* ~~Support multiline inputs~~.
* Handle ^C interrupts.
* Support coloured outputs.
* Command history.
* Tab completion.
* Testing, testing, testing.

### Contribution
1. Create an issue to discuss it.
2. Send in Pull Request.

### License
MIT

### Credits
Library | Use
------- | -----
[github.com/flynn/go-shlex](http://github.com/flynn/go-shlex) | splitting input into command and args.
[github.com/howeyc/gopass](http://github.com/howeyc/gopass) | reading passwords.
