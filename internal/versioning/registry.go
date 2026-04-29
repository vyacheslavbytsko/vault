package versioning

import (
	"net/http"

	"vault/api/auth"
	"vault/api/chain"
	"vault/api/misc"
	"vault/internal/app"

	"github.com/gin-gonic/gin"
)

type MethodHandlers map[string]gin.HandlerFunc
type EndpointHandlers map[string]MethodHandlers

func NewHandlersByVersion(deps *app.Dependencies) map[string]EndpointHandlers {
	return map[string]EndpointHandlers{
		VersionV1dot0: {
			http.MethodGet: {
				EndpointPing:        misc.PingV1dot0,
				EndpointMe:          auth.MeV1dot0(deps),
				EndpointChain:       chain.ChainsV1dot0(deps),
				EndpointChainEvents: chain.FetchEventsV1dot0(deps),
			},
			http.MethodPost: {
				EndpointRegister:    auth.RegisterV1dot0(deps),
				EndpointLogin:       auth.LoginV1dot0(deps),
				EndpointRefresh:     auth.RefreshV1dot0(deps),
				EndpointChain:       chain.CreateChainV1dot0(deps),
				EndpointChainEvents: chain.AppendEventV1dot0(deps),
			},
		},
	}
}
