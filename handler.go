package vesper

type Handler[T any] func(T) error
