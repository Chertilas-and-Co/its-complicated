package users

import (
	"bytes"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"go.uber.org/zap"

	passwd "main/internal/auth/password"
	"main/internal/pg"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

func RegisterUser(c *gin.Context) {
	username := c.PostForm("login")
	email := c.PostForm("email") // Get email from form
	password := c.PostForm("password")
	passwordConfirm := c.PostForm("passwordConfirm")

	if username == "" || password == "" || email == "" {
		c.String(http.StatusBadRequest, "Login, email, and password are required")
		return
	}

	if password != passwordConfirm {
		zap.S().Warn("Register: passwords do not match")
		c.String(http.StatusBadRequest, "Passwords do not match")
		return
	}

	// Check if user or email already exists
	var exists bool
	err := pg.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = $1 OR email = $2)", username, email).
		Scan(&exists)
	if err != nil {
		zap.S().Errorw("Failed to check if user exists", "error", err)
		c.String(http.StatusInternalServerError, "Database error")
		return
	}
	if exists {
		zap.S().Warnw("Register: user or email already exists", "username", username, "email", email)
		c.String(http.StatusConflict, "User with this username or email already exists")
		return
	}

	salt, _ := passwd.GenerateSalt(32)
	hash := passwd.HashPassword(password, salt)

	// Call the updated InsertInDB function
	err = pg.InsertInDB(username, email, hash, salt)
	if err != nil {
		zap.S().Errorw("Register: failed to insert user", "error", err)
		c.String(http.StatusInternalServerError, "Failed to create user")
		return
	}

	zap.S().Infow("Register: insertion successful!", "username", username)
	c.String(http.StatusCreated, "User created successfully")
}

type AuthRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserData struct {
	ID           int
	PasswordHash []byte
	Salt         []byte
}

func AuthorizeUser(c *gin.Context, sessionManager *scs.SessionManager) {
	var req AuthRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, "Invalid JSON")
		return
	}
	username := req.Login
	password := req.Password

	var userData UserData
	err := pg.DB.QueryRow("SELECT id, password_hash, salt FROM users WHERE username = $1", username).
		Scan(&userData.ID, &userData.PasswordHash, &userData.Salt)
	if err != nil {
		zap.S().
			Warnw("Authorization: User not found", "username", username, "error", err)
		c.String(http.StatusUnauthorized, "Invalid credentials")
		return
	}

	if bytes.Equal(
		passwd.HashPassword(password, userData.Salt),
		userData.PasswordHash,
	) {
		// Renew token to prevent session fixation
		err = sessionManager.RenewToken(c.Request.Context())
		if err != nil {
			zap.S().Errorw("Failed to renew session token", "error", err)
			c.String(http.StatusInternalServerError, "Session error")
			return
		}

		// Put user ID in session
		sessionManager.Put(c.Request.Context(), "userID", userData.ID)
		zap.S().Infow("Authorization: success!", "userID", userData.ID)
		c.String(http.StatusOK, "authorize success")
	} else {
		zap.S().Warnw("Authorization: passwords do not match, aborting", "username", username)
		c.String(http.StatusUnauthorized, "Invalid credentials")
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

	users := []User{}
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
