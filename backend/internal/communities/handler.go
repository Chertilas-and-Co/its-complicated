package communities

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"main/internal/pg"
)

type CreateCommunityPayload struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IsPrivate   bool   `json:"is_private"`
}

func CreateCommunityHandler(
	c *gin.Context,
	sessionManager *scs.SessionManager,
) {
	creatorID := sessionManager.GetInt64(c.Request.Context(), "userID")
	if creatorID == 0 {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{"error": "User not authenticated"},
		)
		return
	}

	var payload CreateCommunityPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "Invalid request body. 'name' and 'description' are required.",
			},
		)
		return
	}

	if payload.Name == "" {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Community name cannot be empty."},
		)
		return
	}

	newID, err := pg.CreateCommunity(
		payload.Name,
		payload.Description,
		creatorID,
	)
	if err != nil {
		zap.S().
			Errorw("Failed to create community", "error", err, "creator_id", creatorID)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to create community in database."},
		)
		return
	}

	// Return the newly created community's ID, which the frontend needs
	c.JSON(http.StatusCreated, gin.H{
		"id":          newID,
		"name":        payload.Name,
		"description": payload.Description,
		"creator_id":  creatorID,
	})
}
