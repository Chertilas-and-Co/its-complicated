package pg

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"main/internal/models"
	"net/http"
	"strconv"
	"time"

	"log"
)

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
	result, err := DB.Exec("INSERT INTO communities (name, description, is_private, created_by, created_at) VALUES ($1, $2, $3, $4, $5)",
		com.Name, com.Description, com.IsPrivate, com.CreatedBy, com.CreatedAt)
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
	query := `SELECT id, name, description, is_private FROM communities WHERE id = $1`

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
	subsQuery := `SELECT users.id, users.username FROM community_subscriptions JOIN users ON users.id = community_subscriptions.user_id WHERE community_id = $1`
	println(id)
	rows, err := DB.Query(subsQuery, id)
	if err != nil {
		c.JSON(500, gin.H{"error": "cannot fetch subscribers"})
		return
	}
	defer rows.Close()

	var subscribers []models.User
	for rows.Next() {
		var userID int64
		var username string
		if err := rows.Scan(&userID, &username); err != nil {
			c.JSON(500, gin.H{"error": "error reading subscriber"})
			return
		}
		subscribers = append(subscribers, models.User{ID: userID, Username: username})
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
	err := DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", subReq.UserId).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no user with such id"})
		return
	}

	err = DB.QueryRow("SELECT EXISTS(SELECT 1 FROM commuinties WHERE id = $1)", subReq.CommunityID).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no community with such id"})
		return
	}

	result, err := DB.Exec("INSERT INTO community_subscriptions (user_id, community_id) VALUES ($1, $2)", subReq.UserId, subReq.CommunityID)
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot fetch communities"})
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error reading communities"})
			return
		}
		// По желанию можно загрузить подписчиков для каждого сообщества, но это будет сложнее
		communities = append(communities, community)
	}
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error processing communities"})
		return
	}

	c.JSON(http.StatusOK, communities)
}

// DenormalizedLink is the response structure the frontend expects.
type DenormalizedLink struct {
	ID1               int64  `json:"id_1"`
	ID2               int64  `json:"id_2"`
	Subscribers1      int    `json:"subscribers_1"`
	Subscribers2      int    `json:"subscribers_2"`
	CommonSubscribers int    `json:"common_subscribers"`
	Name1             string `json:"name_1"`
	Desc1             string `json:"desc_1"`
	Name2             string `json:"name_2"`
	Desc2             string `json:"desc_2"`
}

// Internal structs for fetching data from DB
type graphNode struct {
	ID   int64
	Name string
	Size int
}

type graphLink struct {
	Source int64
	Target int64
	Value  int
}

func GetGraphData(c *gin.Context) {
	// 1. Get all nodes (communities and their sizes)
	nodesQuery := `
		SELECT c.id, c.name, COUNT(s.user_id) as size
		FROM communities c
		LEFT JOIN community_subscriptions s ON c.id = s.community_id
		GROUP BY c.id, c.name;`

	rows, err := DB.Query(nodesQuery)
	if err != nil {
		log.Printf("Error querying graph nodes: %v", err)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to fetch community data"},
		)
		return
	}
	defer rows.Close()

	nodeMap := make(map[int64]graphNode)
	for rows.Next() {
		var node graphNode
		if err := rows.Scan(&node.ID, &node.Name, &node.Size); err != nil {
			log.Printf("Error scanning graph node: %v", err)
			c.JSON(
				http.StatusInternalServerError,
				gin.H{"error": "Failed to process community data"},
			)
			return
		}
		nodeMap[node.ID] = node
	}
	if err = rows.Err(); err != nil {
		log.Printf("Error iterating graph nodes: %v", err)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to process community data"},
		)
		return
	}

	// 2. Get all links (intersections)
	linksQuery := `
		SELECT s1.community_id as source, s2.community_id as target, COUNT(s1.user_id) AS value
		FROM community_subscriptions s1
		JOIN community_subscriptions s2 ON s1.user_id = s2.user_id AND s1.community_id < s2.community_id
		GROUP BY s1.community_id, s2.community_id
		HAVING COUNT(s1.user_id) > 0;`

	rows, err = DB.Query(linksQuery)
	if err != nil {
		log.Printf("Error querying graph links: %v", err)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to fetch community links"},
		)
		return
	}
	defer rows.Close()

	var links []graphLink
	for rows.Next() {
		var link graphLink
		if err := rows.Scan(&link.Source, &link.Target, &link.Value); err != nil {
			log.Printf("Error scanning graph link: %v", err)
			c.JSON(
				http.StatusInternalServerError,
				gin.H{"error": "Failed to process community links"},
			)
			return
		}
		links = append(links, link)
	}
	if err = rows.Err(); err != nil {
		log.Printf("Error iterating graph links: %v", err)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to process community links"},
		)
		return
	}

	// 3. Transform data into the flat format the frontend expects
	var response []DenormalizedLink
	for _, link := range links {
		sourceNode := nodeMap[link.Source]
		targetNode := nodeMap[link.Target]

		denormalized := DenormalizedLink{
			ID1:               sourceNode.ID,
			ID2:               targetNode.ID,
			Subscribers1:      sourceNode.Size,
			Subscribers2:      targetNode.Size,
			CommonSubscribers: link.Value,
			Name1:             sourceNode.Name,
			Desc1:             "Участников: " + strconv.Itoa(sourceNode.Size),
			Name2:             targetNode.Name,
			Desc2:             "Участников: " + strconv.Itoa(targetNode.Size),
		}
		response = append(response, denormalized)
	}

	if response == nil {
		response = make([]DenormalizedLink, 0)
	}

	c.JSON(http.StatusOK, response)
}

func GetCommunitySubscribers(c *gin.Context) {
	idStr := c.Param("id")
	communityID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid community id"})
		return
	}

	var exists bool
	err = DB.QueryRow("SELECT EXISTS(SELECT 1 FROM communities WHERE id = $1)", communityID).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error checking community"})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "community not found"})
		return
	}

	query := `
		SELECT u.id, u.username, u.email, u.bio, u.avatar_url
		FROM community_subscriptions cs
		JOIN users u ON cs.user_id = u.id
		WHERE cs.community_id = $1
		ORDER BY u.id;
	`

	rows, err := DB.Query(query, communityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot fetch subscribers"})
		return
	}
	defer rows.Close()

	var subscribers []models.User
	for rows.Next() {
		var user models.User
		var bio sql.NullString
		var avatarURL sql.NullString

		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&bio,
			&avatarURL,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error reading subscriber"})
			return
		}

		// Преобразуем sql.NullString в string
		if bio.Valid {
			user.Bio = bio.String
		} else {
			user.Bio = ""
		}
		if avatarURL.Valid {
			user.AvatarURL = avatarURL.String
		} else {
			user.AvatarURL = ""
		}

		subscribers = append(subscribers, user)
	}

	c.JSON(http.StatusOK, subscribers)
}
