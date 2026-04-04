package handler

import "github.com/gin-gonic/gin"

func success(c *gin.Context, status int, data interface{}) {
	c.JSON(status, gin.H{"data": data})
}

func fail(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"error": message})
}
