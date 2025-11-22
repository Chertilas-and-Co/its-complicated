package communities

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"main/internal/models"
)

var DB *sql.DB

// CREATE TABLE communities (
//
//	id BIGSERIAL PRIMARY KEY,
//	name VARCHAR(255) NOT NULL,
//	description TEXT,
//	is_private BOOLEAN DEFAULT FALSE,
//	created_by BIGINT NOT NULL REFERENCES users(id),
//	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
//
// );

func InsertCommunityInDB(c *gin.Context) {
	var com models.Community

	if err := c.ShouldBindJSON(&com); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	var exists bool
	err := DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)", com.CreatedBy).
		Scan(&exists)
	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"message": "No user with this ID"})
		log.Println(err)
		return
	}
	result, err := DB.Exec(
		"INSERT INTO communities (name, description, is_private, created_by, created_at) VALUES ($1, $2, $3, $4, $5)",
		com.Name,
		com.Description,
		com.IsPrivate,
		com.CreatedBy,
		com.CreatedAt,
	)
	if err != nil {
		log.Println(err)
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		log.Println(err)
	}
	log.Println("Last inserted id:", lastInsertID)
}

type CommunityResponse struct {
	ID          int           `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	IsPrivate   bool          `json:"is_private"`
	Subscribers []models.User `json:"subscribers"`
	CreatedBy   int64         `json:"created_by"`
	CreatedAt   time.Time     `json:"created_at"` // или time.Time
}

func GetCommunityByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid id"})
		return
	}

	var community CommunityResponse
	query := `SELECT id, name, description, is_private, created_by, created_at FROM communities WHERE id = $1`

	err = DB.QueryRow(query, id).Scan(
		&community.ID,
		&community.Name,
		&community.Description,
		&community.IsPrivate,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(404, gin.H{"error": "community not found"})
			return
		}
		c.JSON(500, gin.H{"error": "unknown error"})
		return
	}

	// Получаем user_id подписчиков для этого сообщества
	subsQuery := `SELECT user_id FROM community_subscribers WHERE community_id = $1`
	rows, err := DB.Query(subsQuery, id)
	if err != nil {
		c.JSON(500, gin.H{"error": "cannot fetch subscribers"})
		return
	}
	defer rows.Close()

	var subscribers []models.User
	for rows.Next() {
		var userID int64
		if err := rows.Scan(&userID); err != nil {
			c.JSON(500, gin.H{"error": "error reading subscriber"})
			return
		}
		subscribers = append(subscribers, models.User{ID: userID})
	}
	if err := rows.Err(); err != nil {
		c.JSON(500, gin.H{"error": "error processing subscribers"})
		return
	}
	community.Subscribers = subscribers

	c.JSON(200, community)
}

type SubscribeRequest struct {
	UserId      int `json:"user_id"`
	CommunityID int `json:"community_id"`
}

func SubscribeToCommunity(c *gin.Context) {
	var subReq SubscribeRequest

	if err := c.ShouldBindJSON(&subReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	var exists bool
	err := DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", subReq.UserId).
		Scan(&exists)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no user with such id"})
		return
	}

	err = DB.QueryRow("SELECT EXISTS(SELECT 1 FROM commuinties WHERE id = $1)", subReq.CommunityID).
		Scan(&exists)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if !exists {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "no community with such id"},
		)
		return
	}

	result, err := DB.Exec(
		"INSERT INTO community_subscriptions (user_id, community_id) VALUES ($1, $2)",
		subReq.UserId,
		subReq.CommunityID,
	)
	if err != nil {
		log.Println(err)
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		log.Println(err)
	}
	log.Println("Last inserted id:", lastInsertID)
}

func GetAllCommunities(c *gin.Context) {
	query := `SELECT id, name, description, is_private, created_by, created_at FROM communities`

	rows, err := DB.Query(query)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "cannot fetch communities"},
		)
		return
	}
	defer rows.Close()

	var communities []CommunityResponse
	for rows.Next() {
		var community CommunityResponse
		err := rows.Scan(
			&community.ID,
			&community.Name,
			&community.Description,
			&community.IsPrivate,
			&community.CreatedBy,
			&community.CreatedAt,
		)
		if err != nil {
			c.JSON(
				http.StatusInternalServerError,
				gin.H{"error": "error reading communities"},
			)
			return
		}
		// По желанию можно загрузить подписчиков для каждого сообщества, но это будет сложнее
		communities = append(communities, community)
	}
	if err := rows.Err(); err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "error processing communities"},
		)
		return
	}

	c.JSON(http.StatusOK, communities)
}
