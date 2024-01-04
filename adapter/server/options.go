package server

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

type Option func(*Server)

// WithTimeout is a short-hand middleware to set a timeout value on the request context
// https://pkg.go.dev/github.com/go-chi/chi/middleware#Timeout
func WithTimeout(timeout time.Duration) Option {
	return func(i *Server) {
		i.Router.Use(
			middleware.Timeout(timeout),
		)
	}
}

// WithValue is a short-hand middleware to set a key/value on the request context
// this middleware utilized `github.com/go-chi/chi/middleware`
func WithValue(key string, value interface{}) Option {
	return func(i *Server) {
		i.Router.Use(
			middleware.WithValue(key, value),
		)
	}
}

// AllowCORS is a middleware to allow HTTP request from JS (web browser)
func WithCORS(AllowedOrigins []string, AllowedMethod []string, AllowedHeaders []string) Option {
	return func(i *Server) {
		corsHTTP := cors.New(cors.Options{
			AllowedOrigins: AllowedOrigins,
			AllowedMethods: AllowedMethod,
			AllowedHeaders: AllowedHeaders,
		})

		i.Router.Use(corsHTTP.Handler)
	}
}

func WithProfiler(path string) Option {
	return func(i *Server) {
		if !strings.HasPrefix(path, "/") {
			path = fmt.Sprintf("/%s", path)
		}

		i.profilerPath = path
	}
}

// WithMiddleware will add middleware to the chi http server we can then put our custom middleware or we can use middleware
// that being provided by go-chi.
// https://pkg.go.dev/github.com/go-chi/chi/middleware
func WithMiddleware(middlewares ...func(handler http.Handler) http.Handler) Option {
	return func(i *Server) {
		i.Router.Use(middlewares...)
	}
}

type HandlerOption func(*Handler)
