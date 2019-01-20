package main

import "github.com/gin-gonic/gin"

func authorizationMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
	}
}
