package nim

import (
	"net/http"

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
