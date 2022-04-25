package vesper

import (
	"go.uber.org/zap"
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

func ErrorHandler[T Context](logger *zap.Logger) Middleware[T] {
	return func(next Handler[T]) Handler[T] {
		return func(ctx T) error {
			if err := next(ctx); err != nil {
				logger.Sugar().Errorw("error from handler", "error", err)
				ctx.ResponseWriter().WriteHeader(500)
				ctx.ResponseWriter().Write([]byte("Internal server error"))
			}
			return nil
		}
	}
}

type RequestFormatter func(r *http.Request) string

func RequestLogger[T Context](formatter RequestFormatter, logger *zap.Logger) Middleware[T] {
	return func(next Handler[T]) Handler[T] {
		return func(ctx T) error {
			logger.Info(formatter(ctx.Request()))
			return next(ctx)
		}
	}
}
