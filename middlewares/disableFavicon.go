package middlewares

import (
	"github.com/gin-gonic/gin"
)

// DisableFavicon ignore favicon requests.
func DisableFavicon() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.URL.Path == "/favicon.ico" {
			return
		}
	}
}
