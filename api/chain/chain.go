package chain

import (
	"errors"
	"net/http"

	"vault/internal/app"
	"vault/internal/database/models"
	"vault/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type createChainRequest struct {
	Name string `json:"name" binding:"required,min=1,max=128"`
}

func CreateChainV1dot0(deps *app.Dependencies) gin.HandlerFunc {
	return func(c *gin.Context) {
		if deps == nil || deps.DB == nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "database is not configured",
			})
			return
		}

		userID, ok := middleware.GetCurrentUser(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "unauthorized",
			})
			return
		}

		ownerID, err := uuid.Parse(userID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "unauthorized",
			})
			return
		}

		var request createChainRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "invalid request body",
			})
			return
		}

		chain := models.Chain{
			Name:    request.Name,
			OwnerID: ownerID,
		}

		if err := chain.Validate(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		if err := deps.DB.Create(&chain).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				c.JSON(http.StatusConflict, gin.H{
					"message": "you already have a chain with this name",
				})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "failed to create chain",
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"chain": chain.Name,
		})
	}
}

func ChainsV1dot0(deps *app.Dependencies) gin.HandlerFunc {
	return func(c *gin.Context) {
		if deps == nil || deps.DB == nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "database is not configured",
			})
			return
		}

		userID, ok := middleware.GetCurrentUser(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "unauthorized",
			})
			return
		}

		ownerID, err := uuid.Parse(userID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "unauthorized",
			})
			return
		}

		chains := make([]models.Chain, 0)
		if err := deps.DB.Where("owner_id = ?", ownerID).Order("created_at desc").Find(&chains).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "failed to load chains",
			})
			return
		}

		names := make([]string, len(chains))
		for i, chain := range chains {
			names[i] = chain.Name
		}

		c.JSON(http.StatusOK, gin.H{
			"chains": names,
		})
	}
}
