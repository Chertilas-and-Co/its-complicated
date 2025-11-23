package models

import "time"

// User - пользователь
type User struct {
	ID           int64  `json:"id"         db:"id"`
	Username     string `json:"username"   db:"username"`
	Email        string `json:"email"      db:"email"`
	PasswordHash string `json:"-"          db:"password_hash"`
	Salt         []byte `json:"-"          db:"salt"`
	Bio          string `json:"bio"        db:"bio"`
	Avatar       []byte `json:"-"          db:"avatar"`
	AvatarURL    string `json:"avatar_url" db:"avatar_url"`
}

// Friendship - отношение дружбы
type Friendship struct {
	ID       int64  `json:"id"        db:"id"`
	UserID   int64  `json:"user_id"   db:"user_id"`
	FriendID int64  `json:"friend_id" db:"friend_id"`
	Status   string `json:"status"    db:"status"` // pending, accepted, blocked
}

// Community - сообщество
type Community struct {
	ID          int64     `json:"id"          db:"id"`
	Name        string    `json:"name"        db:"name"`
	Description string    `json:"description" db:"description"`
	IsPrivate   bool      `json:"is_private"  db:"is_private"`
	CreatedBy   int64     `json:"created_by"  db:"created_by"`
	CreatedAt   time.Time `json:"created_at"  db:"created_at"`
	Admins      []int64   `json:"admins"`
	Writers     []int64   `json:"writers"`
}

// CommunitySubscription - подписка на сообщество
type CommunitySubscription struct {
	ID          int64 `json:"id"           db:"id"`
	UserID      int64 `json:"user_id"      db:"user_id"`
	CommunityID int64 `json:"community_id" db:"community_id"`
}

// Post - пост в сообществе
type Post struct {
	ID             int64     `json:"id"`
	Title          string    `json:"title"`
	Text           string    `json:"text"`
	PicURL         string    `json:"pic_url"`
	CommunityID    int64     `json:"community_id"`
	AuthorID       int64     `json:"author_id"`
	AuthorUsername string    `json:"author_username"` // Добавляем имя пользователя автора
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// Roles
const (
	RoleAdmin      = "admin"
	RoleWriter     = "writer"
	RoleSubscriber = "subscriber"
)

// Friendship statuses
const (
	FriendshipPending  = "pending"
	FriendshipAccepted = "accepted"
	FriendshipBlocked  = "blocked"
)

// Comment represents a comment on a post
type Comment struct {
	ID        int64     `json:"id"`
	PostID    int64     `json:"post_id"`
	UserID    int64     `json:"user_id"`
	Username  string    `json:"username"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}
