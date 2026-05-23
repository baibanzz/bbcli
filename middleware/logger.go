package middleware

import (
	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: 实现中间件逻辑
		c.Next()
	}
}
