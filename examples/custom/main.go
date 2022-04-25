package main

import (
	"context"
	"fmt"
	"github.com/coreyvan/vesper"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const port = "4000"

type CustomContext struct {
	vesper.Context

	Logger *log.Logger
}

func reqFormatter(req *http.Request) string {
	return fmt.Sprintf("[%s] %s | %s", req.Method, req.RequestURI, req.RemoteAddr)
}

func run() error {
	logger := log.New(os.Stdout, "[server] ", 0)
	ctxFn := func(c vesper.Context) CustomContext {
		return CustomContext{
			Context: c,
			Logger:  logger,
		}
	}

	cfg := vesper.ServerConfig{
		Host: "localhost",
		Port: port,
	}

	srv := vesper.NewServerWithCustomContext(cfg, ctxFn)

	srv.UseMiddleware(
		vesper.ErrorHandler[CustomContext](logger),
		vesper.RequestLogger[CustomContext](reqFormatter, logger),
	)

	srv.Handle("/", func(c CustomContext) error {
		if _, err := c.ResponseWriter().Write([]byte("Hello World!")); err != nil {
			return fmt.Errorf("writing response: %w", err)
		}

		return nil
	})

	srv.Handle("/error", func(c CustomContext) error {
		w := c.ResponseWriter()

		c.Logger.Println("handled error")
		w.WriteHeader(500)
		w.Write([]byte("Internal Server Error"))

		return nil
	})

	srv.Handle("/unhandled", func(c CustomContext) error {
		return fmt.Errorf("unhandled error")
	})

	shutdown := make(chan os.Signal)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	serverErrors := make(chan error)

	go func() {
		logger.Printf("server listening on port %s...\n", port)
		serverErrors <- srv.Serve(shutdown)
	}()

	select {
	case sig := <-shutdown:
		logger.Printf("received os signal: %s\n", sig)
		defer logger.Println("shutdown complete")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			return srv.Close()
		}
		return nil
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)
	}
}

func main() {
	if err := run(); err != nil {
		// uh... why does this not exit the process?
		log.Fatal(err)
	}
}
