// Package ishell implements an interactive shell.
package ishell

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strings"
	"sync"

	"github.com/flynn/go-shlex"
	"github.com/howeyc/gopass"
)

const (
	defaultPrompt = ">> "
)

type Shell struct {
	prompt      string
	showPrompt  bool
	functions   map[string]CmdFunc
	generic     CmdFunc
	reader      *shellReader
	writer      io.Writer
	active      bool
	activeMutex sync.RWMutex
	ignoreCase  bool
	haltChan    chan struct{}
}

// New creates a new shell with default settings. Uses standard output and default prompt ">>".
func New() *Shell {
	shell := &Shell{
		prompt:     defaultPrompt,
		showPrompt: true,
		functions:  make(map[string]CmdFunc),
		reader: &shellReader{
			scanner: bufio.NewReader(os.Stdin),
		},
		writer:   os.Stdout,
		haltChan: make(chan struct{}),
	}
	addDefaultFuncs(shell)
	return shell
}

// Start starts the shell. It reads inputs from standard input and calls registered functions
// accordingly. This function blocks until the shell is stopped.
func (s *Shell) Start() {
	s.start()
}

func (s *Shell) start() {
	if s.Active() {
		return
	}
	s.activeMutex.Lock()
	s.active = true
	s.activeMutex.Unlock()

shell:
	for s.Active() {
		var line string
		var err error
		read := make(chan struct{})
		go func() {
			line, err = s.readLine()
			read <- struct{}{}
		}()
		select {
		case <-read:
			break
		case <-s.haltChan:
			continue shell
		}
		if err == io.EOF {
			fmt.Println("EOF")
			break
		} else if err != nil {
			s.Println("Error:", err)
			break
		}

		if line == "" {
			continue
		}

		line = strings.TrimSpace(line)

		err = handleInput(s, line)
		if err1, ok := err.(shellError); ok && err != nil {
			switch err1.level {
			case LevelWarn:
				s.Println("Warning:", err)
				continue shell
			case LevelStop:
				s.Println(err)
				break shell
			case LevelExit:
				s.Println(err)
				os.Exit(1)
			case LevelPanic:
				panic(err)
			}
		} else if !ok && err != nil {
			s.Println("Error:", err)
		}
	}
}

// Active tells if the shell is active. i.e. Start is previously called.
func (s *Shell) Active() bool {
	s.activeMutex.RLock()
	defer s.activeMutex.RUnlock()
	return s.active
}

func handleInput(s *Shell, line string) error {
	handled, err := s.handleCommand(line)
	if handled || err != nil {
		return err
	}

	// Generic handler
	if s.generic == nil {
		return errNoHandler
	}
	output, err := s.generic(line)
	if err != nil {
		return err
	}
	if output != "" {
		s.Println(output)
	}
	return nil
}

func (s *Shell) handleCommand(line string) (bool, error) {
	str := strings.SplitN(line, " ", 2)
	cmd := str[0]
	if s.ignoreCase {
		cmd = strings.ToLower(cmd)
	}
	var args []string
	if _, ok := s.functions[cmd]; !ok {
		return false, nil
	}
	if len(str) > 1 {
		args1, err := shlex.Split(str[1])
		if err != nil {
			return false, err
		}
		args = args1
	}
	output, err := s.functions[cmd](args...)
	if err != nil {
		return true, err
	}
	if output != "" {
		s.Println(output)
	}
	return true, nil
}

// Stop stops the shell. This will stop the shell from auto reading inputs and calling
// registered functions. A stopped shell is only inactive but totally functional.
// Its functions can still be called.
func (s *Shell) Stop() {
	if !s.Active() {
		return
	}
	s.activeMutex.Lock()
	s.active = false
	s.activeMutex.Unlock()
	go func() {
		s.haltChan <- struct{}{}
	}()
}

// ReadLine reads a line from standard input.
func (s *Shell) ReadLine() string {
	line, _ := s.readLine()
	return line
}

func (s *Shell) readLine() (line string, err error) {
	if s.showPrompt {
		s.Print(s.prompt)
	}
	consumer := make(chan lineString)
	s.reader.ReadLine(consumer)
	ls := <-consumer
	return ls.line, ls.err
}

// ReadMultiLinesFunc reads multiple lines from standard input. It passes each read line to
// f and stops reading when f returns false.
func (s *Shell) ReadMultiLinesFunc(f func(string) bool) string {
	lines := bytes.NewBufferString("")
	for {
		line := s.ReadLine()
		fmt.Fprint(lines, line)
		if !f(line) {
			break
		}
		fmt.Fprintln(lines)
	}
	return lines.String()
}

// ReadMultiLines reads multiple lines from standard input. It stops reading when terminator
// is encountered at the end of the line. It returns the lines read including terminator.
// For more control, use ReadMultiLinesFunc.
func (s *Shell) ReadMultiLines(terminator string) string {
	return s.ReadMultiLinesFunc(func(line string) bool {
		if strings.HasSuffix(strings.TrimSpace(line), terminator) {
			return false
		}
		return true
	})
}

// ReadPassword reads password from standard input without echoing the characters.
// If mask is true, each character will be represented with asterisks '*'. Note that
// this only works as expected when the standard input is a terminal.
func (s *Shell) ReadPassword(mask bool) string {
	if s.showPrompt {
		s.Print(s.prompt)
	}
	if mask {
		return string(gopass.GetPasswdMasked())
	}
	return string(gopass.GetPasswd())
}

// Println prints to output and ends with newline character.
func (s *Shell) Println(val ...interface{}) {
	fmt.Fprintln(s.writer, val...)
}

// Print prints to output.
func (s *Shell) Print(val ...interface{}) {
	fmt.Fprint(s.writer, val...)
}

// Register registers a function for command. It overwrites existing function, if any.
func (s *Shell) Register(command string, function CmdFunc) {
	s.functions[command] = function
}

// Unregister unregisters a function for a command
func (s *Shell) Unregister(command string) {
	delete(s.functions, command)
}

// RegisterGeneric registers a generic function for all inputs.
// It is called if the shell input could not be handled by any of the
// registered functions. Unlike Register, the entire line is passed as
// first argument to CmdFunc.
func (s *Shell) RegisterGeneric(function CmdFunc) {
	s.generic = function
}

// SetPrompt sets the prompt string. The string to be displayed before the cursor.
func (s *Shell) SetPrompt(prompt string) {
	s.prompt = prompt
}

// ShowPrompt sets whether prompt should show when requesting input for ReadLine and ReadPassword.
// Defaults to true.
func (s *Shell) ShowPrompt(show bool) {
	s.showPrompt = show
}

// SetOut sets the writer to write outputs to.
func (s *Shell) SetOut(writer io.Writer) {
	s.writer = writer
}

// PrintCommands prints a space separated list of registered commands to the shell.
func (s *Shell) PrintCommands() {
	out := strings.Join(s.Commands(), " ")
	if out != "" {
		s.Println("Commands:")
		s.Println(out)
	}
}

// Commands returns a sorted list of all registered commands.
func (s *Shell) Commands() []string {
	var commands []string
	for command := range s.functions {
		commands = append(commands, command)
	}
	sort.Strings(commands)
	return commands
}

// IgnoreCase specifies whether commands should not be case sensitive.
// Defaults to false i.e. commands are case sensitive.
// If true, commands must be registered in lower cases. e.g. shell.Register("cmd", ...)
func (s *Shell) IgnoreCase(ignore bool) {
	s.ignoreCase = ignore
}

// ClearScreen clears the screen. Same behaviour as running 'clear' in unix terminal or 'cls' in windows cmd.
func (s *Shell) ClearScreen() error {
	return clearScreen(s)
}

func clearScreen(s *Shell) error {
	c := "clear"
	if runtime.GOOS == "windows" {
		c = "cls"
	}
	cmd := exec.Command(c)
	cmd.Stdout = s.writer
	return cmd.Run()
}
