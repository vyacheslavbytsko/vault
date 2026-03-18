package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterV1dot0(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"message": "register endpoint not implemented yet",
	})
}
