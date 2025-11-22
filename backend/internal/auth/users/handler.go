package users

import (
	"bytes"
	"encoding/hex"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"main/internal/auth/password"
	"main/internal/pg"
)

type RegisterRequest struct {
	Login    string `json:"login"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

type UserData struct {
	ID           int
	PasswordHash string
	Salt         string
}

func RegisterUser(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON provided"})
		return
	}

	if req.Login == "" || req.Password == "" || req.Email == "" {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Username, email, and password are required"},
		)
		println("asdad")
		return
	}

	var exists bool
	err := pg.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = $1 OR email = $2)", req.Login, req.Email).
		Scan(&exists)
	if err != nil {
		zap.S().Errorw("Failed to check if user exists", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if exists {
		zap.S().
			Warnw("Register: user or email already exists", "username", req.Login, "email", req.Email)
		c.JSON(
			http.StatusConflict,
			gin.H{"error": "User with this username or email already exists"},
		)
		return
	}

	salt, _ := password.GenerateSalt(32)
	hash := password.HashPassword(req.Password, salt)

	err = pg.InsertInDB(req.Login, req.Email, hash, salt)
	if err != nil {
		zap.S().Errorw("Register: failed to insert user", "error", err)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to create user"},
		)
		return
	}

	zap.S().Infow("Register: insertion successful!", "username", req.Login)
	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

func AuthorizeUser(c *gin.Context, sessionManager *scs.SessionManager) {
	var req AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON provided"})
		return
	}

	var userData UserData
	// Read hex strings from the DB
	err := pg.DB.QueryRow("SELECT id, password_hash, salt FROM users WHERE username = $1 OR email = $1", req.Login).
		Scan(&userData.ID, &userData.PasswordHash, &userData.Salt)
	if err != nil {
		zap.S().
			Warnw("Authorization: User not found", "login", req.Login, "error", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Decode hex strings back to bytes
	saltBytes, err := hex.DecodeString(userData.Salt)
	if err != nil {
		zap.S().Errorw("Failed to decode salt from hex", "error", err)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Internal server error"},
		)
		return
	}
	hashBytes, err := hex.DecodeString(userData.PasswordHash)
	if err != nil {
		zap.S().Errorw("Failed to decode hash from hex", "error", err)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Internal server error"},
		)
		return
	}

	// Compare the byte slices
	if bytes.Equal(password.HashPassword(req.Password, saltBytes), hashBytes) {
		err = sessionManager.RenewToken(c.Request.Context())
		if err != nil {
			zap.S().Errorw("Failed to renew session token", "error", err)
			c.JSON(
				http.StatusInternalServerError,
				gin.H{"error": "Session error"},
			)
			return
		}

		sessionManager.Put(c.Request.Context(), "userID", userData.ID)
		zap.S().Infow("Authorization: success!", "userID", userData.ID)
		c.JSON(http.StatusOK, gin.H{"message": "authorize success"})
	} else {
		zap.S().Warnw("Authorization: passwords do not match", "login", req.Login)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
	}
}

func GetAllUsers(c *gin.Context) {
	rows, err := pg.DB.Query("SELECT id, username FROM users ORDER BY id ASC")
	if err != nil {
		zap.S().Errorw("Failed to query users from database", "error", err)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Database query failed"},
		)
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Username); err != nil {
			zap.S().Errorw("Failed to scan user row", "error", err)
			c.JSON(
				http.StatusInternalServerError,
				gin.H{"error": "Failed to process database results"},
			)
			return
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		zap.S().Errorw("Error during rows iteration", "error", err)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Database iteration failed"},
		)
		return
	}

	c.JSON(http.StatusOK, users)
}

func LogoutUser(c *gin.Context, sessionManager *scs.SessionManager) {
	// Destroy the session
	err := sessionManager.Destroy(c.Request.Context())
	if err != nil {
		zap.S().Errorw("Failed to destroy session", "error", err)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to log out"},
		)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}
