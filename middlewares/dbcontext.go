package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// DBContext set database context to the request.
func DBContext(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set("database", db)
		ctx.Next()
	}
}
