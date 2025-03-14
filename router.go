package lit

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Router interface {
	Use(middleware ...func(ctx Context))

	Handle(method string, relativePath string, handler HandlerFunc)

	HandleWithErr(method string, relativePath string, handler ErrHandlerFunc)

	Get(relativePath string, handler ErrHandlerFunc)

	Post(relativePath string, handler ErrHandlerFunc)

	Put(relativePath string, handler ErrHandlerFunc)

	Patch(relativePath string, handler ErrHandlerFunc)

	Delete(relativePath string, handler ErrHandlerFunc)

	Group(relativePath string, routerFunc func(Router))
}

type router struct {
	ginRouter gin.IRouter
}

func NewRouter() (Router, http.Handler) {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.ContextWithFallback = true

	return router{
		ginRouter: engine,
	}, engine.Handler()
}

func (rtr router) Use(middleware ...func(ctx Context)) {
	handlers := make([]gin.HandlerFunc, len(middleware))
	for idx, m := range middleware {
		handlers[idx] = func(ctx *gin.Context) {
			m(litContext{Context: ctx})
		}
	}

	rtr.ginRouter.Use(handlers...)
}

func (rtr router) Handle(method string, relativePath string, handler HandlerFunc) {
	rtr.ginRouter.Handle(method, relativePath, func(ctx *gin.Context) {
		handler(litContext{Context: ctx})
	})
}

func (rtr router) HandleWithErr(method string, relativePath string, handler ErrHandlerFunc) {
	rtr.ginRouter.Handle(method, relativePath, wrapErrHandler(handler))
}

func (rtr router) Get(relativePath string, handler ErrHandlerFunc) {
	rtr.HandleWithErr(http.MethodGet, relativePath, handler)
}

func (rtr router) Post(relativePath string, handler ErrHandlerFunc) {
	rtr.HandleWithErr(http.MethodPost, relativePath, handler)
}

func (rtr router) Put(relativePath string, handler ErrHandlerFunc) {
	rtr.HandleWithErr(http.MethodPut, relativePath, handler)
}

func (rtr router) Patch(relativePath string, handler ErrHandlerFunc) {
	rtr.HandleWithErr(http.MethodPatch, relativePath, handler)
}

func (rtr router) Delete(relativePath string, handler ErrHandlerFunc) {
	rtr.HandleWithErr(http.MethodDelete, relativePath, handler)
}

func (rtr router) Group(relativePath string, routerFunc func(Router)) {
	routerGroup := rtr.ginRouter.Group(relativePath)
	wrappedRoute := router{
		ginRouter: routerGroup,
	}

	routerFunc(wrappedRoute)
}
