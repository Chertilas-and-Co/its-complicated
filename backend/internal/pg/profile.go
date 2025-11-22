package pg

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Errors for user profile operations
var (
	ErrUserNotFound  = errors.New("user not found")
	ErrInvalidUserID = errors.New("invalid user id")
	ErrUnauthorized  = errors.New("unauthorized - can only edit your own profile")
)

// UpdateUserProfileRequest - JSON структура для обновления профиля
type UpdateUserProfileRequest struct {
	Username  string `json:"username" binding:"max=255"`
	Email     string `json:"email" binding:"max=255"`
	Bio       string `json:"bio" binding:"max=500"`
	AvatarURL string `json:"avatar_url" binding:"max=500"`
}

// UserProfileResponse - расширенный ответ с информацией о профиле
type UserProfileResponse struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Bio       string    `json:"bio"`
	AvatarURL string    `json:"avatar_url"`
	CreatedAt time.Time `json:"created_at"`
}

// UserStatsResponse - статистика пользователя
type UserStatsResponse struct {
	UserID         int64 `json:"user_id"`
	PostCount      int   `json:"post_count"`
	FriendsCount   int   `json:"friends_count"`
	FollowersCount int   `json:"followers_count"`
}

// UserPublicResponse - публичный профиль (без email)
type UserPublicResponse struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Bio       string `json:"bio"`
	AvatarURL string `json:"avatar_url"`
}

// GetUserProfile - получить профиль пользователя (публичная информация)
// GET /api/users/:userID
// Не требует авторизацию
func GetUserProfile(c *gin.Context) {
	userIDParam := c.Param("userID")
	userID, err := strconv.ParseInt(userIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	var user UserPublicResponse
	err = DB.QueryRow(
		`SELECT id, username, bio, avatar_url FROM users WHERE id = $1`,
		userID,
	).Scan(&user.ID, &user.Username, &user.Bio, &user.AvatarURL)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		zap.S().Errorf("Failed to fetch user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateProfile - обновить профиль пользователя
// PUT /api/me/profile
// Требует авторизацию
func UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req UpdateUserProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Получаем текущие данные
	var user UserProfileResponse
	var email sql.NullString

	err := DB.QueryRow(
		`SELECT id, username, email, bio, avatar_url, created_at FROM users WHERE id = $1`,
		userID.(int64),
	).Scan(&user.ID, &user.Username, &email, &user.Bio, &user.AvatarURL, &user.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		zap.S().Errorf("Failed to fetch user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user"})
		return
	}

	if email.Valid {
		user.Email = email.String
	}

	// Обновляем только переданные поля
	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Bio != "" {
		user.Bio = req.Bio
	}
	if req.AvatarURL != "" {
		user.AvatarURL = req.AvatarURL
	}

	// Обновляем в БД
	_, err = DB.Exec(
		`UPDATE users SET username = $1, email = $2, bio = $3, avatar_url = $4 WHERE id = $5`,
		user.Username,
		user.Email,
		user.Bio,
		user.AvatarURL,
		userID.(int64),
	)

	if err != nil {
		zap.S().Errorf("Failed to update user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "profile updated successfully",
		"user":    user,
	})
}

// SearchUsers - поиск пользователей по username
// GET /api/users/search?query=john
// Не требует авторизацию
func SearchUsers(c *gin.Context) {
	query := c.Query("query")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "search query is required"})
		return
	}

	// Защита от слишком коротких поисков
	if len(query) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "search query must be at least 2 characters"})
		return
	}

	searchQuery := `
	SELECT id, username, bio, avatar_url
	FROM users
	WHERE username ILIKE $1
	ORDER BY username
	LIMIT 20
	`

	rows, err := DB.Query(searchQuery, "%"+query+"%")
	if err != nil {
		zap.S().Errorf("Failed to search users: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to search users"})
		return
	}
	defer rows.Close()

	var users []UserPublicResponse
	for rows.Next() {
		var user UserPublicResponse
		if err := rows.Scan(&user.ID, &user.Username, &user.Bio, &user.AvatarURL); err != nil {
			zap.S().Errorf("Error scanning user: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error reading users"})
			return
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		zap.S().Errorf("Error iterating users: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error processing users"})
		return
	}

	if users == nil {
		users = make([]UserPublicResponse, 0)
	}

	c.JSON(http.StatusOK, users)
}
