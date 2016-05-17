package nim

import (
	"net/http"
	"path"
	"strings"
)

// NewStatic returns a new instance of Static
func NewStatic(directory http.FileSystem) *static {
	return &static{
		dir:       directory,
		prefix:    "",
		indexFile: "index.html",
	}
}

// Static is a middleware to serves static files in the given directory/filesystem.
type static struct {
	// Dir is the directory to serve static files from
	dir       http.FileSystem
	// Prefix is the optional prefix used to serve the static directory content
	prefix    string
	// IndexFile defines which file to serve as index if it exists.
	indexFile string
}

func (s *static) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if r.Method != "GET" && r.Method != "HEAD" {
		next(rw, r)
		return
	}
	file := r.URL.Path
	// if we have a prefix, filter requests by stripping the prefix
	if s.prefix != "" {
		if !strings.HasPrefix(file, s.prefix) {
			next(rw, r)
			return
		}
		file = file[len(s.prefix):]
		if file != "" && file[0] != '/' {
			next(rw, r)
			return
		}
	}
	f, err := s.dir.Open(file)
	if err != nil {
		// Handle multiple requests from modern browsers if missing the favicon.ico
		if file != "/favicon.ico" {
			next(rw, r)
		}

		// discard the error?
		return
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		next(rw, r)
		return
	}

	// try to serve index file
	if fi.IsDir() {
		// redirect if missing trailing slash
		if !strings.HasSuffix(r.URL.Path, "/") {
			http.Redirect(rw, r, r.URL.Path + "/", http.StatusFound)
			return
		}

		file = path.Join(file, s.indexFile)
		f, err = s.dir.Open(file)
		if err != nil {
			next(rw, r)
			return
		}
		defer f.Close()

		fi, err = f.Stat()
		if err != nil || fi.IsDir() {
			next(rw, r)
			return
		}
	}

	http.ServeContent(rw, r, file, fi.ModTime(), f)
}
