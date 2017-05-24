package nimble

import (
	"log"
	"os"
	"net/http"
)

// New returns a new Nimble instance with no middleware preconfigured.
func New() *Nimble {
	return &Nimble{}
}

// Default returns a new Nimble instance with the default middleware already
// in the stack.
//
// Recovery - Panic Recovery Middleware
// Logger - Request/Response Logging
// Static - Static File Serving
func Default() *Nimble {
	return New().
		UseHandler(NewRecovery()).
		UseHandler(NewLogger()).
		UseHandler(NewStatic(http.Dir("static")))
}

// Nimble is a stack of Middleware Handlers that can be invoked as an http.Handler.
// The middleware stack is run in the sequence that they are added to the stack.
type Nimble struct {
	middleware middleware
	handlers   []HandlerFunc
}

// Run is a convenience function that runs the nimble stack as an HTTP
// server. The addr string takes the same format as http.ListenAndServe.
func (n *Nimble) Run(addr ...string) {
	l := log.New(os.Stdout, "[n.] ", 0)
	address := detectAddress(addr...)
	l.Printf("Server listening on %s", address)
	l.Fatal(http.ListenAndServe(address, n))
}

// Nimble itself is a http.Handler. This allows it to used as a substack manager
func (n *Nimble) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, ok := w.(Writer); ok { // handle substacks
		n.middleware.serve(w, r)
	} else {
		n.middleware.serve(newWriter(w), r)
	}
}

// UseHandler adds a http.Handler onto the middleware stack.
func (n *Nimble) Use(handler http.Handler) *Nimble {
	return n.UseHandlerFunc(wrap(handler))
}

// UseHandlerFunc adds a http.HandlerFunc onto the middleware stack.
func (n *Nimble) UseFunc(handlerFunc http.HandlerFunc) *Nimble {
	return n.UseHandlerFunc(wrapHandlerFunc(handlerFunc))
}

// Use adds a nimble.Handler onto the middleware stack.
func (n *Nimble) UseHandler(handler Handler) *Nimble {
	return n.UseHandlerFunc(handler.ServeHTTP)
}

// UseFunc adds a nimble.HandlerFunc function onto the middleware stack.
func (n *Nimble) UseHandlerFunc(handlerFunc HandlerFunc) *Nimble {
	if (handlerFunc == nil) {
		panic("NimbleHandleFunc cannot be nil")
	}
	n.handlers = append(n.handlers, handlerFunc)
	n.middleware = build(n.handlers)
	return n
}

const (
	// DefaultAddress is used if no other is specified.
	defaultServerAddress = ":8080"
)

func detectAddress(addr ...string) string {
	if len(addr) > 0 {
		return addr[0]
	}
	if port := os.Getenv("PORT"); port != "" {
		return ":" + port
	}
	return defaultServerAddress
}
