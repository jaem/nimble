package nimble

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"net/http"
)

// Recovery is a middleware that attempts to recover from panics and writes a 500 if there was one.
type Recovery struct {
	Logger     *log.Logger
	PrintStack bool
	StackAll   bool
	StackSize  int
}

// NewRecovery returns a new instance of Recovery
func NewRecovery() *Recovery {
	return &Recovery{
		Logger:     log.New(os.Stdout, "[nim.] ", 0),
		PrintStack: true,
		StackAll:   false,
		StackSize:  1024 * 8,
	}
}

func (rec *Recovery) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	defer func() {
		if err := recover(); err != nil {
			stack := make([]byte, rec.StackSize)
			stack = stack[:runtime.Stack(stack, rec.StackAll)]

			f := "RECOVER: %s\n%s"
			rec.Logger.Printf(f, err, stack)
			if rec.PrintStack {
				fmt.Fprintf(w, f, err, stack)
			}

			w.WriteHeader(http.StatusInternalServerError)
		}
	}()

	next(w, r)
}
