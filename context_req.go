package invoke

import (
	"io"
	"net/http"
	"net/url"
)

// RemoteAddr returns the network address of the client sending the request.
func (ctx *HttpContext) RemoteAddr() string {
	return ctx.Req.RemoteAddr
}

// Method returns the HTTP request method.
func (ctx *HttpContext) Method() string {
	return ctx.Req.Method
}

// URL returns the URL of the request.
func (ctx *HttpContext) URL() *url.URL {
	return ctx.Req.URL
}

// ReqHeader returns the request header.
func (ctx *HttpContext) ReqHeader() http.Header {
	return ctx.Req.Header
}

// Body returns the request body.
func (ctx *HttpContext) Body() io.ReadCloser {
	return ctx.Req.Body
}

// ContentLength returns the length of the request body.
func (ctx *HttpContext) ContentLength() int64 {
	return ctx.Req.ContentLength
}

// Host returns the host name provided by the request.
func (ctx *HttpContext) Host() string {
	return ctx.Req.Host
}

// FormValue returns the first value for the named component of the query.
func (ctx *HttpContext) FormValue(key string) string {
	return ctx.Req.FormValue(key)
}

// PostFormValue returns the first value for the named component of the POST or PUT request body.
func (ctx *HttpContext) PostFormValue(key string) string {
	return ctx.Req.PostFormValue(key)
}

// UserAgent returns the user agent string provided in the request header.
func (ctx *HttpContext) UserAgent() string {
	return ctx.Req.UserAgent()
}
