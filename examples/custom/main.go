package main

import (
	"fmt"
	"github.com/coreyvan/vesper"
	"log"
	"net/http"
	"os"
)

const port = "4000"

type CustomContext struct {
	vesper.Context

	Logger *log.Logger
}

func reqFormatter(req *http.Request) string {
	return fmt.Sprintf("[%s] %s | %s", req.Method, req.RequestURI, req.RemoteAddr)
}

func main() {
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

	logger.Printf("server listening on port %s...\n", port)
	if err := srv.Serve(); err != nil {
		logger.Printf("received error: %v\n", err)
	}
}
