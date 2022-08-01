package utils

import "github.com/gin-gonic/gin"

func ApplyError(c *gin.Context, err error) {
	if err != nil {
		c.JSON(500, gin.H{
			"status": "failed",
			"error":  err.Error(),
		})
	}
}
