package handlers

import (
	"its-complicated/internal/storage"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserHandler handles user-related HTTP requests.
type UserHandler struct {
	storage *storage.UserStorage
}

// NewUserHandler creates a new UserHandler.
func NewUserHandler(storage *storage.UserStorage) *UserHandler {
	return &UserHandler{storage: storage}
}

// RegisterUser is the handler for the user registration endpoint.
func (h *UserHandler) RegisterUser(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if req.Username == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username and password are required"})
		return
	}

	user, err := h.storage.CreateUser(req.Username, req.Password)
	if err != nil {
		if err == storage.ErrUserExists {
			c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user_id": user.ID,
	})
}
