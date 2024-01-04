package server

//go:generate mockgen -destination=mock/server.go -package=mock . IServer

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
)

// Server contains information about this module.
type Server struct {
	Context context.Context
	Config  *Configuration

	Router *chi.Mux
	server *http.Server

	skipHeaders []string

	swaggerPattern string

	prometheusPath string
	profilerPath   string
}

// IServer defines the functions of this module.
type IServer interface {
	Close() error
	Serve() error

	Configure(options ...Option)
	GetRouter() *chi.Mux

	Get(pattern string, f http.HandlerFunc, opts ...HandlerOption)
	Post(pattern string, f http.HandlerFunc, opts ...HandlerOption)
	Put(pattern string, f http.HandlerFunc, opts ...HandlerOption)
	Patch(pattern string, f http.HandlerFunc, opts ...HandlerOption)
	Options(pattern string, f http.HandlerFunc, opts ...HandlerOption)
	Connect(pattern string, f http.HandlerFunc, opts ...HandlerOption)
	Head(pattern string, f http.HandlerFunc, opts ...HandlerOption)
	Delete(pattern string, f http.HandlerFunc, opts ...HandlerOption)
	Trace(pattern string, f http.HandlerFunc, opts ...HandlerOption)

	Parse(request *http.Request, data interface{}) error
	Response(w http.ResponseWriter, body interface{}, code int) error
}

// TracingData is a struct to define data for tracing
type TracingData struct {
	appName   string
	urlPath   string
	requestID string
}

type Handler struct {
	disableLogging bool
}
