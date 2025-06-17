//go:build darwin || dragonfly || freebsd || (linux && !appengine) || netbsd || openbsd || solaris
// +build darwin dragonfly freebsd linux,!appengine netbsd openbsd solaris

package ishell

func clearScreen(s *Shell) error {
	_, err := s.writer.Write([]byte("\033[H\033[2J"))
	return err
}
