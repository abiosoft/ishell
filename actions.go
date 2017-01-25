package ishell

import (
	"fmt"
	"os/exec"
	"runtime"
	"sort"
	"strings"
)

// Actions are actions that can be performed by a shell.
type Actions interface {
	// ReadLine reads a line from standard input.
	ReadLine() string
	// ReadPassword reads password from standard input without echoing the characters.
	// Note that this only works as expected when the standard input is a terminal.
	ReadPassword() string
	// ReadMultiLinesFunc reads multiple lines from standard input. It passes each read line to
	// f and stops reading when f returns false.
	ReadMultiLinesFunc(f func(string) bool) string
	// ReadMultiLines reads multiple lines from standard input. It stops reading when terminator
	// is encountered at the end of the line. It returns the lines read including terminator.
	// For more control, use ReadMultiLinesFunc.
	ReadMultiLines(terminator string) string
	// Println prints to output and ends with newline character.
	Println(val ...interface{})
	// Print prints to output.
	Print(val ...interface{})
	// Printf prints to output using string format.
	Printf(format string, val ...interface{})
	// SetPrompt sets the prompt string. The string to be displayed before the cursor.
	SetPrompt(prompt string)
	// SetMultiPrompt sets the prompt string used for multiple lines. The string to be displayed before
	// the cursor; starting from the second line of input.
	SetMultiPrompt(prompt string)
	// ShowPrompt sets whether prompt should show when requesting input for ReadLine and ReadPassword.
	// Defaults to true.
	ShowPrompt(show bool)
	// PrintCommands prints a space separated list of registered commands to the shell.
	PrintCommands()
	// Commands returns a sorted list of all registered commands.
	Commands() []string
	// ClearScreen clears the screen. Same behaviour as running 'clear' in unix terminal or 'cls' in windows cmd.
	ClearScreen() error
	// Stop stops the shell. This will stop the shell from auto reading inputs and calling
	// registered functions. A stopped shell is only inactive but totally functional.
	// Its functions can still be called.
	Stop()
}

type shellActionsImpl struct {
	*Shell
}

// ReadLine reads a line from standard input.
func (s *shellActionsImpl) ReadLine() string {
	line, _ := s.readLine()
	return line
}

func (s *shellActionsImpl) ReadPassword() string {
	return s.reader.readPassword()
}

func (s *shellActionsImpl) ReadMultiLinesFunc(f func(string) bool) string {
	lines, _ := s.readMultiLinesFunc(f)
	return lines
}

func (s *shellActionsImpl) ReadMultiLines(terminator string) string {
	return s.ReadMultiLinesFunc(func(line string) bool {
		if strings.HasSuffix(strings.TrimSpace(line), terminator) {
			return false
		}
		return true
	})
}

func (s *shellActionsImpl) Println(val ...interface{}) {
	s.reader.buf.Truncate(0)
	fmt.Fprintln(s.writer, val...)
}

func (s *shellActionsImpl) Print(val ...interface{}) {
	s.reader.buf.Truncate(0)
	fmt.Fprint(s.reader.buf, val...)
	fmt.Fprint(s.writer, val...)
}

func (s *shellActionsImpl) Printf(format string, val ...interface{}) {
	s.reader.buf.Truncate(0)
	fmt.Fprintf(s.reader.buf, format, val...)
	fmt.Fprintf(s.writer, format, val...)
}

func (s *shellActionsImpl) SetPrompt(prompt string) {
	s.reader.prompt = prompt
	s.reader.scanner.SetPrompt(s.reader.rlPrompt())
}

func (s *shellActionsImpl) SetMultiPrompt(prompt string) {
	s.reader.multiPrompt = prompt
}

func (s *shellActionsImpl) ShowPrompt(show bool) {
	s.reader.showPrompt = show
	s.reader.scanner.SetPrompt(s.reader.rlPrompt())
}

func (s *shellActionsImpl) PrintCommands() {
	out := strings.Join(s.Commands(), " ")
	if out != "" {
		s.Println("Commands:")
		s.Println(out)
	}
}

func (s *shellActionsImpl) Commands() []string {
	var commands []string
	for command := range s.functions {
		commands = append(commands, command)
	}
	sort.Strings(commands)
	return commands
}

func (s *shellActionsImpl) ClearScreen() error {
	return clearScreen(s.Shell)
}

func (s *shellActionsImpl) Stop() {
	s.reader.scanner.Close()
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

func clearScreen(s *Shell) error {
	cmd := exec.Command("clear")
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", "cls")
	}
	cmd.Stdout = s.writer
	return cmd.Run()
}
