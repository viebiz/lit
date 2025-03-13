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
	// Read more at https://pkg.go.dev/github.com/gin-gonic/gin#Context.Set
	Set(key string, value any)

	// Get returns the value for the given key
	// Eead more at https://pkg.go.dev/github.com/gin-gonic/gin#Context.Get
	Get(key string) (value any, exists bool)

	// Status sets the HTTP response code
	Status(code int)

	// Header writes header to the response
	// If value is empty, it will remove the key in the header
	Header(key string, value string)

	// JSON serializes the given struct as JSON into the response body
	JSON(code int, obj any)

	// ProtoBuf serializes the given protocol buffer object and writes it to the response
	ProtoBuf(code int, obj any)

	// AbortWithError will abort the remaining handlers and return an error to the client
	// Should be used in middleware
	AbortWithError(err error)

	// Next continues to the next handler in the chain
	Next()
}

type litContext struct {
	*gin.Context
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

func (c litContext) AbortWithError(obj error) {
	// Set JSON header
	c.Header("Content-Type", "application/json")

	var (
		status        = http.StatusInternalServerError
		errBody error = ErrDefaultInternal
	)

	// Determine the error type and choose a response payload
	var litErr Error
	if errors.As(obj, &litErr) {
		sc := litErr.StatusCode()
		if sc < http.StatusInternalServerError || sc == http.StatusServiceUnavailable {
			status, errBody = sc, litErr
		}
	}

	// Marshal the response payload
	respBytes, err := json.Marshal(errBody)
	if err != nil {
		monitoring.FromContext(c).Errorf(err, "[AbortWithError] JSON marshal failed, using ErrDefaultInternal")
		if respBytes, err = json.Marshal(ErrDefaultInternal); err != nil {
			monitoring.FromContext(c).Errorf(err, "[AbortWithError] JSON marshal of ErrDefaultInternal failed") // Should never happen
			respBytes = []byte(`{}`)
		}

		status = http.StatusInternalServerError
	}

	// Write response
	c.AbortWithStatus(status) // Abort and write header now
	if _, writeErr := c.Writer().Write(respBytes); writeErr != nil {
		monitoring.FromContext(c).Errorf(writeErr, "[AbortWithError] Write failed")
	}
}
