package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"main/internal/auth/users"
	"main/internal/pg"
)

func initLogger() {
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not set")
	}

	initLogger()
	zap.S().Debug("Application starting...")

	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		zap.S().Fatal("DATABASE_URL is not set")
	}

	var err error
	pg.DB, err = sql.Open("postgres", connStr)
	if err != nil {
		zap.S().Fatal(err)
	}
	defer pg.DB.Close()

	if err := pg.DB.Ping(); err != nil {
		zap.S().Fatal(err)
	}

	log.Println("Database is ready to accept connections")

	r := gin.Default()

	// CORS Middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().
			Set("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Writer.Header().
			Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().
			Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}
		c.Next()
	})

	r.GET("/register", func(c *gin.Context) {
		c.File("static/form.html")
	})

	r.POST("/auth", gin.HandlerFunc(users.AuthorizeUser))
	r.POST("/register", gin.HandlerFunc(users.RegisterHandler))
	profileRepo := profile.NewRepository(pg.DB)
	profileHandler := profile.NewHandler(profileRepo)
	profile.RegisterRoutes(
		r,
		profileHandler,
	) // ← ВСЕ маршруты регистрируются автоматически!

	r.POST("/profile/posts", handler.CreatePost)
	r.GET("/profile/:userID/posts", handler.GetUserPosts)
	r.GET("/profile/posts/:postID", handler.GetPost)
	r.PUT("/profile/posts/:postID", handler.UpdatePost)
	r.DELETE("/profile/posts/:postID", handler.DeletePost)
	r.POST("/profile/posts/:postID/like", handler.LikePost)
	r.DELETE("profile/posts/:postID/like", handler.UnlikePost)

	r.NoRoute(func(c *gin.Context) {
		c.String(404, "not found")
	})

	log.Println("Starting server on :8080")
	r.Run(":8080")
}
