package main

import (
	"context"
	"fmt"
	"github.com/coreyvan/vesper"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const port = "4000"

func run(logger *zap.Logger) error {
	srv := vesper.NewServer(vesper.ServerConfig{
		Host:   "localhost",
		Port:   port,
		Logger: logger,
	})

	srv.Handle("/", func(c vesper.Context) error {
		c.ResponseWriter().WriteHeader(200)
		_, err := c.ResponseWriter().Write([]byte("Hello World!"))
		return err
	})

	srv.Handle("/error", func(c vesper.Context) error {
		return fmt.Errorf("oopsie")
	})

	shutdown := make(chan os.Signal)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	serverErrors := make(chan error)

	go func() {
		logger.Sugar().Infof("server listening on port %s...\n", port)
		serverErrors <- srv.Serve(shutdown)
	}()

	select {
	case sig := <-shutdown:
		logger.Sugar().Infof("received os signal: %s\n", sig)
		defer logger.Info("shutdown complete")

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
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	if err := run(logger); err != nil {
		logger.Fatal(err.Error())
	}
}
