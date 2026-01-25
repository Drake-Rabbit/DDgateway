package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
	})
}

func Error(c *gin.Context, message string) {
	c.JSON(http.StatusOK, gin.H{
		"success": false,
		"error":   message,
	})
}

func BadRequest(c *gin.Context, message string) {
	Error(c, message)
}

func Unauthorized(c *gin.Context, message string) {
	Error(c, message)
}

func Forbidden(c *gin.Context, message string) {
	Error(c, message)
}

func NotFound(c *gin.Context, message string) {
	Error(c, message)
}

func InternalError(c *gin.Context, message string) {
	Error(c, message)
}

// Setup sets up response middleware
func Setup(r *gin.Engine) {
	r.Use(func(c *gin.Context) {
		c.Next()
	})
}
