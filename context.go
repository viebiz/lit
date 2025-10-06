package lit

import (
	"context"
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/viebiz/lit/monitoring"
)

const (
	plainContentType = "text/plain; charset=utf-8"
	jsonContentType  = "application/json"
)

type Context interface {
	context.Context

	// Request returns the underlying http.Request object
	Request() *http.Request

	// Writer returns the underlying ResponseWriter object
	Writer() ResponseWriter

	// SetRequest updates the current request
	SetRequest(*http.Request)

	// SetRequestContext updates the current request context
	SetRequestContext(ctx context.Context) *http.Request

	// SetWriter updates the current writer
	SetWriter(w ResponseWriter)

	// Bind binds the incoming request body and URI parameters to the provided object
	// Return error if the got error when binding and validating the object
	// Support validation tags from https://github.com/go-playground/validator/v10
	Bind(obj interface{}) error

	// FormFile returns the first file for the provided form key
	FormFile(name string) (*multipart.FileHeader, error)

	// MultipartForm returns the parsed multipart form, including file uploads
	// Default memory maxMemory = 32MB
	MultipartForm() (*multipart.Form, error)

	// Set stores a new key/value pair exclusively for this context
	Set(key string, value any)

	// Get returns the value for the given key
	Get(key string) (value any, exists bool)

	// Status sets the HTTP response code
	Status(code int)

	// Header writes header to the response
	// If value is empty, it will remove the key in the header
	Header(key string, value string)

	// String sends a string response with status code.
	String(code int, str string) error

	// JSON serializes the given struct as JSON into the response body
	JSON(code int, obj any) error

	// NoContent sends a no content response with status code
	NoContent(status int) error

	// FullPath returns the full path of the request
	FullPath() string

	// Param gets the URL path parameter value by key
	Param(key string) string

	// ParamWithDefault behaves like Param but returns the defaultValue if the parameter is missing or empty
	ParamWithDefault(key string, defaultValue string) string

	// ParamWithCallback behaves like Param but returns the value from the callback if the parameter is missing or empty
	ParamWithCallback(key string, callback func() string) string

	// Query gets the URL query parameter value by key
	Query(key string) string

	// QueryWithDefault behaves like Query but returns the defaultValue if the parameter is missing or empty
	QueryWithDefault(key string, defaultValue string) string

	// QueryWithCallback behaves like Query but returns the value from the callback if the parameter is missing or empty
	QueryWithCallback(key string, callback func() string) string

	// Next continues to the next handler in the chain
	Next()

	// Error handles request error
	// Should be called in middleware to render error response
	Error(err error)

	// Abort prevents pending handlers from being called
	Abort()
}

type litContext struct {
	*gin.Context
}

func newContext(c *gin.Context) Context {
	return &litContext{
		Context: c,
	}
}

func (c litContext) Value(key any) any {
	return c.Context.Value(key)
}

func (c litContext) Request() *http.Request {
	return c.Context.Request
}

func (c litContext) Writer() ResponseWriter {
	return c.Context.Writer
}

func (c litContext) SetRequest(r *http.Request) {
	c.Context.Request = r
}

func (c litContext) SetRequestContext(ctx context.Context) *http.Request {
	c.Context.Request = c.Context.Request.WithContext(ctx)
	return c.Context.Request
}

func (c litContext) SetWriter(w ResponseWriter) {
	c.Context.Writer = w
}

func (c litContext) ParamWithDefault(key string, defaultValue string) string {
	value := c.Context.Param(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func (c litContext) ParamWithCallback(key string, callback func() string) string {
	value := c.Context.Param(key)
	if value == "" {
		return callback()
	}
	return value
}

func (c litContext) QueryWithDefault(key string, defaultValue string) string {
	value := c.Context.Query(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func (c litContext) QueryWithCallback(key string, callback func() string) string {
	value := c.Context.Query(key)
	if value == "" {
		return callback()
	}
	return value
}

func (c litContext) String(status int, s string) error {
	c.Writer().Header().Set("Content-Type", plainContentType)
	c.Writer().WriteHeader(status)
	c.Writer().WriteHeaderNow()

	if _, err := c.Writer().Write([]byte(s)); err != nil {
		return err
	}
	return nil
}

func (c litContext) JSON(status int, obj any) error {
	c.Writer().Header().Set("Content-Type", jsonContentType)
	c.Writer().WriteHeader(status)
	c.Writer().WriteHeaderNow()

	return json.NewEncoder(c.Writer()).Encode(obj)
}

func (c litContext) NoContent(status int) error {
	c.Context.Writer.WriteHeader(status)
	return nil
}

func (c litContext) Get(key string) (value any, exists bool) {
	return c.Context.Get(key)
}

func (c litContext) Set(key string, value any) {
	c.Context.Set(key, value)
}

func (c litContext) Error(err error) {
	// Gin parse error and add to Context.Errors
	// But we don't use this way, just response this error
	// c.Context.Error(err)

	// If given error not in type HTTPError, prepare an internal error
	var httpErr Error
	if !errors.As(err, &httpErr) {
		httpErr = &HTTPError{
			Status: http.StatusInternalServerError,
			Code:   http.StatusText(http.StatusInternalServerError),
			Desc:   "Internal Server Error",
		}

		monitoring.FromContext(c.Request().Context()).Errorf(err, "Server error")
	}

	switch c.Request().Method {
	case http.MethodHead:
		err = c.NoContent(httpErr.StatusCode())
	default:
		err = c.JSON(httpErr.StatusCode(), httpErr)
	}
	if err != nil {
		monitoring.
			FromContext(c.Request().Context()).
			Errorf(err, "Failed to write response")
	}
}
