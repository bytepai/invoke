package invoke

import (
	"bufio"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"path/filepath"
	"runtime"
	"strconv"
)

// HttpContext represents the HTTP context.
type HttpContext struct {
	W      http.ResponseWriter
	Req    *http.Request
	Params map[string]string
}

// ResponseResult represents a unified response structure.
type ResponseResult struct {
	Code int         `json:"code"`
	URL  string      `json:"url"`
	Desc string      `json:"desc"`
	Data interface{} `json:"data"`
}

// ErrorCode represents custom error codes.
type ErrorCode int

const (
	// AuthError indicates an error related to permission.
	AuthError ErrorCode = iota*1000 + 1000
	// ParamError indicates an error related to parameters.
	ParamError
	// BizError indicates a business logic error.
	BizError
	// NetError indicates a network-related error.
	NetError
	// DBError indicates a database-related error.
	DBError
	// IOError indicates an I/O-related error.
	IOError
	// OtherError indicates an error that does not fall into any specific category.
	OtherError
)

// ErrorCodeToString maps ErrorCode values to their corresponding strings.
var ErrorCodeToString = map[ErrorCode]string{
	AuthError:  "AuthError",
	ParamError: "ParamError",
	BizError:   "BizError",
	NetError:   "NetError",
	DBError:    "DBError",
	IOError:    "IOError",
	OtherError: "OtherError",
}

// ErrorCodeToString converts an ErrorCode value to its corresponding string.
func (ec ErrorCode) String() string {
	if str, ok := ErrorCodeToString[ec]; ok {
		return str
	}
	return "Unknown"
}

// getCallerInfo retrieves the filename and line number of the caller.
func getCallerInfo() string {
	_, file, line, ok := runtime.Caller(2) // Adjust the skip value as needed
	if !ok {
		return "unknown"
	}
	return filepath.Base(file) + ":" + strconv.Itoa(line)
}

// WriteErrorJSON writes an error message as JSON to the response with the specified status code.
func (ctx *HttpContext) WriteErrorJSON(statusCode interface{}, errMsg interface{}) {
	ctx.W.Header().Set("Content-Type", "application/json")
	ctx.W.WriteHeader(http.StatusOK)

	var code int
	var errMsgFormatted = errMsg
	// Determine the type of statusCode and handle accordingly
	switch v := statusCode.(type) {
	case int:
		code = int(v)
	case int64:
		code = int(64)
	case ErrorCode:
		code = int(v)
		if str, ok := errMsg.(string); ok {
			errMsgFormatted = fmt.Sprintf("%v: %v", v.String(), str)
		}
	default:
		code = int(OtherError)
	}

	response := ResponseResult{
		Code: code,
		URL:  ctx.Req.URL.Path,
		Desc: getCallerInfo(),
		Data: errMsgFormatted,
	}

	json.NewEncoder(ctx.W).Encode(response)
}

// WriteSuccessJSON writes an object as JSON to the response with a 200 status code.
func (ctx *HttpContext) WriteSuccessJSON(data interface{}) {
	ctx.W.Header().Set("Content-Type", "application/json")
	ctx.W.WriteHeader(http.StatusOK)

	response := ResponseResult{
		Code: http.StatusOK,
		URL:  ctx.Req.URL.Path,
		Desc: getCallerInfo(),
		Data: data,
	}

	json.NewEncoder(ctx.W).Encode(response)
}

// WriteErrorXML writes an error message as XML to the response with the specified status code.
func (ctx *HttpContext) WriteErrorXML(statusCode int, errMsg interface{}) {
	ctx.W.Header().Set("Content-Type", "application/xml")
	ctx.W.WriteHeader(statusCode)
	response := ResponseResult{
		Code: statusCode,
		URL:  ctx.Req.URL.Path,
		Desc: getCallerInfo(),
		Data: errMsg,
	}
	xml.NewEncoder(ctx.W).Encode(response)
}

// WriteSuccessXML writes an object as XML to the response with a 200 status code.
func (ctx *HttpContext) WriteSuccessXML(data interface{}) {
	ctx.W.Header().Set("Content-Type", "application/xml")
	ctx.W.WriteHeader(http.StatusOK)
	response := ResponseResult{
		Code: http.StatusOK,
		URL:  ctx.Req.URL.Path,
		Desc: getCallerInfo(),
		Data: data,
	}
	xml.NewEncoder(ctx.W).Encode(response)
}

// WriteString writes a string to the response.
func (ctx *HttpContext) WriteString(s string) {
	ctx.W.WriteHeader(http.StatusOK)
	ctx.W.Write([]byte(s))
}

// WriteByte writes byte to the client.
func (ctx *HttpContext) WriteByte(data []byte) {
	ctx.W.WriteHeader(http.StatusOK)
	ctx.W.Write(data)
}

// ParseJSONBody reads the JSON body from the request and unmarshals it into the target object.
func (ctx *HttpContext) ParseJSONBody(target interface{}) error {
	// Read the JSON body from the request
	body, err := ioutil.ReadAll(ctx.Req.Body)
	if err != nil {
		return fmt.Errorf("failed to read request body: %v", err)
	}

	// Unmarshal the JSON data into the target object
	if err := json.Unmarshal(body, target); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	return nil
}

// WriteHeader sends an HTTP response header with the provided status code.
func (ctx *HttpContext) WriteHeader(statusCode int) {
	ctx.W.WriteHeader(statusCode)
}

// Write writes the data to the connection as part of an HTTP reply.
func (ctx *HttpContext) Write(data []byte) (int, error) {
	return ctx.W.Write(data)
}

// Header returns the header map that will be sent by WriteHeader.
func (ctx *HttpContext) Header() http.Header {
	return ctx.W.Header()
}

// SetCookie adds a Set-Cookie header to the provided ResponseWriter's headers.
func (ctx *HttpContext) SetCookie(cookie *http.Cookie) {
	http.SetCookie(ctx.W, cookie)
}

// AddCookie adds a Set-Cookie header to the provided ResponseWriter's headers.
func (ctx *HttpContext) AddCookie(cookie *http.Cookie) {
	ctx.W.Header().Add("Set-Cookie", cookie.String())
}

// Flush sends any buffered data to the client.
func (ctx *HttpContext) Flush() {
	ctx.W.(http.Flusher).Flush()
}

// Hijack hijacks the connection.
func (ctx *HttpContext) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return ctx.W.(http.Hijacker).Hijack()
}

// CloseNotify returns a channel that receives a single value when the client connection has gone away.
func (ctx *HttpContext) CloseNotify() <-chan bool {
	return ctx.W.(http.CloseNotifier).CloseNotify()
}

// Pusher returns the HTTP/2 Pusher for the provided ResponseWriter.
func (ctx *HttpContext) Pusher() http.Pusher {
	return ctx.W.(http.Pusher)
}
