/*
 * Copyright (c) 2019. Octofox.io
 */

package http

import (
	"context"
	"github.com/gorilla/mux"
	_ "github.com/gorilla/mux"
	"github.com/rs/xid"
	"net/http"
)

const (
	HeaderRequestIDKey     = "RequestID"
	HeaderAuthorizationKey = "Authorization"

	RequestContextKey        = "request"
	ResponseWriterContextKey = "response"
	InputBodyContextKey      = "input"
)

type Server struct {
	r *mux.Router
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.r.ServeHTTP(w, req)
}

type Handler func(ctx context.Context) (interface{}, error)
type EndpointFunc func(ctx context.Context, request interface{}) (interface{}, error)
type Middleware func(next Handler) Handler

func EndpointHandler(endpoint EndpointFunc) Handler {
	return func(ctx context.Context) (i interface{}, e error) {
		var request = ctx.Value(InputBodyContextKey)
		return endpoint(ctx, request)
	}
}

func (s *Server) registerHTTPHandler(method string, path string, handler Handler, middleware ...Middleware) {
	s.r.Methods(method).Path(path).
		HandlerFunc(execute(method, path, handler, middleware...))
}
func (s *Server) Get(path string, handler Handler, middleware ...Middleware) {
	s.registerHTTPHandler("Get", path, handler, middleware...)
}
func (s *Server) Post(path string, handler Handler, middleware ...Middleware) {
	s.registerHTTPHandler("Post", path, handler, middleware...)
}
func (s *Server) Put(path string, handler Handler, middleware ...Middleware) {
	s.registerHTTPHandler("Put", path, handler, middleware...)
}
func (s *Server) Delete(path string, handler Handler, middleware ...Middleware) {
	s.registerHTTPHandler("Delete", path, handler, middleware...)
}

func NewServer() *Server {
	r := mux.NewRouter()

	// Mess with request ID and recovery
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var requestID string
			if id := w.Header().Get("x-amzn-RequestId"); requestID != "" {
				requestID = id
			} else {
				requestID = xid.New().String()
			}
			w.Header().Set(HeaderRequestIDKey, requestID)
			next.ServeHTTP(w, r)
		})
	})
	return &Server{
		r: r,
	}
}
