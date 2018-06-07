// Package ishell implements an interactive shell.
package ishell

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/abiosoft/readline"
	"github.com/fatih/color"
	"github.com/flynn-archive/go-shlex"
)

const (
	defaultPrompt      = ">>> "
	defaultMultiPrompt = "... "
)

var (
	errNoHandler          = errors.New("incorrect input, try 'help'")
	errNoInterruptHandler = errors.New("no interrupt handler")
)

// Shell is an interactive cli shell.
type Shell struct {
	rootCmd           *Cmd
	generic           func(*Context)
	interrupt         func(*Context, int, string)
	interruptCount    int
	eof               func(*Context)
	reader            *shellReader
	writer            io.Writer
	active            bool
	activeMutex       sync.RWMutex
	ignoreCase        bool
	customCompleter   bool
	multiChoiceActive bool
	haltChan          chan struct{}
	historyFile       string
	autoHelp          bool
	rawArgs           []string
	progressBar       ProgressBar
	pager             string
	pagerArgs         []string
	contextValues
	Actions
}

// New creates a new shell with default settings. Uses standard output and default prompt ">> ".
func New() *Shell {
	return NewWithConfig(&readline.Config{Prompt: defaultPrompt})
}

// NewWithConfig creates a new shell with custom readline config.
func NewWithConfig(conf *readline.Config) *Shell {
	rl, err := readline.NewEx(conf)
	if err != nil {
		log.Println("Shell or operating system not supported.")
		log.Fatal(err)
	}
	shell := &Shell{
		rootCmd: &Cmd{},
		reader: &shellReader{
			scanner:     rl,
			prompt:      rl.Config.Prompt,
			multiPrompt: defaultMultiPrompt,
			showPrompt:  true,
			buf:         &bytes.Buffer{},
			completer:   readline.NewPrefixCompleter(),
		},
		writer:   conf.Stdout,
		autoHelp: true,
	}
	shell.Actions = &shellActionsImpl{Shell: shell}
	shell.progressBar = newProgressBar(shell)
	addDefaultFuncs(shell)
	return shell
}

// Start starts the shell but does not wait for it to stop.
func (s *Shell) Start() {
	s.prepareRun()
	go s.run()
}

// Run starts the shell and waits for it to stop.
func (s *Shell) Run() {
	s.prepareRun()
	s.run()
}

// Wait waits for the shell to stop.
func (s *Shell) Wait() {
	<-s.haltChan
}

func (s *Shell) stop() {
	if !s.Active() {
		return
	}
	s.activeMutex.Lock()
	s.active = false
	s.activeMutex.Unlock()
	close(s.haltChan)
}

// Close stops the shell (if required) and closes the shell's input.
// This should be called when done with reading inputs.
// Unlike `Stop`, a closed shell cannot be restarted.
func (s *Shell) Close() {
	s.stop()
	s.reader.scanner.Close()
}

func (s *Shell) prepareRun() {
	if s.Active() {
		return
	}
	if !s.customCompleter {
		s.initCompleters()
	}
	s.activeMutex.Lock()
	s.active = true
	s.activeMutex.Unlock()

	s.haltChan = make(chan struct{})
}

func (s *Shell) run() {
shell:
	for s.Active() {
		var line []string
		var err error
		read := make(chan struct{})
		go func() {
			line, err = s.read()
			read <- struct{}{}
		}()
		select {
		case <-read:
			break
		case <-s.haltChan:
			continue shell
		}

		if err == io.EOF {
			if s.eof == nil {
				fmt.Println("EOF")
				break
			}
			if err := handleEOF(s); err != nil {
				s.Println("Error:", err)
				continue
			}
		} else if err != nil && err != readline.ErrInterrupt {
			s.Println("Error:", err)
			continue
		}

		if err == readline.ErrInterrupt {
			// interrupt received
			err = handleInterrupt(s, line)
		} else {
			// reset interrupt counter
			s.interruptCount = 0

			// normal flow
			if len(line) == 0 {
				// no input line
				continue
			}

			err = handleInput(s, line)
		}
		if err != nil {
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

// Process runs shell using args in a non-interactive mode.
func (s *Shell) Process(args ...string) error {
	return handleInput(s, args)
}

func handleInput(s *Shell, line []string) error {
	handled, err := s.handleCommand(line)
	if handled || err != nil {
		return err
	}

	// Generic handler
	if s.generic == nil {
		return errNoHandler
	}
	c := newContext(s, nil, line)
	s.generic(c)
	return c.err
}

func handleInterrupt(s *Shell, line []string) error {
	if s.interrupt == nil {
		return errNoInterruptHandler
	}
	c := newContext(s, nil, line)
	s.interruptCount++
	s.interrupt(c, s.interruptCount, strings.Join(line, " "))
	return c.err
}

func handleEOF(s *Shell) error {
	c := newContext(s, nil, nil)
	s.eof(c)
	return c.err
}

func (s *Shell) handleCommand(str []string) (bool, error) {
	if s.ignoreCase {
		for i := range str {
			str[i] = strings.ToLower(str[i])
		}
	}
	cmd, args := s.rootCmd.FindCmd(str)
	if cmd == nil {
		return false, nil
	}
	// trigger help if func is not registered or auto help is true
	if cmd.Func == nil || (s.autoHelp && len(args) == 1 && args[0] == "help") {
		s.Println(cmd.HelpText())
		return true, nil
	}
	c := newContext(s, cmd, args)
	cmd.Func(c)
	return true, c.err
}

func (s *Shell) readLine() (line string, err error) {
	consumer := make(chan lineString)
	defer close(consumer)
	go s.reader.readLine(consumer)
	ls := <-consumer
	return ls.line, ls.err
}

func (s *Shell) read() ([]string, error) {
	s.rawArgs = nil
	heredoc := false
	eof := ""
	// heredoc multiline
	lines, err := s.readMultiLinesFunc(func(line string) bool {
		if !heredoc {
			if strings.Contains(line, "<<") {
				s := strings.SplitN(line, "<<", 2)
				if eof = strings.TrimSpace(s[1]); eof != "" {
					heredoc = true
					return true
				}
			}
		} else {
			return line != eof
		}
		return strings.HasSuffix(strings.TrimSpace(line), "\\")
	})

	s.rawArgs = strings.Fields(lines)

	if heredoc {
		s := strings.SplitN(lines, "<<", 2)
		args, err1 := shlex.Split(s[0])

		arg := strings.TrimSuffix(strings.SplitN(s[1], "\n", 2)[1], eof)
		args = append(args, arg)
		if err1 != nil {
			return args, err1
		}
		return args, err
	}

	lines = strings.Replace(lines, "\\\n", " \n", -1)

	args, err1 := shlex.Split(lines)
	if err1 != nil {
		return args, err1
	}

	return args, err
}

func (s *Shell) readMultiLinesFunc(f func(string) bool) (string, error) {
	var lines bytes.Buffer
	currentLine := 0
	var err error
	for {
		if currentLine == 1 {
			// from second line, enable next line prompt.
			s.reader.setMultiMode(true)
		}
		var line string
		line, err = s.readLine()
		fmt.Fprint(&lines, line)
		if !f(line) || err != nil {
			break
		}
		fmt.Fprintln(&lines)
		currentLine++
	}
	if currentLine > 0 {
		// if more than one line is read
		// revert to standard prompt.
		s.reader.setMultiMode(false)
	}
	return lines.String(), err
}

func (s *Shell) initCompleters() {
	s.setCompleter(iCompleter{cmd: s.rootCmd, disabled: func() bool { return s.multiChoiceActive }})
}

func (s *Shell) setCompleter(completer readline.AutoCompleter) {
	config := s.reader.scanner.Config.Clone()
	config.AutoComplete = completer
	s.reader.scanner.SetConfig(config)
}

// CustomCompleter allows use of custom implementation of readline.Autocompleter.
func (s *Shell) CustomCompleter(completer readline.AutoCompleter) {
	s.customCompleter = true
	s.setCompleter(completer)
}

// AddCmd adds a new command handler.
// This only adds top level commands.
func (s *Shell) AddCmd(cmd *Cmd) {
	s.rootCmd.AddCmd(cmd)
}

// DeleteCmd deletes a top level command.
func (s *Shell) DeleteCmd(name string) {
	s.rootCmd.DeleteCmd(name)
}

// NotFound adds a generic function for all inputs.
// It is called if the shell input could not be handled by any of the
// added commands.
func (s *Shell) NotFound(f func(*Context)) {
	s.generic = f
}

// AutoHelp sets if ishell should trigger help message if
// a command's arg is "help". Defaults to true.
//
// This can be set to false for more control on how help is
// displayed.
func (s *Shell) AutoHelp(enable bool) {
	s.autoHelp = enable
}

// Interrupt adds a function to handle keyboard interrupt (Ctrl-c).
// count is the number of consecutive times that Ctrl-c has been pressed.
// i.e. any input apart from Ctrl-c resets count to 0.
func (s *Shell) Interrupt(f func(c *Context, count int, input string)) {
	s.interrupt = f
}

// EOF adds a function to handle End of File input (Ctrl-d).
// This overrides the default behaviour which terminates the shell.
func (s *Shell) EOF(f func(c *Context)) {
	s.eof = f
}

// SetHistoryPath sets where readlines history file location. Use an empty
// string to disable history file. It is empty by default.
func (s *Shell) SetHistoryPath(path string) {
	// Using scanner.SetHistoryPath doesn't initialize things properly and
	// history file is never written. Simpler to just create a new readline
	// Instance.
	config := s.reader.scanner.Config.Clone()
	config.HistoryFile = path
	s.reader.scanner, _ = readline.NewEx(config)
}

// SetHomeHistoryPath is a convenience method that sets the history path
// in user's home directory.
func (s *Shell) SetHomeHistoryPath(path string) {
	home := os.Getenv("HOME")
	if runtime.GOOS == "windows" {
		home = os.Getenv("USERPROFILE")
	}
	abspath := filepath.Join(home, path)
	s.SetHistoryPath(abspath)
}

// SetOut sets the writer to write outputs to.
func (s *Shell) SetOut(writer io.Writer) {
	s.writer = writer
}

// SetPager sets the pager and its arguments for paged output
func (s *Shell) SetPager(pager string, args []string) {
	s.pager = pager
	s.pagerArgs = args
}

func initSelected(init []int, max int) []int {
	selectedMap := make(map[int]bool)
	for _, i := range init {
		if i < max {
			selectedMap[i] = true
		}
	}
	selected := make([]int, len(selectedMap))
	i := 0
	for k := range selectedMap {
		selected[i] = k
		i++
	}
	return selected
}

func toggle(selected []int, cur int) []int {
	for i, s := range selected {
		if s == cur {
			return append(selected[:i], selected[i+1:]...)
		}
	}
	return append(selected, cur)
}

func (s *Shell) multiChoice(options []string, text string, init []int, multiResults bool) []int {
	s.multiChoiceActive = true
	defer func() { s.multiChoiceActive = false }()

	conf := s.reader.scanner.Config.Clone()

	conf.DisableAutoSaveHistory = true

	conf.FuncFilterInputRune = func(r rune) (rune, bool) {
		switch r {
		case 16:
			return -1, true
		case 14:
			return -2, true
		case 32:
			return -3, true

		}
		return r, true
	}

	var selected []int
	if multiResults {
		selected = initSelected(init, len(options))
	}

	s.ShowPrompt(false)
	defer s.ShowPrompt(true)

	// TODO this may not work on windows.
	s.Print("\033[?25l")
	defer s.Print("\033[?25h")

	cur := 0
	if len(selected) > 0 {
		cur = selected[len(selected)-1]
	}

	fd := int(os.Stdout.Fd())
	_, maxRows, err := readline.GetSize(fd)
	if err != nil {
		return nil
	}

	// move cursor to the top
	// TODO it happens on every update, however, some trash appears in history without this line
	s.Print("\033[0;0H")

	offset := fd

	update := func() {
		strs := buildOptionsStrings(options, selected, cur)
		if len(strs) > maxRows-1 {
			strs = strs[offset : maxRows+offset-1]
		}
		s.Print("\033[0;0H")
		// clear from the cursor to the end of the screen
		s.Print("\033[0J")
		s.Println(text)
		s.Print(strings.Join(strs, "\n"))
	}
	var lastKey rune
	refresh := make(chan struct{}, 1)
	listener := func(line []rune, pos int, key rune) (newline []rune, newPos int, ok bool) {
		lastKey = key
		if key == -2 {
			cur++
			if cur >= maxRows+offset-1 {
				offset++
			}
			if cur >= len(options) {
				offset = fd
				cur = 0
			}
		} else if key == -1 {
			cur--
			if cur < offset {
				offset--
			}
			if cur < 0 {
				if len(options) > maxRows-1 {
					offset = len(options) - maxRows + 1
				} else {
					offset = fd
				}
				cur = len(options) - 1
			}
		} else if key == -3 {
			if multiResults {
				selected = toggle(selected, cur)
			}
		}
		refresh <- struct{}{}
		return
	}
	conf.Listener = readline.FuncListener(listener)
	oldconf := s.reader.scanner.SetConfig(conf)

	stop := make(chan struct{})
	defer func() {
		stop <- struct{}{}
		s.Println()
	}()
	t := time.NewTicker(time.Millisecond * 200)
	defer t.Stop()
	go func() {
		for {
			select {
			case <-stop:
				return
			case <-refresh:
				update()
			case <-t.C:
				_, rows, _ := readline.GetSize(fd)
				if maxRows != rows {
					maxRows = rows
					update()
				}
			}
		}
	}()
	s.ReadLine()

	s.reader.scanner.SetConfig(oldconf)

	// only handles Ctrl-c for now
	// this can be broaden later
	switch lastKey {
	// Ctrl-c
	case 3:
		return []int{-1}
	}
	if multiResults {
		return selected
	}
	return []int{cur}
}

func buildOptionsStrings(options []string, selected []int, index int) []string {
	var strs []string
	symbol := " ❯"
	if runtime.GOOS == "windows" {
		symbol = " >"
	}
	for i, opt := range options {
		mark := "⬡ "
		if selected == nil {
			mark = " "
		}
		for _, s := range selected {
			if s == i {
				mark = "⬢ "
			}
		}
		if i == index {
			cyan := color.New(color.FgCyan).Add(color.Bold).SprintFunc()
			strs = append(strs, cyan(symbol+mark+opt))
		} else {
			strs = append(strs, "  "+mark+opt)
		}
	}
	return strs
}

// IgnoreCase specifies whether commands should not be case sensitive.
// Defaults to false i.e. commands are case sensitive.
// If true, commands must be registered in lower cases.
func (s *Shell) IgnoreCase(ignore bool) {
	s.ignoreCase = ignore
}

// ProgressBar returns the progress bar for the shell.
func (s *Shell) ProgressBar() ProgressBar {
	return s.progressBar
}

func newContext(s *Shell, cmd *Cmd, args []string) *Context {
	if cmd == nil {
		cmd = &Cmd{}
	}
	return &Context{
		Actions:     s.Actions,
		progressBar: copyShellProgressBar(s),
		Args:        args,
		RawArgs:     s.rawArgs,
		Cmd:         *cmd,
		contextValues: func() contextValues {
			values := contextValues{}
			for k := range s.contextValues {
				values[k] = s.contextValues[k]
			}
			return values
		}(),
	}
}

func copyShellProgressBar(s *Shell) ProgressBar {
	sp := s.progressBar.(*progressBarImpl)
	p := newProgressBar(s)
	p.Indeterminate(sp.indeterminate)
	p.Display(sp.display)
	p.Prefix(sp.prefix)
	p.Suffix(sp.suffix)
	p.Final(sp.final)
	p.Interval(sp.interval)
	return p
}

func getPosition() (int, int, error) {
	fd := int(os.Stdout.Fd())
	state, err := readline.MakeRaw(fd)
	if err != nil {
		return 0, 0, err
	}
	defer readline.Restore(fd, state)
	fmt.Printf("\033[6n")
	var out string
	reader := bufio.NewReader(os.Stdin)
	if err != nil {
		return 0, 0, err
	}
	for {
		b, err := reader.ReadByte()
		if err != nil || b == 'R' {
			break
		}
		if unicode.IsPrint(rune(b)) {
			out += string(b)
		}
	}
	var row, col int
	_, err = fmt.Sscanf(out, "[%d;%d", &row, &col)
	if err != nil {
		return 0, 0, err
	}

	return col, row, nil
}
