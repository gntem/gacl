package middlewares

import (
	"github.com/gin-gonic/gin"
)

// AuthorizationMiddleware handle authorize of requests.
func AuthorizationMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
	}
}
