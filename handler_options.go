package lit

// HandlerOption is an optional config used to modify the Handler's behaviour
type HandlerOption func(*handlerConfig)

// HandlerWithProfilingDisabled disables the Handler's router's default pprof routes
func HandlerWithProfilingDisabled() HandlerOption {
	return func(c *handlerConfig) {
		c.profilingDisabled = true
	}
}
