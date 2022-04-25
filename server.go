package vesper

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"syscall"
)

type ContextFn[T Context] func(Context) T

type Server[T Context] struct {
	host     string
	port     string
	mux      *http.ServeMux
	fn       ContextFn[T]
	mw       []Middleware[T]
	mu       *sync.RWMutex
	shutdown chan os.Signal
	srv      *http.Server
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
	srv := http.Server{
		Addr: fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		//	TODO: add server config like timeouts here
	}

	return &Server[T]{
		host:     cfg.Host,
		port:     cfg.Port,
		mux:      http.NewServeMux(),
		fn:       fn,
		mw:       []Middleware[T]{},
		mu:       &sync.RWMutex{},
		shutdown: make(chan os.Signal),
		srv:      &srv,
	}
}

func (s *Server[T]) Serve(shutdown chan os.Signal) error {
	if s.srv == nil {
		return fmt.Errorf("server HTTP server has not been specified")
	}

	s.shutdown = shutdown
	s.srv.Handler = s.mux

	return s.srv.ListenAndServe()
}

func (s *Server[T]) Handle(route string, handler Handler[T], mw ...Middleware[T]) {
	handler = wrapMiddleware(mw, handler)
	handler = wrapMiddleware(s.mw, handler)

	h := func(w http.ResponseWriter, req *http.Request) {
		outCtx := s.fn(&ctx{
			w: w,
			r: req,
		})

		if err := handler(outCtx); err != nil {
			log.Printf("error from handler... shutting down: %v", err)
			s.SignalShutdown()
		}
	}

	s.mux.HandleFunc(route, h)
}

func (s *Server[T]) UseMiddleware(middleware ...Middleware[T]) {
	s.mu.Lock()
	s.mw = append(s.mw, middleware...)
	s.mu.Unlock()
}

func (s *Server[T]) SignalShutdown() {
	s.shutdown <- syscall.SIGTERM
}

func (s *Server[T]) Close() error {
	return s.srv.Close()
}

func (s *Server[T]) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
