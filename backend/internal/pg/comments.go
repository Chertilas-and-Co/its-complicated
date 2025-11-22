// internal/pg/comments.go
package pg

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"main/internal/models"
	"time"
)

// CreateComment inserts a new comment into the database
func CreateComment(ctx context.Context, comment *models.Comment) error {
	query := `
		INSERT INTO comments (post_id, user_id, username, content, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`

	err := DB.QueryRowContext(ctx, query,
		comment.PostID,
		comment.UserID,
		comment.Username,
		comment.Content,
		time.Now().UTC(),
	).Scan(&comment.ID, &comment.CreatedAt)

	if err != nil {
		log.Printf("Error creating comment: %v", err)
		return err
	}

	return nil
}

// GetCommentByID retrieves a single comment by ID
func GetCommentByID(ctx context.Context, commentID int64) (*models.Comment, error) {
	query := `
		SELECT id, post_id, user_id, username, content, created_at
		FROM comments
		WHERE id = $1
	`

	var comment models.Comment

	err := DB.QueryRowContext(ctx, query, commentID).Scan(
		&comment.ID,
		&comment.PostID,
		&comment.UserID,
		&comment.Username,
		&comment.Content,
		&comment.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("comment not found")
		}
		log.Printf("Error getting comment: %v", err)
		return nil, err
	}

	return &comment, nil
}

// GetCommentsByPostID retrieves all comments for a specific post with pagination
func GetCommentsByPostID(ctx context.Context, postID int64, limit, offset int) ([]models.Comment, int, error) {
	// Get total count
	countQuery := `SELECT COUNT(*) FROM comments WHERE post_id = $1`
	var total int
	err := DB.QueryRowContext(ctx, countQuery, postID).Scan(&total)
	if err != nil {
		log.Printf("Error counting comments: %v", err)
		return nil, 0, err
	}

	// Get paginated results
	query := `
		SELECT id, post_id, user_id, username, content, created_at
		FROM comments
		WHERE post_id = $1
		ORDER BY created_at ASC
		LIMIT $2 OFFSET $3
	`

	rows, err := DB.QueryContext(ctx, query, postID, limit, offset)
	if err != nil {
		log.Printf("Error querying comments: %v", err)
		return nil, 0, err
	}
	defer rows.Close()

	var comments []models.Comment

	for rows.Next() {
		var comment models.Comment

		err := rows.Scan(
			&comment.ID,
			&comment.PostID,
			&comment.UserID,
			&comment.Username,
			&comment.Content,
			&comment.CreatedAt,
		)

		if err != nil {
			log.Printf("Error scanning comment: %v", err)
			return nil, 0, err
		}

		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating comments: %v", err)
		return nil, 0, err
	}

	if comments == nil {
		comments = []models.Comment{}
	}

	return comments, total, nil
}

// UpdateComment updates an existing comment
func UpdateComment(ctx context.Context, comment *models.Comment) error {
	query := `
		UPDATE comments
		SET content = $1
		WHERE id = $2
	`

	result, err := DB.ExecContext(ctx, query, comment.Content, comment.ID)
	if err != nil {
		log.Printf("Error updating comment: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected: %v", err)
		return err
	}

	if rowsAffected == 0 {
		return errors.New("comment not found")
	}

	return nil
}

// DeleteComment deletes a comment by ID
func DeleteComment(ctx context.Context, commentID int64) error {
	query := `DELETE FROM comments WHERE id = $1`

	result, err := DB.ExecContext(ctx, query, commentID)
	if err != nil {
		log.Printf("Error deleting comment: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected: %v", err)
		return err
	}

	if rowsAffected == 0 {
		return errors.New("comment not found")
	}

	return nil
}

// GetCommentCount returns the number of comments for a post
func GetCommentCount(ctx context.Context, postID int64) (int, error) {
	query := `SELECT COUNT(*) FROM comments WHERE post_id = $1`

	var count int
	err := DB.QueryRowContext(ctx, query, postID).Scan(&count)
	if err != nil {
		log.Printf("Error counting comments: %v", err)
		return 0, err
	}

	return count, nil
}
