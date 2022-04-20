package main

import (
	"github.com/coreyvan/vesper"
	"log"
	"os"
)

const port = "4000"

type CustomContext struct {
	vesper.Context

	Logger *log.Logger
}

func main() {
	ctxFn := func(c vesper.Context) CustomContext {
		return CustomContext{
			Context: c,
			Logger:  log.New(os.Stdout, "[server] ", 0),
		}
	}

	cfg := vesper.ServerConfig{
		Host: "localhost",
		Port: port,
	}

	srv := vesper.NewServerWithCustomContext(cfg, ctxFn)

	srv.Handle("/", func(c CustomContext) {
		c.Logger.Printf("received request: %v", c.Request())

		c.ResponseWriter().WriteHeader(200)
		c.ResponseWriter().Write([]byte("Hello World!"))
	})

	log.Printf("server listening on port %s...", port)
	if err := srv.Serve(); err != nil {
		log.Printf("received error: %v", err)
	}
}
