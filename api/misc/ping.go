package misc

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func PingV1dot0(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
