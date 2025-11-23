package main

import (
	"database/sql"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"main/internal/auth/users"
	"main/internal/comments"
	"main/internal/communities"
	community "main/internal/community/posts"
	"main/internal/middleware"
	"main/internal/pg"
	"main/internal/profile/friends"
	profile "main/internal/profile/posts"
)

var (
	sessionManager *scs.SessionManager
	redisPool      *redis.Pool
)

func CreateCommunityHandlerWrapper(
	sm *scs.SessionManager,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		communities.CreateCommunityHandler(c, sm)
	}
}

func SubscribeToCommunityHandlerWrapper(sm *scs.SessionManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		pg.SubscribeToCommunity(c, sm)
	}
}

func UnsubscribeFromCommunityHandlerWrapper(sm *scs.SessionManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		pg.UnsubscribeFromCommunity(c, sm)
	}
}

func CheckSubscriptionStatusHandlerWrapper(sm *scs.SessionManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		pg.CheckSubscriptionStatus(c, sm)
	}
}

func initLogger() {
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)
}

func main() {
	if err := godotenv.Load(); err != nil {
		zap.S().Warn(".env file not found")
	}

	initLogger()

	sessionManager = scs.New()
	redisAddr := os.Getenv("DRAGONFLY_URL")
	if redisAddr == "" {
		redisAddr = "redis://localhost:6379"
	}

	redisPool = &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.DialURL(redisAddr)
		},
	}
	zap.S().Info("Successfully configured Redis connection pool for Dragonfly")

	sessionManager.Store = redisstore.New(redisPool)
	sessionManager.Lifetime = 24 * time.Hour
	sessionManager.Cookie.Name = "session_id"
	sessionManager.Cookie.HttpOnly = true
	sessionManager.Cookie.Persist = true
	sessionManager.Cookie.SameSite = http.SameSiteLaxMode
	sessionManager.Cookie.Secure = false
	// Removed explicit Domain and Path settings

	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		zap.S().Fatal("DATABASE_URL is not set")
	}

	var err error
	pg.DB, err = sql.Open("postgres", connStr)
	if err != nil {
		zap.S().Fatalf("Failed to open database connection: %v", err)
	}
	defer pg.DB.Close()

	if err := pg.DB.Ping(); err != nil {
		zap.S().Fatalf("Failed to ping database: %v", err)
	}
	zap.S().Info("Database is ready to accept connections")

	r := gin.Default()

	// CORS Middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().
			Set("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Writer.Header().
			Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().
			Set("Access-Control-Allow-Headers", "Content-Type, X-Requested-With")
		c.Writer.Header().
			Set("Access-Control-Allow-Credentials", "true")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}
		c.Next()
	})

	// Public routes
	r.POST("/register", users.RegisterUser)
	r.POST("/auth", func(c *gin.Context) {
		users.AuthorizeUser(c, sessionManager)
	})

	r.GET("/community/:id/subscribers", pg.GetCommunitySubscribers)
	r.GET("/community/:id", pg.GetCommunityByID)
	r.GET("/communities", pg.GetAllCommunities) // Added route to get all communities

	r.GET("/user/:userID/posts", profile.GetUserPosts)
	r.GET("/user/posts/:postID", profile.GetPost)

	r.GET("/community/:id/posts", community.GetCommunityPosts)
	r.GET("/community/:id/posts/:postID", community.GetPost)

	r.GET("/graph-data", pg.GetGraphData)

	r.GET("/community/:id/posts/:postID/comments?limit=20&offset=0", comments.GetCommentsByPostID)
	r.GET("/user/posts/:postID/comments?limit=20&offset=0", comments.GetCommentsByPostID)
	r.GET("/community/:id/posts/:postID/comments/:commentID", comments.GetComment)
	r.GET("user/posts/:postID/comments/:commentID", comments.GetComment)

	// Protected routes
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware(sessionManager))
	{
		api.GET("/user", pg.GetUserProfile)
		api.PUT("/user", pg.UpdateProfile)
		api.GET("/user/search?query=john", pg.SearchUsers)

		// New route for creating communities
		api.POST("/communities", CreateCommunityHandlerWrapper(sessionManager))

		// Community subscription routes
		api.POST("/community/:id/subscribe", SubscribeToCommunityHandlerWrapper(sessionManager))
		api.DELETE("/community/:id/subscribe", UnsubscribeFromCommunityHandlerWrapper(sessionManager))
		api.GET("/community/:id/is_subscribed", CheckSubscriptionStatusHandlerWrapper(sessionManager))

		api.POST("/community/:id/posts/:postID/comments", comments.CreateComment)
		api.POST("/user/posts/:postID/comments", comments.CreateComment)
		api.PUT("/community/:id/posts/:postID/comments/:commentID", comments.UpdateComment)
		api.PUT("/user/posts/:postID/comments/:commentID", comments.UpdateComment)
		api.DELETE("/community/:id/posts/:postID/comments/:commentID", comments.DeleteComment)
		api.DELETE("/user/posts/:postID/comments/:commentID", comments.DeleteComment)

		api.GET("/users", users.GetAllUsers)

		api.POST("/user/posts", profile.CreatePost)
		api.PUT("/user/posts/:postID", profile.UpdatePost)
		api.DELETE("/user/posts/:postID", profile.DeletePost)
		api.POST("/user/posts/:postID/like", profile.LikePost)
		api.DELETE("/user/posts/:postID/like", profile.UnlikePost)

		api.POST("/community/:id/posts", community.CreatePost)
		api.PUT("/community/:id/posts/:postID", community.UpdatePost)
		api.DELETE("/community/:id/posts/:postID", community.DeletePost)
		api.POST("/community/:id/posts/:postID/like", community.LikePost)
		api.DELETE("/community/:id/posts/:postID/like", community.UnlikePost)

		api.POST("/logout", func(c *gin.Context) {
			users.LogoutUser(c, sessionManager)
		})

		// Friend routes
		api.POST("/friends/requests", func(c *gin.Context) {
			friends.SendFriendRequestHandler(c, sessionManager)
		})
		api.GET("/friends", func(c *gin.Context) {
			friends.GetFriendsHandler(c, sessionManager)
		})
		api.DELETE("/friends/:friend_id", func(c *gin.Context) {
			friends.DeleteFriendHandler(c, sessionManager)
		})
		api.GET("/friends/requests/incoming", func(c *gin.Context) {
			friends.GetIncomingRequestsHandler(c, sessionManager)
		})
		api.PUT("/friends/requests/:request_id", func(c *gin.Context) {
			friends.UpdateFriendRequestHandler(c, sessionManager)
		})
	}

	r.NoRoute(func(c *gin.Context) {
		c.String(404, "not found")
	})

	zap.S().Info("Starting server on :8080")
	// Use http.ListenAndServe with the scs middleware wrapping the gin router
	if err := http.ListenAndServe(":8080", sessionManager.LoadAndSave(r)); err != nil {
		zap.S().Fatalf("Failed to start server: %v", err)
	}
}

// ABOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOBA
