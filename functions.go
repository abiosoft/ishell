package ishell

// CmdFunc represents a command function that is called after an input to the shell.
// The shell input is split into command and arguments like cli args.
// The shell will print output if output != "".
type CmdFunc func(command string, args []string) (output string, err error)

func exitFunc(s *Shell) CmdFunc {
	return func(command string, args []string) (output string, err error) {
		s.Stop()
		return "", nil
	}
}

func helpFunc(s *Shell) CmdFunc{
	return func(command string, args []string) (output string, err error) {
		s.PrintCommands()
		return "", nil
	}
}

func addDefaultFuncs(s *Shell){
	s.Register("exit", exitFunc(s))
	s.Register("help", helpFunc(s))
}