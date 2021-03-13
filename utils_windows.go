// +build windows

package ishell

import (
	"github.com/arrow2nd/readline"
)

func clearScreen(s *Shell) error {
	return readline.ClearScreen(s.writer)
}
