//go:build darwin || dragonfly || freebsd || (linux && !appengine) || netbsd || openbsd || solaris
// +build darwin dragonfly freebsd linux,!appengine netbsd openbsd solaris

package ishell

import (
	"github.com/arrow2nd/readline"
)

func clearScreen(s *Shell) error {
	_, err := s.writer.Write([]byte("\033[H\033[2J"))
	return err
}
