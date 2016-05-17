package nim

import (
	"log"
	"os"
	"time"
	"net/http"
)

// NewLogger returns a new Logger instance
func NewLogger() *logger {
	return &logger{ log.New(os.Stdout, "[n.] ", 0), false }
}

func NewColorLogger() *logger {
	return &logger{ log.New(os.Stdout, "[n.] ", 0), true }
}

// Logger is a middleware that logs per request.
type logger struct {
	*log.Logger
	color bool
}

func (l *logger) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()

	next(w, r)

	ww := w.(Writer)

	clientIP := "0.0.0.0"
	latency := time.Since(start)
	method := r.Method
	path := r.URL.Path
	statusCode := ww.Status()
	status := http.StatusText(ww.Status())

	if l.color {
		statusColor := colorForStatus(statusCode)
		methodColor := colorForMethod(method)
		l.Printf("%v |%s %s %3d | %13v | %s |%s %s %-7s [%-2s] %s \n",
			start.Format("2006/01/02 - 15:04:05 -0700"),
			statusColor, reset, statusCode,
			latency,
			clientIP,
			methodColor, reset, method,
			status,
			path,
		)
	} else {
		textColor := colorForText(statusCode)
		l.Printf("%s%v | %13v | %15s | %-7s | %3d [%-2s] %s %s\n",
			textColor,
			start.Format("2006/01/02 - 15:04:05 -0700"),
			latency,
			clientIP,
			method,
			statusCode,
			status,
			path,
			reset,
		)
	}
}

// https://github.com/shiena/ansicolor
var (
	green   = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	white   = string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
	yellow  = string([]byte{27, 91, 57, 55, 59, 52, 51, 109})
	red     = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	blue    = string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
	magenta = string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
	cyan    = string([]byte{27, 91, 57, 55, 59, 52, 54, 109})
	reset   = string([]byte{27, 91, 48, 109})

	greenTT  = "\x1b[92m"
	whiteTT  = "\x1b[37m"
	yellowTT = "\x1b[33m"
	redTT		 = "\x1b[31m"
)

func colorForText(code int) string {
	switch {
	case code >= 200 && code < 300:
		return greenTT
	case code >= 300 && code < 400:
		return whiteTT
	case code >= 400 && code < 500:
		return yellowTT
	default:
		return redTT
	}
}

func colorForStatus(code int) string {
	switch {
	case code >= 200 && code < 300:
		return green
	case code >= 300 && code < 400:
		return white
	case code >= 400 && code < 500:
		return yellow
	default:
		return red
	}
}

func colorForMethod(method string) string {
	switch method {
	case "GET":
		return green
	case "POST":
		return cyan
	case "PUT":
		return yellow
	case "DELETE":
		return red
	case "PATCH":
		return green
	case "HEAD":
		return magenta
	case "OPTIONS":
		return white
	default:
		return reset
	}
}