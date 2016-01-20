package ishell

import (
	"bufio"
	"sync"

	"github.com/bobappleyard/readline"
)

type (
	lineString struct {
		line string
		err  error
	}

	shellReader struct {
		scanner   *bufio.Reader
		consumers []chan lineString
		reading   bool
		sync.Mutex
	}
)

func (s *shellReader) ReadLine(consumer chan lineString, prompt string) {
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
		line, err := readline.String(prompt)
		ls := lineString{line, err}
		s.Lock()
		defer s.Unlock()
		for i := range s.consumers {
			c := s.consumers[i]
			go func(c chan lineString) {
				c <- ls
			}(c)
		}
		readline.AddHistory(line)
		s.reading = false
	}()

}
