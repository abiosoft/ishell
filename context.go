package ishell

// Context is an ishell context. It embeds ishell.Actions.
type Context struct {
	values      map[string]interface{}
	progressBar ProgressBar
	err         error

	// Args is command arguments.
	Args []string

	// RawArgs is unprocessed command arguments.
	RawArgs []string

	// Cmd is the currently executing command. This is empty for NotFound and Interrupt.
	Cmd Cmd

	Actions
}

// Err informs ishell that an error occurred in the current
// function.
func (c *Context) Err(err error) {
	c.err = err
}

// Get returns the value associated with this context for key, or nil
// if no value is associated with key. Successive calls to Get with
// the same key returns the same result.
func (c *Context) Get(key string) interface{} {
	return c.values[key]
}

// Set sets the key in this context to value.
func (c *Context) Set(key string, value interface{}) {
	if c.values == nil {
		c.values = make(map[string]interface{})
	}
	c.values[key] = value
}

// Del deletes key and its value in this context.
func (c *Context) Del(key string) {
	delete(c.values, key)
}

// ProgressBar returns the progress bar for the current shell context.
func (c *Context) ProgressBar() ProgressBar {
	return c.progressBar
}
