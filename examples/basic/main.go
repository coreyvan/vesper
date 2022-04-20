package main

import (
	"github.com/coreyvan/vesper"
	"log"
)

const port = "4000"

func main() {
	srv := vesper.NewServer(vesper.ServerConfig{
		Host: "localhost",
		Port: port,
	})

	srv.Handle("/", func(c vesper.Context) {
		c.ResponseWriter().WriteHeader(200)
		c.ResponseWriter().Write([]byte("Hello World!"))
	})

	log.Printf("server listening on port %s...", port)
	if err := srv.Serve(); err != nil {
		log.Printf("received error: %v", err)
	}
}
