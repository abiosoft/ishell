package ishell

import (
	"bytes"
	"strings"
	"sync"

	"github.com/abiosoft/readline"
)

type (
	lineString struct {
		line string
		err  error
	}

	shellReader struct {
		scanner      *readline.Instance
		consumers    chan lineString
		reading      bool
		readingMulti bool
		buf          *bytes.Buffer
		prompt       string
		multiPrompt  string
		showPrompt   bool
		completer    readline.AutoCompleter
		sync.Mutex
	}
)

// rlPrompt returns the proper prompt for readline based on showPrompt and
// prompt members.
func (s *shellReader) rlPrompt() string {
	if s.showPrompt {
		if s.readingMulti {
			return s.multiPrompt
		}
		return s.prompt
	}
	return ""
}

func (s *shellReader) readPasswordErr() (string, error) {
	prompt := ""
	if s.buf.Len() > 0 {
		prompt = s.buf.String()
		s.buf.Truncate(0)
	}
	password, err := s.scanner.ReadPassword(prompt)
	return string(password), err
}

func (s *shellReader) readPassword() string {
	password, _ := s.readPasswordErr()
	return password
}

func (s *shellReader) setMultiMode(use bool) {
	s.readingMulti = use
}

func (s *shellReader) readLine(consumer chan lineString) {
	s.Lock()
	defer s.Unlock()

	// already reading
	if s.reading {
		return
	}
	s.reading = true
	// start reading

	// detect if print is called to
	// prevent readline lib from clearing line.
	// use the last line as prompt.
	// TODO find better way.
	shellPrompt := s.prompt
	prompt := s.rlPrompt()
	if s.buf.Len() > 0 {
		lines := strings.Split(s.buf.String(), "\n")
		if p := lines[len(lines)-1]; strings.TrimSpace(p) != "" {
			prompt = p
		}
		s.buf.Truncate(0)
	}

	// use printed statement as prompt
	s.scanner.SetPrompt(prompt)

	line, err := s.scanner.Readline()

	// reset prompt
	s.scanner.SetPrompt(shellPrompt)

	ls := lineString{string(line), err}
	consumer <- ls
	s.reading = false
}
