package lit

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRouterForTest(w http.ResponseWriter) (Router, Context, func()) {
	gin.SetMode(gin.TestMode)
	route := gin.New()
	route.ContextWithFallback = true
	rtr := router{
		ginRouter: route,
	}

	ginCtx := gin.CreateTestContextOnly(w, route)
	ctx := litContext{
		Context: ginCtx,
	}

	return rtr, ctx, func() {
		route.HandleContext(ginCtx)
	}
}

func CreateTestContext(w http.ResponseWriter) Context {
	gin.SetMode(gin.TestMode)
	route := gin.New()
	route.ContextWithFallback = true
	ginCtx := gin.CreateTestContextOnly(w, route)
	ctx := litContext{
		Context: ginCtx,
	}

	return ctx
}
