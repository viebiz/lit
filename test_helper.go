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
	ctx := lightningContext{
		Context: ginCtx,
	}

	return rtr, ctx, func() {
		route.HandleContext(ginCtx)
	}
}

func NewTestRoute(w http.ResponseWriter) (Router, func(func(Context))) {
	gin.SetMode(gin.TestMode)
	route := gin.New()
	rtr := router{
		ginRouter: route,
	}

	return rtr, func(cb func(Context)) {
		ginCtx := gin.CreateTestContextOnly(w, route)
		cb(lightningContext{
			Context: ginCtx,
		})

		route.HandleContext(ginCtx)
	}
}
