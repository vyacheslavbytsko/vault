package auth

import (
	"strings"
	"vault/internal/app"
	"vault/internal/database/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const maxSessionNameLength = 255

type TokenPair struct {
	AccessToken      string
	AccessExpiresAt  int64
	RefreshToken     string
	RefreshExpiresAt int64
}

func IssueTokenPairForNewSession(deps *app.Dependencies, user models.User, sessionName string) (TokenPair, error) {
	sessionName = normalizeSessionName(sessionName)

	refreshTokenID := uuid.NewString()
	session := models.Session{
		AccountID:      user.ID,
		Name:           sessionName,
		RefreshTokenID: refreshTokenID,
	}

	var tokens TokenPair
	err := deps.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&session).Error; err != nil {
			return err
		}

		pair, err := generateTokenPairForSession(deps, user, session.ID, refreshTokenID)
		if err != nil {
			return err
		}

		tokens = pair
		return nil
	})
	if err != nil {
		return TokenPair{}, err
	}

	return tokens, nil
}

func IssueTokenPairForExistingSession(deps *app.Dependencies, user models.User, session models.Session, expectedRefreshTokenID string) (TokenPair, error) {
	refreshTokenID := uuid.NewString()

	var tokens TokenPair
	err := deps.DB.Transaction(func(tx *gorm.DB) error {
		pair, err := generateTokenPairForSession(deps, user, session.ID, refreshTokenID)
		if err != nil {
			return err
		}

		updateResult := tx.Model(&models.Session{}).
			Where("id = ? AND account_id = ? AND refresh_token_id = ?", session.ID, user.ID, expectedRefreshTokenID).
			Update("refresh_token_id", refreshTokenID)
		if updateResult.Error != nil {
			return updateResult.Error
		}

		if updateResult.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		tokens = pair
		return nil
	})
	if err != nil {
		return TokenPair{}, err
	}

	return tokens, nil
}

func generateTokenPairForSession(deps *app.Dependencies, user models.User, sessionID string, refreshTokenID string) (TokenPair, error) {
	accessToken, accessExpiresAt, err := deps.AccessJWTManager.GenerateToken(sessionID, user.ID, refreshTokenID)
	if err != nil {
		return TokenPair{}, err
	}

	refreshToken, refreshExpiresAt, err := deps.RefreshJWTManager.GenerateToken(sessionID, user.ID, refreshTokenID)
	if err != nil {
		return TokenPair{}, err
	}

	return TokenPair{
		AccessToken:      accessToken,
		AccessExpiresAt:  accessExpiresAt,
		RefreshToken:     refreshToken,
		RefreshExpiresAt: refreshExpiresAt,
	}, nil
}

func normalizeSessionName(sessionName string) string {
	sessionName = strings.TrimSpace(sessionName)
	if sessionName == "" {
		return "unknown"
	}

	runes := []rune(sessionName)
	if len(runes) > maxSessionNameLength {
		return string(runes[:maxSessionNameLength])
	}

	return sessionName
}
