package ishell

import "sync"

// Context is an ishell context. It embeds ishell.Actions.
type Context struct {
	*contextValues
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

// ProgressBar returns the progress bar for the current shell context.
func (c *Context) ProgressBar() ProgressBar {
	return c.progressBar
}

// contextValues is the map for values in the context.
type contextValues struct {
	vals map[string]interface{}
	*sync.RWMutex
}

// Get returns the value associated with this context for key, or nil
// if no value is associated with key. Successive calls to Get with
// the same key returns the same result.
func (c *contextValues) Get(key string) interface{} {
	c.RLock()
	defer c.RUnlock()
	return c.vals[key]
}

// Set sets the key in this context to value.
func (c *contextValues) Set(key string, value interface{}) {
	if c.vals == nil {
		c.vals = make(map[string]interface{})
		c.RWMutex = &sync.RWMutex{}
	}
	c.Lock()
	c.vals[key] = value
	c.Unlock()
}

// Del deletes key and its value in this context.
func (c *contextValues) Del(key string) {
	c.Lock()
	delete(c.vals, key)
	c.Unlock()
}

// Keys returns all keys in the context.
func (c *contextValues) Keys() (keys []string) {
	c.RLock()
	for key := range c.vals {
		keys = append(keys, key)
	}
	c.RUnlock()
	return
}
