package ishell

// ErrLevel is the severity of an error.
type ErrLevel int

const (
	LevelWarn ErrLevel = iota + 1
	LevelStop
	LevelExit
	LevelPanic
)

var (
	errNoHandler = WarnErr("No handler registered for input.")
)

// ShellError is an interractive shell error
type shellError struct {
	err   string
	level ErrLevel
}

func (s shellError) Error() string {
	return s.err
}

// NewErr creates a new error with specified level
func NewErr(err string, level ErrLevel) error {
	return shellError{
		err:   err,
		level: LevelWarn,
	}
}

// WarnErr creates a Warn level error
func WarnErr(err string) error {
	return shellError{
		err:   err,
		level: LevelWarn,
	}
}

// StopErr creates a Stop level error. Shell stops if encountered.
func StopErr(err string) error {
	return shellError{
		err:   err,
		level: LevelStop,
	}
}

// ExitErr creates a Exit level error. Program terminates if encountered.
func ExitErr(err string) error {
	return shellError{
		err:   err,
		level: LevelExit,
	}
}

// PanicErr creates a Panic level error. Program panics if encountered.
func PanicErr(err string) error {
	return shellError{
		err:   err,
		level: LevelPanic,
	}
}
