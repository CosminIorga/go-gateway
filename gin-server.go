package main

import "github.com/gin-gonic/gin"

func main() {
	router := gin.Default()

	router.GET("/blocks/*path", func(c *gin.Context) {
		forwardPath := c.Param("path")

		c.JSON(200, gin.H{
			"path": forwardPath,
		})
	})

	router.Run() // listen and serve on 0.0.0.0:8080
}
