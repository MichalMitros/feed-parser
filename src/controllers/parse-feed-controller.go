package controllers

import "github.com/gin-gonic/gin"

func PostParseFeed(c *gin.Context) {
	c.IndentedJSON(200, gin.H{
		"message": "OK",
	})
}
