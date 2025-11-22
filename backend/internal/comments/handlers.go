package comments

import (
	"errors"
	"net/http"
	"strconv"

	"main/internal/models"
	"main/internal/pg"

	"github.com/gin-gonic/gin"
)

var (
	ErrCommentNotFound = errors.New("comment not found")
)

// CreateComment creates a new comment on a post
// POST /api/posts/:postID/comments
// Requires: UserID in context
func CreateComment(c *gin.Context) {
	userID, exists := c.Get("UserID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	username, exists := c.Get("Username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "username not found"})
		return
	}

	postIDParam := c.Param("postID")
	postID, err := strconv.ParseInt(postIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post ID"})
		return
	}

	var req struct {
		Content string `json:"content" binding:"required,min=1,max=5000"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	comment := &models.Comment{
		PostID:   postID,
		UserID:   userID.(int64),
		Username: username.(string),
		Content:  req.Content,
	}

	if err := pg.CreateComment(c.Request.Context(), comment); err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "failed to create comment"},
		)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "comment created successfully",
		"comment": comment,
	})
}

// GetCommentsByPostID retrieves all comments for a specific post
// GET /api/posts/:postID/comments?limit=20&offset=0
func GetCommentsByPostID(c *gin.Context) {
	postIDParam := c.Param("postID")
	postID, err := strconv.ParseInt(postIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post ID"})
		return
	}

	limit := 20
	offset := 0

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	if o := c.Query("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	comments, total, err := pg.GetCommentsByPostID(
		c.Request.Context(),
		postID,
		limit,
		offset,
	)

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "failed to fetch comments"},
		)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"comments": comments,
		"total":    total,
		"limit":    limit,
		"offset":   offset,
	})
}

// GetComment retrieves a single comment by ID
// GET /api/posts/:postID/comments/:commentID
func GetComment(c *gin.Context) {
	commentIDParam := c.Param("commentID")
	commentID, err := strconv.ParseInt(commentIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid comment ID"})
		return
	}

	comment, err := pg.GetCommentByID(c.Request.Context(), commentID)
	if err != nil {
		if errors.Is(err, ErrCommentNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "comment not found"})
			return
		}
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "failed to fetch comment"},
		)
		return
	}

	c.JSON(http.StatusOK, comment)
}

// UpdateComment updates a comment
// PUT /api/posts/:postID/comments/:commentID
// Requires: UserID in context (must be comment author)
func UpdateComment(c *gin.Context) {
	userID, exists := c.Get("UserID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	commentIDParam := c.Param("commentID")
	commentID, err := strconv.ParseInt(commentIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid comment ID"})
		return
	}

	comment, err := pg.GetCommentByID(c.Request.Context(), commentID)
	if err != nil {
		if errors.Is(err, ErrCommentNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "comment not found"})
			return
		}
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "failed to fetch comment"},
		)
		return
	}

	// Check ownership
	if comment.UserID != userID.(int64) {
		c.JSON(
			http.StatusForbidden,
			gin.H{"error": "you can only edit your own comments"},
		)
		return
	}

	var req struct {
		Content string `json:"content" binding:"required,min=1,max=5000"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	comment.Content = req.Content

	if err := pg.UpdateComment(c.Request.Context(), comment); err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "failed to update comment"},
		)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "comment updated successfully",
		"comment": comment,
	})
}

// DeleteComment deletes a comment
// DELETE /api/posts/:postID/comments/:commentID
// Requires: UserID in context (must be comment author)
func DeleteComment(c *gin.Context) {
	userID, exists := c.Get("UserID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	commentIDParam := c.Param("commentID")
	commentID, err := strconv.ParseInt(commentIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid comment ID"})
		return
	}

	comment, err := pg.GetCommentByID(c.Request.Context(), commentID)
	if err != nil {
		if errors.Is(err, ErrCommentNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "comment not found"})
			return
		}
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "failed to fetch comment"},
		)
		return
	}

	// Check ownership
	if comment.UserID != userID.(int64) {
		c.JSON(
			http.StatusForbidden,
			gin.H{"error": "you can only delete your own comments"},
		)
		return
	}

	if err := pg.DeleteComment(c.Request.Context(), commentID); err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "failed to delete comment"},
		)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "comment deleted successfully"})
}
