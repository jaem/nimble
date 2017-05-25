package interfaces

import (
	"net/http"
)

// Writer is the interface response wrapper that provides extra information about
// the response. It is recommended that middleware handlers use this construct to wrap a responsewriter
// if the functionality calls for it.
type (
	Writer interface {
		http.ResponseWriter
		http.Flusher
		// Status returns the status code of the response or 0 if the response has not been written.
		Status() int
		// Written returns whether or not the ResponseWriter has been written.
		Written() bool
		// Size returns the size of the response body.
		Size() int
		// Before allows for a function to be called before the ResponseWriter has been written to. This is
		// useful for setting headers or any other operations that must happen before a response has been written.
		Before(func(Writer))
	}
)
