//go:build darwin || dragonfly || freebsd || (linux && !appengine) || netbsd || openbsd || solaris
// +build darwin dragonfly freebsd linux,!appengine netbsd openbsd solaris

package ishell

func clearScreen(s *Shell) error {
	// PR: ClearScreen function corrected (@shadowshot-x)
	// link: https://github.com/abiosoft/ishell/pull/146
	_, err := s.writer.Write([]byte("\033[H\033[2J"))
	return err
}
