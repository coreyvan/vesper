package vesper

import (
	"fmt"
	"net"
	"net/http"
	"sync"
)

type ContextFn[T Context] func(Context) T

type Server[T Context] struct {
	host string
	port string
	mux  *http.ServeMux
	fn   ContextFn[T]
	mw   []Middleware[T]
	mu   *sync.RWMutex
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

func NewServerWithCustomContext[T Context](cfg ServerConfig, fn ContextFn[T]) *Server[T] {
	return newServer(cfg, fn)
}

func newServer[T Context](cfg ServerConfig, fn ContextFn[T]) *Server[T] {
	return &Server[T]{
		host: cfg.Host,
		port: cfg.Port,
		mux:  http.NewServeMux(),
		fn:   fn,
		mw:   []Middleware[T]{},
		mu:   &sync.RWMutex{},
	}
}

func (s *Server[T]) Serve() error {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", s.host, s.port))
	if err != nil {
		return err
	}

	return http.Serve(lis, s.mux)
}

func (s *Server[T]) Handle(route string, handler Handler[T], mw ...Middleware[T]) {
	handler = wrapMiddleware(mw, handler)
	handler = wrapMiddleware(s.mw, handler)

	h := func(w http.ResponseWriter, req *http.Request) {
		outCtx := s.fn(&context{
			w: w,
			r: req,
		})

		if err := handler(outCtx); err != nil {
			// TODO: this shouldn't panic, it should gracefully shutdown the server
			panic(err)
		}
	}
	
	s.mux.HandleFunc(route, h)
}

func (s *Server[T]) UseMiddleware(middleware ...Middleware[T]) {
	s.mu.Lock()
	s.mw = append(s.mw, middleware...)
	s.mu.Unlock()
}
