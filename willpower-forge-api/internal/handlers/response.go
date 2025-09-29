package handlers

import "github.com/gin-gonic/gin"

func respondSuccess(c *gin.Context, statusCode int, message string, data interface{}) {
	if message == "" {
		message = "Success"
	}

	payload := gin.H{
		"code":    0,
		"message": message,
	}

	if data != nil {
		payload["data"] = data
	}

	c.JSON(statusCode, payload)
}

func respondError(c *gin.Context, statusCode int, code int, message string) {
	c.JSON(statusCode, gin.H{
		"code":    code,
		"message": message,
	})
}
