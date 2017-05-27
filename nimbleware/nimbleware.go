package nimbleware

import (
	"net/http"

	"github.com/nimgo/nimble"
)

// Default returns a new Nimble instance with the default middleware already
// in the stack.
//
// Recovery - Panic Recovery Middleware
// Logger - Request/Response Logging
// Static - Static File Serving
func Default() *nimble.Nimble {
	return nimble.New().
		WithHandler(NewRecovery()).
		WithHandler(NewColorLogger()).
		WithHandler(NewStatic(http.Dir("static")))
}
