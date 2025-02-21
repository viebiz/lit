package lit

import (
	"context"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Context interface {
	context.Context

	Request() *http.Request

	Writer() ResponseWriter

	SetRequest(*http.Request)

	SetRequestContext(ctx context.Context) *http.Request

	SetWriter(w ResponseWriter)

	Bind(obj interface{}) error

	FormFile(name string) (*multipart.FileHeader, error)

	MultipartForm() (*multipart.Form, error)

	Set(key string, value any)

	Get(key string) (value any, exists bool)

	Status(code int)

	Header(key string, value string)

	JSON(code int, obj any)

	ProtoBuf(code int, obj any)

	AbortWithError(err error)

	Next()
}

type litContext struct {
	*gin.Context
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

func (c litContext) Bind(obj interface{}) error {
	if err := c.Context.Bind(obj); err != nil {
		return err
	}

	// Because Context.Bind does not binding URI by default
	if err := c.Context.BindUri(obj); err != nil {
		return err
	}

	return nil
}

func (c litContext) AbortWithError(err error) {
	respondJSON(c.Request().Context(), c.Writer(), err)
}
