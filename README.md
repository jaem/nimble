# Nim (aka Nimble) [![GoDoc](https://godoc.org/github.com/nimgo/nim?status.svg)](http://godoc.org/github.com/nimgo/nim) [![Build Status](https://travis-ci.org/nimgo/nim.svg?branch=master)](https://travis-ci.org/nimgo/nim)

Nim is a lightweight middleware stack engine that encourages dev authors to use `net/http` Handlers.

It uses a chaining approach to web middleware in Go. It is inspired by Negroni and is similar
to how Express (nodejs) handles its middleware.

##### What it is:

* Lightweight and non-intrusive to use as an stack engine.
* Helps you to manage your webstack cleanly.
* Allows you to combine the specialization of various open-source middleware packages.
* (NOT) Nimble is not a web framework. It allows you to build your custom framework with
  your selected middleware. It is a no-frills engine to manage the middleware stack.

  An analogy would be that a web framework is a fixed 10-course dinner set; Nimble lets you have
  a 10-course dinner set, but you can custom choose your dishes. So it depends on what you need.

##### What is middleware?

Middleware is actually a very broad fancy term. For starters, you can see it basically
as layers of reusable codes that run sequentially per request served. Nimble can manage
the flow before it gets to your router, or within your router for specific routes
(like authentication checks).

For Nimble, it can be as small as a this function that looks like this:
~~~ go
func iAmMiddleware(w http.ResponseWriter, r *http.Request) {
    fmt.Println("I am middleware")
}
~~~
or a fat middleware that looks like this:
~~~ go
struct FatMiddleWare {}
func (mw *FatMiddleWare) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    ... holy cow I'm doing plenty of stuff ...
    ... call this
    ... check that
    ... phew done ('sweat)
}
~~~

##### Features:

* Stacks and sub-stacks: Because of its small size and compatibility with net/http,
  Nimble allows you to create large stacks at the application-level, and/or
  independent sub-stacks you can apply at the router-level.

* Routing: Nimble is compatible with routers such as gorilla/mux, net/http, httprouter, etc.
  Nimble keeps the router separate, so you can use any router easily according to your needs.
  The Go community already has a number of great http routers available, and using them is
  usually straightforward if they use the signature from net/http.Handler.

  For example in a small application,
  net/http, httprouter, fastrouter could easily handle your needs.

  However, in a mid-large application,
  you might choose gorilla/mux (beyond just speed benchmarks), because it can handle things like:
  `"/articles/{category}/{id:[0-9]+}"`

* Context: Nimble provisions a net/context per request, using gorilla/context.
  This keeps the stack compatible with net/http without altering existing apis, but allows
  the freedom to use context information. The context is stored once and is available
  throughout the request lifecycle. Having context is useful for authentication/authorization
  processes, or pre-request handling.

  The purpose of a request context that uses net/context to allow dev authors
  to start familiarizing with the use of net/context. Nimble avoids trying to create
  a Nimble specific context to allow for generality in your codes.

  In addition, recent updates in the golang sphere (as of May 2016) indicates that
  net/context will be attached to http.request. So sticking to
  `func(w http.ResponseWriter, r *http.Request)` as close as possible
  will be better moving forward.

## Getting Started

Step 1: Install Go and setting up your [GOPATH](http://golang.org/doc/code.html#GOPATH).

Step 2: Then install the nimble package (**go 1.1** and greater is required):
~~~
go get github.com/jaem/nimble
~~~
Step 3: Alright, here we go with examples.

##### Example 1 (Basic): A simple web server

Here, we are using golang/net/http to route your request.

* Copy the following code to a file `myserver.go`
* Run your server with: `# go run myserver.go`
* Open `http://localhost:3000` with your browser

Results: There, you now have a Golang net/http webserver running on `localhost:3000`.

~~~ go
package main

import (
  "net/http"
  "github.com/jaem/nimble"
)

func main() {
  router := http.NewServeMux()
  router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Welcome to your server!"))
  })

  n := nimble.Default()
  n.Use(router)
  n.Run(":3000")
}
~~~

##### Example 2 (Basic): Creating middleware

Create a simple middleware function, using gorilla/mux for routing this time.

* Similarly, copy the code to `myserver.go` and run ```# go run myserver.go```
* Open `http://localhost:3000` with your browser
* Open `http://localhost:3000/about` with your browser

Results: You should see that the middleware was run on both webpages.

~~~ go
package main

import (
  "net/http"

  "github.com/gorilla/mux"
  "github.com/jaem/nimble"
)

func main() {
  router := mux.NewRouter()
  router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Welcome your server!"))
  })
  router.HandleFunc("/about", aboutFunc)

  n := nimble.Default()
  n.UseFunc(myMiddleware)
  n.Use(router)
  n.Run(":3000")
}

func aboutFunc(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("This is a lean, mean server."))
}

func myMiddleware(w http.ResponseWriter, r *http.Request) {
  w.Write([]byte("A middleware that always runs per http request.\n\n"))
}
~~~

##### Example 3 (Advanced): Sub-stacks

Use sub-stacks for certain routes. Use cases like in authentication/admin access controls, etc.

* Similarly, copy the code to `myserver.go` and run ```# go run myserver.go```
* Open `http://localhost:3000` with your browser
* Open `http://localhost:3000/about` with your browser
* Open `http://localhost:3000/p/iron_man` with your browser
* Open `http://localhost:3000/p/captain_america` with your browser

Results: You should see that the sub-Middleware only ran on `http://localhost:3000/p/iron_man`
and on `http://localhost:3000/p/captain_america`, as defined by the subrouter.

~~~ go
package main

import (
  "net/http"

  "github.com/gorilla/mux"
  "github.com/jaem/nimble"
)

func main() {
  router := mux.NewRouter()
  router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Welcome to Nimble!"))
  })
  router.HandleFunc("/about", aboutFunc)

  subrouter := mux.NewRouter()
  subrouter.HandleFunc("/p/iron_man", saysHi("Iron Man"))
  subrouter.HandleFunc("/p/captain_america", saysHi("Captain America"))
  router.PathPrefix("/p").Handler(nimble.New().
		UseFunc(subMiddleware).
		Use(subrouter),
	)

  n := nimble.Default()
  n.UseFunc(myMiddleware)
  n.Use(router)
  n.Run(":3000")
}

func aboutFunc(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("This is a lean, mean server."))
}

func myMiddleware(w http.ResponseWriter, r *http.Request) {
  w.Write([]byte("A middleware that always runs per http request.\n\n"))
}

func saysHi(who string) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte(who + " says, 'Hi y'all!'"))
  }
}

func subMiddleware(w http.ResponseWriter, r *http.Request) {
  w.Write([]byte("SubMiddleware: Presenting to you \n\n"))
}
~~~

##### Example 4 (Advanced): Using context

A net/context is provisioned per request.

* Similarly, copy the code to `myserver.go` and run ```# go run myserver.go```
* Open `http://localhost:3000/p/iron_man` with your browser
* Open `http://localhost:3000/p/captain_america` with your browser

Results: You should see that the value stored in the request context in
myMiddleware was retrieved separately in sub-Middleware as per the request.

~~~ go
package main

import (
  "net/http"

  "golang.org/x/net/context"

  "github.com/gorilla/mux"
  "github.com/jaem/nimble"
)

func main() {
  router := mux.NewRouter()
  router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Welcome to Nimble!"))
  })
  router.HandleFunc("/about", aboutFunc)

  subrouter := mux.NewRouter()
  subrouter.HandleFunc("/p/iron_man", saysHi("Iron Man"))
  subrouter.HandleFunc("/p/captain_america", saysHi("Captain America"))
  router.PathPrefix("/p").Handler(nimble.New().
		UseFunc(subMiddleware).
		Use(subrouter),
	)

  n := nimble.Default()
  n.UseFunc(myMiddleware)
  n.Use(router)
  n.Run(":3000")
}

func aboutFunc(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("This is a lean, mean server."))
}

func myMiddleware(w http.ResponseWriter, r *http.Request) {
  w.Write([]byte("A middleware that always runs per http request.\n\n"))

  c := nimble.GetContext(r)
  c = context.WithValue(c, "key", "the Avengers")
  nimble.SetContext(r, c)
}

func saysHi(who string) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte(who + " says, 'Hi y'all!'"))
  }
}

func subMiddleware(w http.ResponseWriter, r *http.Request) {
  c := nimble.GetContext(r)
  if value, ok := c.Value("key").(string); ok {
    w.Write([]byte("SubMiddleware: Presenting to you " + value + "\n\n"))
	}
}
~~~

## Additional Notes

##### Note 1: Using a pre-defined context

If you need to pass a pre-defined context, then use
`n := nimble.DefaultWithContext(appContext)`

~~~ go
func main() {
    appContext := context.TODO()
    appContext = context.WithValue(appContext, "dbname", "....")
    appContext = context.WithValue(appContext, "dbpass", "....")
    ...

    n := nimble.DefaultWithContext(appContext) // Instead of nimble.Default()
    ...
}
~~~

##### Note 2: Chaining middleware

The aim of nimble is to help you chain your middleware into stacks, to make it easier to create
reusable codes when authoring your web server. This is how you can use it.

~~~
func main() {
    stack := nimble.New()            // creates an empty stack
                .Use(...)            // chains http.Handler
                .UseFunc(...)        // chains fn(http.ResponseWriter, *http.Request)
                .UseHandler(...)     // chains nimble.Handler
                .UseHandlerFunc(...) // chains fn(http.ResponseWriter, *http.Request, next http.HandlerFunc)
}
~~~

##### Note 3: Routing

Nimble is designed to support the magical routers built by the Go community.
They should work so long as net/http is supported.

Benchmarks: https://github.com/julienschmidt/go-http-routing-benchmark.

Speed is one aspect of choosing a router. Generally, the more features in a router,
the longer it would take. Use the router you need, but it should not affect too much
the way you manage your middleware stacks. For instance, integrating with
gorilla/mux or httprouter are as follow:

``` go
// "github.com/gorilla/mux"
func main() {
    router := mux.NewRouter()
    router.HandleFunc("/", indexHandler)

    n := nimble.Default()
    n.Use(middleware) // Use middleware
    n.Use(router) // router goes last in the nimble stack
    n.Run(":3000")
}
```

``` go
// "github.com/julienschmidt/httprouter"
func main() {
    router := httprouter.New()
    router.GET("/", indexHandler)

    n := nimble.Default()
    n.Use(middleware) // Use middleware
    n.Use(router) // again, router goes last in the nimble stack
    n.Run(":3000")
}
```

##### Note 4: Running/Shutdown gracefully

When shutting down the server, you want it to be done when requests are fulfilled.
Otherwise, you can end up with a half-done request that can cause data inconsistency.
Graceful is a lightweight tool that deals with this
http://github.com/tylerb/graceful

Note: During development, you can simply use `n.Run(":8000")`

```
// "github.com/tylerb/graceful"
func main() {

   n := nimble.Default()
   n.Use(middleware)
   n.Use(router)

   graceful.Run(":3000", 10 * time.Second, n)
}
```

##### Note 5: Nimble factory

There are 3 ways to get a Nimble instance:

`nimble.Default()` provides some default middleware that is useful for most applications:
* `Recovery` - Panic Recovery Middleware.
* `Logging` - Request/Response Logging Middleware.
* `Static` - Static File serving under the "public" directory.

`nimble.DefaultWithContext(context)` provides the default middleware, and allows you
 to provide pre-defined context to support your application.

`nimble.New()` creates a no frills emptystack. This is useful in mainly two ways.
 One, if you want to define your own middleware. And two, for sub-middleware stacks
 to be used only in specific routes. See example 3 above.

## Need Help?

Drop a note at GitHub issues for Nimble on faqs, bug reports and pull requests.

## License

The MIT License (MIT). See the LICENSE file for details.

## Links

* Graceful - http://github.com/tylerb/graceful
* Negroni - http://github.com/codegangsta/negroni
* Gorilla/mux - http://github.com/gorilla/mux
* HttpRouter - http://github.com/julienschmidt/httprouter
* Context - http://go-review.googlesource.com/#/c21496
* Context - http://github.com/golang/go/issues/14660
