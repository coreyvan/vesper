package vesper

import (
	"fmt"
	"net"
	"net/http"
)

type ContextFn[T any] func(Context) T

type Server[T any] struct {
	host string
	port string
	mux  *http.ServeMux
	fn   ContextFn[T]
}

type ServerConfig struct {
	Host string
	Port string
}

func NewServer(cfg ServerConfig) *Server[Context] {
	return newServer(cfg, func(c Context) Context {
		return c
	})
}

func NewServerWithCustomContext[T any](cfg ServerConfig, fn ContextFn[T]) *Server[T] {
	return newServer(cfg, fn)
}

func newServer[T any](cfg ServerConfig, fn ContextFn[T]) *Server[T] {
	return &Server[T]{
		host: cfg.Host,
		port: cfg.Port,
		mux:  http.NewServeMux(),
		fn:   fn,
	}
}

func (s *Server[T]) Serve() error {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", s.host, s.port))
	if err != nil {
		return err
	}

	return http.Serve(lis, s.mux)
}

func (s *Server[T]) Handle(route string, handler Handler[T]) {
	s.mux.HandleFunc(route, func(w http.ResponseWriter, req *http.Request) {
		outCtx := s.fn(&context{
			w: w,
			r: req,
		})
		handler(outCtx)
	})
}
