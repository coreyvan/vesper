package vesper

import "net/http"

type Context interface {
	Request() *http.Request
	SetRequest(*http.Request)

	ResponseWriter() http.ResponseWriter
	SetResponseWriter(http.ResponseWriter)
}

type ctx struct {
	w http.ResponseWriter
	r *http.Request
}

func (c *ctx) Request() *http.Request {
	return c.r
}

func (c *ctx) SetRequest(r *http.Request) {
	c.r = r
}

func (c *ctx) ResponseWriter() http.ResponseWriter {
	return c.w
}

func (c *ctx) SetResponseWriter(w http.ResponseWriter) {
	c.w = w
}
