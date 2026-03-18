package versioning

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RouteByVersion(handlersByVersion map[string]EndpointHandlers, method string, endpoint string) gin.HandlerFunc {
	return func(c *gin.Context) {
		version := c.GetHeader(HeaderAPIVersion)
		if version == "" {
			version = DefaultAPIVersion
		}

		handlers, ok := handlersByVersion[version]
		if !ok {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": fmt.Sprintf("unsupported API version: %s", version),
			})
			return
		}

		methodHandlers, methodOk := handlers[method]
		if !methodOk {
			c.AbortWithStatusJSON(http.StatusMethodNotAllowed, gin.H{
				"message": fmt.Sprintf("method %s is not available for version %s", method, version),
			})
			return
		}

		handler, ok := methodHandlers[endpoint]
		if !ok {
			for _, byEndpoint := range handlers {
				if _, exists := byEndpoint[endpoint]; exists {
					c.AbortWithStatusJSON(http.StatusMethodNotAllowed, gin.H{
						"message": fmt.Sprintf("method %s is not allowed for endpoint %s", method, endpoint),
					})
					return
				}
			}

			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"message": fmt.Sprintf("endpoint %s is not available for version %s", endpoint, version),
			})
			return
		}

		handler(c)
	}
}
