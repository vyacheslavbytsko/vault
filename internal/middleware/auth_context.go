package middleware

import "github.com/gin-gonic/gin"

const (
	ContextAuthClaimsKey         = "auth.claims"
	ContextAuthSessionIDKey      = "auth.session_id"
	ContextAuthUserIDKey         = "auth.user_id"
	ContextAuthRefreshTokenIDKey = "auth.refresh_token_id"
)

func GetCurrentUser(c *gin.Context) (string, bool) {
	userIDValue, userIDExists := c.Get(ContextAuthUserIDKey)
	if !userIDExists {
		return "", false
	}

	userID, userIDOk := userIDValue.(string)
	if !userIDOk || userID == "" {
		return "", false
	}

	return userID, true
}

func GetCurrentSession(c *gin.Context) (string, string, bool) {
	sessionIDValue, sessionIDExists := c.Get(ContextAuthSessionIDKey)
	refreshTokenIDValue, refreshTokenIDExists := c.Get(ContextAuthRefreshTokenIDKey)
	if !sessionIDExists || !refreshTokenIDExists {
		return "", "", false
	}

	sessionID, sessionIDOk := sessionIDValue.(string)
	refreshTokenID, refreshTokenIDOk := refreshTokenIDValue.(string)
	if !sessionIDOk || !refreshTokenIDOk || sessionID == "" || refreshTokenID == "" {
		return "", "", false
	}

	return sessionID, refreshTokenID, true
}
