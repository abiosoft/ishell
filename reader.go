package ishell

import (
	"sync"

	"github.com/chzyer/readline"
)

type (
	lineString struct {
		line string
		err  error
	}

	shellReader struct {
		scanner   *readline.Instance
		consumers []chan lineString
		reading   bool
		sync.Mutex
	}
)

func (s *shellReader) ReadLine(consumer chan lineString) {
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
		line, err := s.scanner.Readline()
		ls := lineString{line, err}
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
