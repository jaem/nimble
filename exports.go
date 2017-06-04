package nim

import (
	"log"
	"net/http"
	"os"

	"github.com/nimgo/nim/nimble"
	"github.com/nimgo/nim/nimware"
)

// Default returns a new Nimble instance with the default middleware already
// in the stack.
//
// Recovery - Panic Recovery Middleware
// Logger - Request/Response Logging
// Static - Static File Serving
func Default() *nimble.Nimble {
	return nimble.New().
		WithHandler(nimware.NewRecovery()).
		WithHandler(nimware.NewColorLogger()).
		WithHandler(nimware.NewStatic(http.Dir("static")))
}

// New returns a new Nimble instance with no middleware preconfigured.
func New() *nimble.Nimble {
	return nimble.New()
}

// Run is a convenience function that runs the nimble stack as an HTTP
// server. The addr string takes the same format as http.ListenAndServe.
func Run(n *nimble.Nimble, addr ...string) {
	l := log.New(os.Stdout, "[n.] ", 0)
	address := detectAddress(addr...)
	l.Printf("Server is listening on %s", address)
	l.Fatal(http.ListenAndServe(address, n))
}
