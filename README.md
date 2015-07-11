# ishell
ishell is an interactive shell library for creating interactive cli applications.

[![Documentation](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/abiosoft/ishell)

### Usage

```go
import "github.com/abiosoft/ishell"

func main(){
    // create new shell.
    shell := ishell.NewShell()

	// display welcome info.
	shell.Println("Sample Interactive Shell")

	// register a function for "exit" command.
	shell.Register("exit", func(cmd string, args []string) (string, error) {
		shell.Stop()
		return "bye!", nil
	})

	// start shell
	shell.Start()
}
```
Execution
```
Sample Interactive Shell
>> exit
bye!
```

#### Let's do more.
```go
// simulate an authentication
shell.Register("login", func(cmd string, args []string) (string, error) {
    // get username
	shell.Println("Username:")
	username, _ := shell.ReadLine()

    // get password. Does not echo characters.
	shell.Println("Password:")
	password := shell.ReadPassword(false)

	... // do something with username and password

    return "Authentication Successful.", nil
})
```
Execution
```
Username:
>> someusername
Password:
>>
Authentication Successful.
```
Check example code for more.

### Note
ishell is in active development and can still change significantly.

### Roadmap (in no particular order)
* Support multiline inputs.
* Handle ^C interrupts.
* Support coloured outputs.
* Testing, testing, testing.

### Contribution
1. Create an issue to discuss it.
2. Send in Pull Request.

### License
MIT

### Credits
* github.com/flynn/go-shlex for splitting input into command and args.
* github.com/howeyc/gopass for reading passwords.
