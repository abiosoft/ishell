// +build windows

package ishell

import (
	"github.com/chzyer/readline"
)

func clearScreen(s *Shell) error {
	return readline.ClearScreen(s.writer)
}
