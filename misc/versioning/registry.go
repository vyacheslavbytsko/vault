package versioning

import (
	"vault/api/auth"
	"vault/api/misc"

	"github.com/gin-gonic/gin"
)

type endpointHandlers map[string]gin.HandlerFunc

var handlersByVersion = map[string]endpointHandlers{
	VersionV1dot0: {
		EndpointPing:     misc.PingV1dot0,
		EndpointRegister: auth.RegisterV1dot0,
		EndpointLogin:    auth.LoginV1dot0,
	},
}
