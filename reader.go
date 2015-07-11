package ishell

import (
	"bufio"
	"sync"
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
		line, err := s.scanner.ReadString('\n')
		// remove training '\n'
		if err == nil {
			line = line[:len(line)-1]
		}
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
