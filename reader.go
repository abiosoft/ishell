package ishell

import (
	"bytes"
	"sync"

	"gopkg.in/readline.v1"
)

type (
	lineString struct {
		line string
		err  error
	}

	shellReader struct {
		scanner      *readline.Instance
		consumers    []chan lineString
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

func (s *shellReader) readPassword() string {
	password, _ := s.scanner.ReadPassword("")
	return string(password)
}

func (s *shellReader) setMultiMode(use bool) {
	s.readingMulti = use
}

func (s *shellReader) readLine(consumer chan lineString) {
	s.Lock()
	defer s.Unlock()
	s.consumers = append(s.consumers, consumer)
	// already reading
	if s.reading {
		return
	}
	s.reading = true
	// start reading
	go func() {
		// detect if print is called to
		// prevent readline lib from clearing line.
		// TODO find better way.
		shellPrompt := s.prompt
		prompt := s.rlPrompt()
		if s.buf.Len() > 0 {
			prompt += s.buf.String()
			s.buf.Truncate(0)
		}

		// use printed statement as prompt
		s.scanner.SetPrompt(prompt)

		line, err := s.scanner.Readline()

		// reset prompt
		s.scanner.SetPrompt(shellPrompt)

		ls := lineString{string(line), err}
		s.Lock()
		defer s.Unlock()
		for i := range s.consumers {
			c := s.consumers[i]
			go func(c chan lineString) {
				c <- ls
			}(c)
		}
		s.reading = false
	}()

}
