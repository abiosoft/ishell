package ishell

// CmdFunc represents a command function that is called after an input to the shell.
// The shell input is split into command and arguments like cli args and the arguments
// are passed to this function. The shell will print output if output != "".
type CmdFunc func(args ...string) (output string, err error)

func exitFunc(s *Shell) CmdFunc {
	return func(args ...string) (string, error) {
		s.Stop()
		return "", nil
	}
}

func helpFunc(s *Shell) CmdFunc {
	return func(args ...string) (string, error) {
		s.PrintCommands()
		return "", nil
	}
}

func clearFunc(s *Shell) CmdFunc {
	return func(args ...string) (string, error) {
		err := s.ClearScreen()
		return "", err
	}
}

func addDefaultFuncs(s *Shell) {
	s.Register("exit", exitFunc(s))
	s.Register("help", helpFunc(s))
	s.Register("clear", clearFunc(s))
}
