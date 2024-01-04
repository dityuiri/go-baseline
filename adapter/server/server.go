package server

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

var (
	// RequestIDHeaderName is a key used to get `X-Request-Id` from request header
	RequestIDHeaderName = middleware.RequestIDHeader
)

// NewServer create Chi instance that contains Chi router
// ctx is passed and store in Chi instance
// appName is used for service name
func NewServer(ctx context.Context, config *Configuration, options ...Option) IServer {
	router := chi.NewRouter()

	server := &Server{
		Context: ctx,
		Config:  config,

		Router: router,
	}

	// Make sure we always include RequestID.
	router.Use(middleware.RequestID)

	router.Use(middleware.RealIP)

	server.Configure(options...)

	return server
}

func (i *Server) GetRouter() *chi.Mux {
	return i.Router
}

// Configure will apply additional options
func (i *Server) Configure(options ...Option) {
	for _, opt := range options {
		opt(i)
	}

	i.Router.Use(middleware.Recoverer)

	if i.profilerPath != "" {
		// enable profiling for service https://github.com/go-chi/chi/blob/master/middleware/profiler.go
		i.Router.Mount(i.profilerPath, middleware.Profiler())
	}
}

// Serve will spawn the server listening on the port
func (i *Server) Serve() error {
	address := fmt.Sprintf("%s:%d", i.Config.Host, i.Config.Port)

	server := &http.Server{
		Addr:    address,
		Handler: i.Router,
	}

	var e error

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			e = err
		}
	}()

	// Wait for the go routine to return from error early
	runtime.Gosched()

	if e == nil {
		i.server = server
	}

	return e
}

// Close shuts down the server
func (i *Server) Close() error {
	// close observer to flush logs and trace

	if i.server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(i.Config.ShutdownTimeout)*time.Second)
		defer cancel()

		if err := i.server.Shutdown(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (i *Server) processHandlerOption(method string, pattern string, opts ...HandlerOption) {
	handler := &Handler{}

	for _, opt := range opts {
		opt(handler)
	}
}

// Get adds the route `pattern` that matches a GET http method to
// execute the `handlerFn` http.HandlerFunc.
func (i *Server) Get(pattern string, f http.HandlerFunc, opts ...HandlerOption) {
	i.processHandlerOption(http.MethodGet, pattern, opts...)

	i.Router.Get(pattern, f)
}

// Post adds the route `pattern` that matches a POST http method to
// execute the `handlerFn` http.HandlerFunc.
func (i *Server) Post(pattern string, f http.HandlerFunc, opts ...HandlerOption) {
	i.processHandlerOption(http.MethodPost, pattern, opts...)

	i.Router.Post(pattern, f)
}

// Put adds the route `pattern` that matches a PUT http method to
// execute the `handlerFn` http.HandlerFunc.
func (i *Server) Put(pattern string, f http.HandlerFunc, opts ...HandlerOption) {
	i.processHandlerOption(http.MethodPut, pattern, opts...)

	i.Router.Put(pattern, f)
}

// Patch adds the route `pattern` that matches a PATCH http method to
// execute the `handlerFn` http.HandlerFunc.
func (i *Server) Patch(pattern string, f http.HandlerFunc, opts ...HandlerOption) {
	i.processHandlerOption(http.MethodPatch, pattern, opts...)

	i.Router.Patch(pattern, f)
}

// Options adds the route `pattern` that matches a OPTIONS http method to
// execute the `handlerFn` http.HandlerFunc.
func (i *Server) Options(pattern string, f http.HandlerFunc, opts ...HandlerOption) {
	i.processHandlerOption(http.MethodOptions, pattern, opts...)

	i.Router.Options(pattern, f)
}

// Connect adds the route `pattern` that matches a CONNECT http method to
// execute the `handlerFn` http.HandlerFunc.
func (i *Server) Connect(pattern string, f http.HandlerFunc, opts ...HandlerOption) {
	i.processHandlerOption(http.MethodConnect, pattern, opts...)

	i.Router.Connect(pattern, f)
}

// Head adds the route `pattern` that matches a HEAD http method to
// execute the `handlerFn` http.HandlerFunc.
func (i *Server) Head(pattern string, f http.HandlerFunc, opts ...HandlerOption) {
	i.processHandlerOption(http.MethodHead, pattern, opts...)

	i.Router.Head(pattern, f)
}

// Delete adds the route `pattern` that matches a DELETE http method to
// execute the `handlerFn` http.HandlerFunc.
func (i *Server) Delete(pattern string, f http.HandlerFunc, opts ...HandlerOption) {
	i.processHandlerOption(http.MethodDelete, pattern, opts...)

	i.Router.Delete(pattern, f)
}

// Trace adds the route `pattern` that matches a TRACE http method to
// execute the `handlerFn` http.HandlerFunc.
func (i *Server) Trace(pattern string, f http.HandlerFunc, opts ...HandlerOption) {
	i.processHandlerOption(http.MethodTrace, pattern, opts...)

	i.Router.Trace(pattern, f)
}

// Parse will parse the body from the reader and supports GZip compression.
// It will try to unmarshal the body from JSON.
func (i *Server) Parse(request *http.Request, data interface{}) error {
	return Parse(request, data)
}

func (i *Server) Response(w http.ResponseWriter, body interface{}, code int) error {
	return Response(w, body, code)
}
