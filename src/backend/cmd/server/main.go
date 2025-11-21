package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not set")
	}

	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	log.Println("fucll")
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	log.Println("Database is ready to accept connections")
	// Здесь дальше код запуска сервера и обработчиков

	r := gin.Default()

	// Вешаем CORS (или ваше middleware)
	r.Use(func(c *gin.Context) {
		c.Writer.Header().
			Set("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Writer.Header().
			Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}
		c.Next()
	})

	// Вместо router.GET, router.POST и т.д.
	// подключите ваши обработчики через gin, например:
	// r.GET("/hello/:name", gin.HandlerFunc(hello))
	r.GET("/register", func(c *gin.Context) {
		c.File("static/form.html")
	})
	// r.POST("/register", gin.HandlerFunc(registerHandler))
	r.POST("/auth", gin.HandlerFunc(authorize))
	r.NoRoute(func(c *gin.Context) {
		c.String(404, "not found")
	})

	// Запуск вашего сервера
	r.Run(":8080")
}
