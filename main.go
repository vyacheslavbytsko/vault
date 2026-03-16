package main

import "github.com/gin-gonic/gin"
import "net/http"

func main() {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	err := r.Run("127.0.0.1:8081")
	if err != nil {
		return
	}
}
