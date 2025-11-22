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
	"main/internal/middleware"
	"main/internal/pg"
	profile "main/internal/profile/posts"
)

var (
	sessionManager *scs.SessionManager
	redisPool      *redis.Pool
)

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
	sessionManager.Cookie.Secure = false // Set to true in production with HTTPS

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
	r.GET("/community/:id", pg.GetCommunityByID)
	r.GET("/communities", pg.GetAllCommunities)

	r.GET("/:userID/posts", profile.GetUserPosts)
	r.GET("/posts/:postID", profile.GetPost)
	r.GET("/graph-data", pg.GetGraphData)

	// Protected routes
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware(sessionManager))
	{
		api.GET("/users", users.GetAllUsers)
		api.POST("/posts", profile.CreatePost)
		api.PUT("/posts/:postID", profile.UpdatePost)
		api.DELETE("/posts/:postID", profile.DeletePost)
		api.POST("/posts/:postID/like", profile.LikePost)
		api.DELETE("/posts/:postID/like", profile.UnlikePost)
		api.POST("/logout", func(c *gin.Context) {
			users.LogoutUser(c, sessionManager)
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
