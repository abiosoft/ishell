package ishell

import (
	"bytes"
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/chzyer/readline"
)

// Actions are actions that can be performed by a shell.
type Actions interface {
	// ReadLine reads a line from standard input.
	ReadLine() string
	// ReadLineErr is ReadLine but returns error as well
	ReadLineErr() (string, error)
	// ReadPassword reads password from standard input without echoing the characters.
	// Note that this only works as expected when the standard input is a terminal.
	ReadPassword() string
	// ReadPasswordErr is ReadPassword but returns error as well
	ReadPasswordErr() (string, error)
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
	// ShowPaged shows a paged text that is scrollable.
	// This leverages on "less" for unix and "more" for windows.
	ShowPaged(text string) error
	// MultiChoice presents options to the user.
	// returns the index of the selection or -1 if nothing is
	// selected.
	// text is displayed before the options.
	MultiChoice(options []string, text string) int
	// Checklist is similar to MultiChoice but user can choose multiple variants using Space.
	// init is initially selected options.
	Checklist(options []string, text string, init []int) []int
	// SetPrompt sets the prompt string. The string to be displayed before the cursor.
	SetPrompt(prompt string)
	// SetMultiPrompt sets the prompt string used for multiple lines. The string to be displayed before
	// the cursor; starting from the second line of input.
	SetMultiPrompt(prompt string)
	// ShowPrompt sets whether prompt should show when requesting input for ReadLine and ReadPassword.
	// Defaults to true.
	ShowPrompt(show bool)
	// Cmds returns all the commands added to the shell.
	Cmds() []*Cmd
	// HelpText returns the computed help of top level commands.
	HelpText() string
	// ClearScreen clears the screen. Same behaviour as running 'clear' in unix terminal or 'cls' in windows cmd.
	ClearScreen() error
	// Stop stops the shell. This will stop the shell from auto reading inputs and calling
	// registered functions. A stopped shell is only inactive but totally functional.
	// Its functions can still be called and can be restarted.
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

func (s *shellActionsImpl) ReadLineErr() (string, error) {
	return s.readLine()
}

func (s *shellActionsImpl) ReadPassword() string {
	return s.reader.readPassword()
}

func (s *shellActionsImpl) ReadPasswordErr() (string, error) {
	return s.reader.readPasswordErr()
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

func (s *shellActionsImpl) MultiChoice(options []string, text string) int {
	choice := s.multiChoice(options, text, nil, false)
	return choice[0]
}
func (s *shellActionsImpl) Checklist(options []string, text string, init []int) []int {
	return s.multiChoice(options, text, init, true)
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

func (s *shellActionsImpl) Cmds() []*Cmd {
	var cmds []*Cmd
	for _, cmd := range s.rootCmd.children {
		cmds = append(cmds, cmd)
	}
	return cmds
}

func (s *shellActionsImpl) ClearScreen() error {
	return clearScreen(s.Shell)
}

func (s *shellActionsImpl) ShowPaged(text string) error {
	return showPaged(s.Shell, text)
}

func (s *shellActionsImpl) Stop() {
	s.stop()
}

func (s *shellActionsImpl) HelpText() string {
	return s.rootCmd.HelpText()
}

func clearScreen(s *Shell) error {
	_, err := readline.ClearScreen(s.writer)
	return err
}

func showPaged(s *Shell, text string) error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("more")
	} else {
		cmd = exec.Command("less")
	}
	cmd.Stdout = s.writer
	cmd.Stderr = s.writer
	cmd.Stdin = bytes.NewBufferString(text)
	return cmd.Run()
}
