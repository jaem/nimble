package nimble

import (
	"log"
	"net/http"
	"os"

	"github.com/nimgo/nimble/middleware"
)

// Default returns a new Nimble instance with the default middleware already
// in the stack.
//
// Recovery - Panic Recovery Middleware
// Logger - Request/Response Logging
// Static - Static File Serving
func Default() *Nimble {

	return New().
		UseHandler(middleware.NewRecovery()).
		UseHandler(middleware.NewLogger()).
		UseHandler(middleware.NewStatic(http.Dir("static")))
}

// New returns a new Nimble instance with no middleware preconfigured.
func New() *Nimble {
	return &Nimble{}
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

// Use adds a http.Handler onto the middleware stack.
func (n *Nimble) Use(handler http.Handler) *Nimble {
	return n.UseHandlerFunc(wrap(handler))
}

// UseFunc adds a http.HandlerFunc onto the middleware stack.
func (n *Nimble) UseFunc(handlerFunc http.HandlerFunc) *Nimble {
	return n.UseHandlerFunc(wrapHandlerFunc(handlerFunc))
}

// UseHandler adds a nimble.Handler onto the middleware stack.
func (n *Nimble) UseHandler(handler Handler) *Nimble {
	return n.UseHandlerFunc(handler.ServeHTTP)
}

// UseHandlerFunc adds a nimble.HandlerFunc function onto the middleware stack.
func (n *Nimble) UseHandlerFunc(handlerFunc HandlerFunc) *Nimble {
	if handlerFunc == nil {
		panic("handlerFunc cannot be nil")
	}
	n.handlers = append(n.handlers, handlerFunc)
	n.middleware = build(n.handlers)
	return n
}
