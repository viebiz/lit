package lit

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRouterForTest(w http.ResponseWriter) (Router, Context, func()) {
	gin.SetMode(gin.TestMode)
	route := gin.New()
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
