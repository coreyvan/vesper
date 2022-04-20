package main

import (
	"fmt"
	"github.com/coreyvan/vesper"
	"log"
)

const port = "4000"

func main() {
	srv := vesper.NewServer(vesper.ServerConfig{
		Host: "localhost",
		Port: port,
	})

	srv.Handle("/", func(c vesper.Context) error {
		c.ResponseWriter().WriteHeader(200)
		_, err := c.ResponseWriter().Write([]byte("Hello World!"))
		return err
	})

	srv.Handle("/error", func(c vesper.Context) error {
		return fmt.Errorf("oopsie")
	})

	log.Printf("server listening on port %s...", port)
	if err := srv.Serve(); err != nil {
		log.Printf("received error: %v", err)
	}
}
