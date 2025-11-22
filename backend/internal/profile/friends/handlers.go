package friends

import (
	"net/http"
	"strconv"

	"github.com/alexedwards/scs/v2"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"main/internal/pg"
)

type FriendRequestPayload struct {
	FriendID int64 `json:"friend_id"`
}

func SendFriendRequestHandler(c *gin.Context, sessionManager *scs.SessionManager) {
	senderID := sessionManager.GetInt64(c.Request.Context(), "userID")
	if senderID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var payload FriendRequestPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body, 'friend_id' is required"})
		return
	}

	err := pg.CreateFriendRequest(senderID, payload.FriendID)
	if err != nil {
		zap.S().Errorw("Failed to create friend request", "error", err, "sender_id", senderID, "receiver_id", payload.FriendID)
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Friend request sent successfully"})
}

func GetFriendsHandler(c *gin.Context, sessionManager *scs.SessionManager) {
	userID := sessionManager.GetInt64(c.Request.Context(), "userID")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	friends, err := pg.GetFriendsByUserID(userID)
	if err != nil {
		zap.S().Errorw("Failed to get friends list", "error", err, "user_id", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve friends list"})
		return
	}

	if friends == nil {
		friends = make([]pg.FriendUser, 0)
	}

	c.JSON(http.StatusOK, friends)
}

func DeleteFriendHandler(c *gin.Context, sessionManager *scs.SessionManager) {
	userID := sessionManager.GetInt64(c.Request.Context(), "userID")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	friendIDStr := c.Param("friend_id")
	friendID, err := strconv.ParseInt(friendIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid friend ID format"})
		return
	}

	err = pg.DeleteFriendship(userID, friendID)
	if err != nil {
		zap.S().Errorw("Failed to delete friendship", "error", err, "user_id", userID, "friend_id", friendID)
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// --- New Handlers ---

func GetIncomingRequestsHandler(c *gin.Context, sessionManager *scs.SessionManager) {
	receiverID := sessionManager.GetInt64(c.Request.Context(), "userID")
	if receiverID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	requests, err := pg.GetIncomingFriendRequests(receiverID)
	if err != nil {
		zap.S().Errorw("Failed to get incoming friend requests", "error", err, "user_id", receiverID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve incoming friend requests"})
		return
	}

	if requests == nil {
		requests = make([]pg.IncomingRequest, 0)
	}

	c.JSON(http.StatusOK, requests)
}

type UpdateRequestPayload struct {
	Status string `json:"status"`
}

func UpdateFriendRequestHandler(c *gin.Context, sessionManager *scs.SessionManager) {
	receiverID := sessionManager.GetInt64(c.Request.Context(), "userID")
	if receiverID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	requestIDStr := c.Param("request_id")
	requestID, err := strconv.ParseInt(requestIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID format"})
		return
	}

	var payload UpdateRequestPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body, 'status' is required"})
		return
	}

	err = pg.UpdateFriendRequestStatus(requestID, receiverID, payload.Status)
	if err != nil {
		zap.S().Errorw("Failed to update friend request", "error", err, "receiver_id", receiverID, "request_id", requestID)
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Friend request status updated successfully"})
}