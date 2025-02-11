package lit

import (
	"context"
	"errors"
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

	JSONP(code int, obj any)

	ProtoBuf(code int, obj any)

	AbortWithStatus(code int)

	AbortWithError(err error)

	Next()
}

type lightningContext struct {
	*gin.Context
}

func (c lightningContext) Request() *http.Request {
	return c.Context.Request
}

func (c lightningContext) Writer() ResponseWriter {
	return c.Context.Writer
}

func (c lightningContext) SetRequest(r *http.Request) {
	c.Context.Request = r
}

func (c lightningContext) SetRequestContext(ctx context.Context) *http.Request {
	c.Context.Request = c.Context.Request.WithContext(ctx)
	return c.Context.Request
}

func (c lightningContext) SetWriter(w ResponseWriter) {
	c.Context.Writer = w
}

func (c lightningContext) AbortWithError(err error) {
	var httpErr HttpError
	if errors.As(err, &httpErr) {
		c.Context.AbortWithStatusJSON(httpErr.Status, httpErr)
		return
	}

	c.Context.AbortWithError(http.StatusInternalServerError, err)
}

func (c lightningContext) Bind(obj interface{}) error {
	if err := c.Context.Bind(obj); err != nil {
		return err
	}

	// Because Context.Bind does not binding URI by default
	if err := c.Context.BindUri(obj); err != nil {
		return err
	}

	return nil
}
