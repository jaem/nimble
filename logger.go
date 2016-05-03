package nimble

import (
	"log"
	"os"
	"time"
	"net/http"
)

// Logger is a middleware that logs per request.
type Logger struct {
	*log.Logger
}

// NewLogger returns a new Logger instance
func NewLogger() *Logger {
	return &Logger{log.New(os.Stdout, "[nim.] ", 0)}
}

func (l *Logger) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()
	l.Printf("Request started %s %s", r.Method, r.URL.Path)

	next(w, r)

	res := w.(Writer)
	l.Printf("Request completed %v %s in %v\n\n", res.Status(), http.StatusText(res.Status()), time.Since(start))
}
