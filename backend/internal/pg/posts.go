package pg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"

	"main/internal/models"
)

var (
	ErrPostNotFound = errors.New("post not found")
	ErrAlreadyLiked = errors.New("already liked")
	ErrNotLiked     = errors.New("not liked")
)

// CreatePost creates a new post in the database
// Возвращает созданный пост с заполненным ID и временем создания
func CreatePost(ctx context.Context, post *models.Post) error {
	const query = `
		INSERT INTO posts (title, text, pic_url, community_id, author_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	err := DB.QueryRowContext(ctx, query,
		post.Title,       // $1 - название поста
		post.Text,        // $2 - содержание поста
		post.PicURL,      // $3 - ссылка на картинку
		post.CommunityID, // $4 - ID сообщества (= ID пользователя для профиля)
		post.AuthorID,    // $5 - ID автора
	).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create post: %w", err)
	}

	return nil
}

// GetPostByID retrieves a single post by its ID
// Возвращает ошибку ErrPostNotFound если пост не найден
func GetPostByID(
	ctx context.Context,
	postID int64,
) (*models.Post, error) {
	const query = `
		SELECT id, title, text, pic_url, community_id, author_id, created_at, updated_at
		FROM posts
		WHERE id = $1
	`

	post := &models.Post{}

	err := DB.QueryRowContext(ctx, query, postID).Scan(
		&post.ID,
		&post.Title,
		&post.Text,
		&post.PicURL,
		&post.CommunityID,
		&post.AuthorID,
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrPostNotFound
		}
		return nil, fmt.Errorf("failed to fetch post: %w", err)
	}

	return post, nil
}

// GetUserPosts retrieves all posts for a specific user with pagination
// CommunityID = UserID для постов в профиле
// Возвращает слайс постов, общее количество постов и ошибку
func GetUserPosts(
	ctx context.Context,
	userID int64,
	limit, offset int,
) ([]*models.Post, int64, error) {
	// Получаем общее количество постов пользователя
	const countQuery = `
		SELECT COUNT(*)
		FROM posts
		WHERE community_id = $1
	`

	var total int64
	err := DB.QueryRowContext(ctx, countQuery, userID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count posts: %w", err)
	}

	// Получаем посты с учетом пагинации
	const postsQuery = `
		SELECT id, title, text, pic_url, community_id, author_id, created_at, updated_at
		FROM posts
		WHERE community_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := DB.QueryContext(ctx, postsQuery, userID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch posts: %w", err)
	}
	defer rows.Close()

	posts := make([]*models.Post, 0, limit)

	for rows.Next() {
		post := &models.Post{}

		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Text,
			&post.PicURL,
			&post.CommunityID,
			&post.AuthorID,
			&post.CreatedAt,
			&post.UpdatedAt,
		)

		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan post: %w", err)
		}

		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("row iteration error: %w", err)
	}

	return posts, total, nil
}

// UpdatePost updates an existing post
// Обновляет title, text, pic_url и updated_at
func UpdatePost(ctx context.Context, post *models.Post) error {
	const query = `
		UPDATE posts
		SET title = $1, text = $2, pic_url = $3, updated_at = NOW()
		WHERE id = $4 AND author_id = $5
		RETURNING updated_at
	`

	err := DB.QueryRowContext(ctx, query,
		post.Title,    // $1 - новое название
		post.Text,     // $2 - новое содержание
		post.PicURL,   // $3 - новая ссылка на картинку
		post.ID,       // $4 - ID поста
		post.AuthorID, // $5 - ID автора (проверка прав)
	).Scan(&post.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return ErrPostNotFound
		}
		return fmt.Errorf("failed to update post: %w", err)
	}

	return nil
}

// DeletePost deletes a post by its ID
func DeletePost(ctx context.Context, postID int64) error {
	const query = `
		DELETE FROM posts
		WHERE id = $1
	`

	result, err := DB.ExecContext(ctx, query, postID)
	if err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrPostNotFound
	}

	return nil
}

// LikePost adds a like to a post
// Возвращает ошибку ErrAlreadyLiked если пользователь уже лайкнул этот пост
func LikePost(ctx context.Context, postID, userID int64) error {
	const query = `
		INSERT INTO post_likes (post_id, user_id, created_at)
		VALUES ($1, $2, NOW())
	`

	_, err := DB.ExecContext(ctx, query, postID, userID)
	if err != nil {
		// Проверяем на нарушение уникальности (уже существует лайк)
		if err, ok := err.(*pq.Error); ok && err.Code == "23505" {
			return ErrAlreadyLiked
		}
		return fmt.Errorf("failed to like post: %w", err)
	}

	return nil
}

// UnlikePost removes a like from a post
// Возвращает ошибку ErrNotLiked если пользователь не лайкал этот пост
func UnlikePost(
	ctx context.Context,
	postID, userID int64,
) error {
	const query = `
		DELETE FROM post_likes
		WHERE post_id = $1 AND user_id = $2
	`

	result, err := DB.ExecContext(ctx, query, postID, userID)
	if err != nil {
		return fmt.Errorf("failed to unlike post: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrNotLiked
	}

	return nil
}

// GetCommunityPosts retrieves all posts for a specific community with pagination
// Возвращает слайс постов, общее количество постов и ошибку
func GetCommunityPosts(
	ctx context.Context,
	communityID int64,
	limit, offset int,
) ([]*models.Post, int64, error) {
	// Получаем общее количество постов сообщества
	const countQuery = `
		SELECT COUNT(*)
		FROM posts
		WHERE community_id = $1 AND author_id IS NOT NULL
	`

	var total int64
	err := DB.QueryRowContext(ctx, countQuery, communityID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count community posts: %w", err)
	}

	// Получаем посты с учетом пагинации
	const postsQuery = `
		SELECT id, title, text, pic_url, community_id, author_id, created_at, updated_at
		FROM posts
		WHERE community_id = $1 AND author_id IS NOT NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := DB.QueryContext(ctx, postsQuery, communityID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch community posts: %w", err)
	}
	defer rows.Close()

	posts := make([]*models.Post, 0, limit)

	for rows.Next() {
		post := &models.Post{}

		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Text,
			&post.PicURL,
			&post.CommunityID,
			&post.AuthorID,
			&post.CreatedAt,
			&post.UpdatedAt,
		)

		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan community post: %w", err)
		}

		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("row iteration error for community posts: %w", err)
	}

	return posts, total, nil
}
