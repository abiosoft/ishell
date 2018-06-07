// +build windows

package ishell

import (
	"github.com/abiosoft/readline"
)

func clearScreen(s *Shell) error {
	return readline.ClearScreen(s.writer)
}
