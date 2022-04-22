package vesper

import (
	"log"
	"net/http"
)

type Middleware[T Context] func(Handler[T]) Handler[T]

func wrapMiddleware[T Context](mw []Middleware[T], handler Handler[T]) Handler[T] {
	for i := len(mw) - 1; i >= 0; i-- {
		h := mw[i]
		if h != nil {
			handler = h(handler)
		}
	}

	return handler
}

func ErrorHandler[T Context](logger *log.Logger) Middleware[T] {
	return func(next Handler[T]) Handler[T] {
		return func(ctx T) error {
			if err := next(ctx); err != nil {
				logger.Printf("error handled: %v", err)
			}
			return nil
		}
	}
}

type RequestFormatter func(r *http.Request) string

func RequestLogger[T Context](formatter RequestFormatter, logger *log.Logger) Middleware[T] {
	return func(next Handler[T]) Handler[T] {
		return func(ctx T) error {
			logger.Println(formatter(ctx.Request()))
			return next(ctx)
		}
	}
}
