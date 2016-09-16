package ishell

// errLevel is the severity of an error.
type errLevel int

const (
	levelWarn errLevel = iota + 1
	levelStop
	levelExit
	levelPanic
)

var (
	errNoHandler = WarnErr("No handler registered for input.")
)

// shellError is an interractive shell error
type shellError struct {
	err   string
	level errLevel
}

func (s shellError) Error() string {
	return s.err
}

// WarnErr creates a Warn level error
func WarnErr(err string) error {
	return shellError{
		err:   err,
		level: levelWarn,
	}
}

// StopErr creates a Stop level error. Shell stops if encountered.
func StopErr(err string) error {
	return shellError{
		err:   err,
		level: levelStop,
	}
}

// ExitErr creates a Exit level error. Program terminates if encountered.
func ExitErr(err string) error {
	return shellError{
		err:   err,
		level: levelExit,
	}
}

// PanicErr creates a Panic level error. Program panics if encountered.
func PanicErr(err string) error {
	return shellError{
		err:   err,
		level: levelPanic,
	}
}
