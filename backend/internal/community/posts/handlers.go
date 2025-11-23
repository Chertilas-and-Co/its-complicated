package Community

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"main/internal/models"
	"main/internal/pg"
)

var (
	ErrPostNotFound = errors.New("post not found")
	ErrAlreadyLiked = errors.New("already liked")
	ErrNotLiked     = errors.New("not liked")
)

// Handler handles HTTP requests for profile posts
// CreatePost creates a new post
// POST /api/profile/posts
// Требует авторизацию (userID в контексте)
func CreatePost(c *gin.Context) {
	communityIDStr := c.Param("id") // Извлекаем communityID из URL
	communityID, err := strconv.ParseInt(communityIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid community ID"})
		return
	}

	userID, exists := c.Get("userID") // Извлекаем userID из контекста, установленного middleware
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req struct {
		Title  string `json:"title" binding:"required,max=255"`
		Text   string `json:"text" binding:"required"`
		PicURL string `json:"pic_url"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post := &models.Post{
		Title:       req.Title,
		Text:        req.Text,
		PicURL:      req.PicURL,
		CommunityID: communityID, // Для профиля - CommunityID = UserID
		AuthorID:    userID.(int64),
	}

	if err := pg.CreatePost(c.Request.Context(), post); err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "failed to create post"},
		)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "post created successfully",
		"post":    post,
	})
}

// GetCommunityPosts retrieves all posts for a community
// GET /community/:id/posts?limit=20&offset=40
func GetCommunityPosts(c *gin.Context) {
	communityIDParam := c.Param("id") // Извлекаем communityID из URL
	communityID, err := strconv.ParseInt(communityIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid community ID"})
		return
	}

	limit := 10
	offset := 0

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 &&
			parsed <= 100 {
			limit = parsed
		}
	}

	if o := c.Query("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	posts, total, err := pg.GetCommunityPosts( // Вызываем новую функцию pg.GetCommunityPosts
		c.Request.Context(),
		communityID,
		limit,
		offset,
	)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "failed to fetch posts"},
		)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"posts":  posts,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

// GetPost retrieves a single post
// GET /api/profile/posts/:postID
func GetPost(c *gin.Context) {
	postIDParam := c.Param("postID")
	postID, err := strconv.ParseInt(postIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post ID"})
		return
	}

	post, err := pg.GetPostByID(c.Request.Context(), postID)
	if err != nil {
		if errors.Is(err, ErrPostNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
			return
		}
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "failed to fetch post"},
		)
		return
	}

	c.JSON(http.StatusOK, post)
}

// UpdatePost updates a post
// PUT /api/profile/posts/:postID
// Требует авторизацию + проверку владельца
func UpdatePost(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	postIDParam := c.Param("postID")
	postID, err := strconv.ParseInt(postIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post ID"})
		return
	}

	post, err := pg.GetPostByID(c.Request.Context(), postID)
	if err != nil {
		if errors.Is(err, ErrPostNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
			return
		}
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "failed to fetch post"},
		)
		return
	}

	// Проверяем владельца поста
	if post.AuthorID != userID.(int64) {
		c.JSON(
			http.StatusForbidden,
			gin.H{"error": "you can only edit your own posts"},
		)
		return
	}

	var req struct {
		Title  string `json:"title" binding:"max=255"`
		Text   string `json:"text"`
		PicURL string `json:"pic_url"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Partial update - обновляем только переданные поля
	if req.Title != "" {
		post.Title = req.Title
	}
	if req.Text != "" {
		post.Text = req.Text
	}
	if req.PicURL != "" {
		post.PicURL = req.PicURL
	}

	if err := pg.UpdatePost(c.Request.Context(), post); err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "failed to update post"},
		)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "post updated successfully",
		"post":    post,
	})
}

// DeletePost deletes a post
// DELETE /api/profile/posts/:postID
// Требует авторизацию + проверку владельца
func DeletePost(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	postIDParam := c.Param("postID")
	postID, err := strconv.ParseInt(postIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post ID"})
		return
	}

	post, err := pg.GetPostByID(c.Request.Context(), postID)
	if err != nil {
		if errors.Is(err, ErrPostNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
			return
		}
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "failed to fetch post"},
		)
		return
	}

	// Проверяем владельца
	if post.AuthorID != userID.(int64) {
		c.JSON(
			http.StatusForbidden,
			gin.H{"error": "you can only delete your own posts"},
		)
		return
	}

	if err := pg.DeletePost(c.Request.Context(), postID); err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "failed to delete post"},
		)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "post deleted successfully"})
}

// LikePost adds a like to a post
// POST /api/profile/posts/:postID/like
// Требует авторизацию
func LikePost(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	postIDParam := c.Param("postID")
	postID, err := strconv.ParseInt(postIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post ID"})
		return
	}

	// Проверяем существует ли пост
	_, err = pg.GetPostByID(c.Request.Context(), postID)
	if err != nil {
		if errors.Is(err, ErrPostNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
			return
		}
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "failed to fetch post"},
		)
		return
	}

	if err := pg.LikePost(c.Request.Context(), postID, userID.(int64)); err != nil {
		if errors.Is(err, ErrAlreadyLiked) {
			c.JSON(
				http.StatusConflict,
				gin.H{"error": "you already liked this post"},
			)
			return
		}
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "failed to like post"},
		)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "post liked successfully"})
}

// UnlikePost removes a like from a post
// DELETE /api/profile/posts/:postID/like
// Требует авторизацию
func UnlikePost(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	postIDParam := c.Param("postID")
	postID, err := strconv.ParseInt(postIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post ID"})
		return
	}

	if err := pg.UnlikePost(c.Request.Context(), postID, userID.(int64)); err != nil {
		if errors.Is(err, ErrNotLiked) {
			c.JSON(
				http.StatusNotFound,
				gin.H{"error": "you haven't liked this post"},
			)
			return
		}
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "failed to unlike post"},
		)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "post unliked successfully"})
}
