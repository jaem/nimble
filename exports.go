package nimble

import (
	"net/http"

	"github.com/nimgo/nimble/nim"
	"github.com/nimgo/nimble/nimware"
)

// Default returns a new Nimble instance with the default middleware already
// in the stack.
//
// Recovery - Panic Recovery Middleware
// Logger - Request/Response Logging
// Static - Static File Serving
func Default() *nim.Nimble {
	return nim.New().
		WithHandler(nimware.NewRecovery()).
		WithHandler(nimware.NewColorLogger()).
		WithHandler(nimware.NewStatic(http.Dir("static")))
}

// New returns a new Nimble instance with no middleware preconfigured.
func New() *nim.Nimble {
	return nim.New()
}
