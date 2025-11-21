package main

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"time"
)

// структура-обёртка для ResponseWriter — перехватывает WriteHeader
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

// middleware для логирования
func handleLogger(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		start := time.Now()
		sr := &statusRecorder{ResponseWriter: w, status: 200}
		h(sr, r, ps)
		duration := time.Since(start)
		log.Printf("%s %s %d %s", r.Method, r.URL.Path, sr.status, duration)
	}
}

func handlerFuncLogger(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sr := &statusRecorder{ResponseWriter: w, status: 200}
		h(sr, r)
		duration := time.Since(start)
		log.Printf("%s %s %d %s", r.Method, r.URL.Path, sr.status, duration)
	}
}

func hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.WriteHeader(200)
	fmt.Fprintf(w, "Hello, %s!", ps.ByName("name"))
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

func getRegPage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	http.ServeFile(w, r, "static/form.html")
}

func getAuthPage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	http.ServeFile(w, r, "static/auth.html")
}

func registerHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Form parse error", http.StatusBadRequest)
		return
	}

	username := r.FormValue("login")

	salt, _ := generateSalt(32)
	password1 := r.FormValue("password")
	password2 := r.FormValue("passwordConfirm")

	hash1 := hashPassword(password1, salt)
	hash2 := hashPassword(password2, salt)

	if !bytes.Equal(hash1, hash2) {
		fmt.Println("Register: passwords do not match")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)", username).Scan(&exists)
	if err != nil {
		fmt.Println(err)
	}
	if !exists {
		insertInDB(username, hash1, salt)
		fmt.Println("Register: insertion succesful!")
		w.WriteHeader(201)
	} else {
		fmt.Println("Register: there is already user with this username, aborting:", username)
		w.WriteHeader(409)
	}

}

type AuthRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func authorize(c *gin.Context) {
	var req AuthRequest

	fmt.Println("assdasd")
	// Декодируем JSON из тела запроса
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, "Invalid JSON")
		return
	}
	username := req.Login

	password := req.Password

	var userExists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)", username).Scan(&userExists)
	if err != nil {
		fmt.Println(err)
	}

	if !userExists {
		fmt.Println("Authorization: there is no user with such username:", username)
		c.String(http.StatusBadRequest, "There is no user")
		return
	}
	var correctHash []byte
	err = db.QueryRow("SELECT password_hash FROM users WHERE username = $1", username).Scan(&correctHash)
	if err != nil {
		fmt.Println(err)
	}

	var salt []byte
	err = db.QueryRow("SELECT salt FROM users WHERE username = $1", username).Scan(&salt)
	if err != nil {
		fmt.Println(err)
	}

	if bytes.Equal(hashPassword(password, salt), correctHash) {
		fmt.Println("Authorization: success!")
		c.String(http.StatusOK, "authorize success")
	} else {
		fmt.Println("Authorization: passwords do not match, aborting")
		c.String(http.StatusUnauthorized, "authorize failure")
	}
}

func generateSalt(length int) ([]byte, error) {
	salt := make([]byte, length)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

func hashPassword(password string, salt []byte) (hash []byte) {
	data := append([]byte(password), salt...)
	hash1 := sha256.Sum256(data)
	hash = hash1[:]
	return
}

func insertInDB(username string, password_hash []byte, salt []byte) {
	result, err := db.Exec("INSERT INTO users (username, password_hash, salt) VALUES ($1, $2, $3)", username, password_hash[:], salt)
	if err != nil {
		log.Println(err)
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		log.Println(err)
	}
	log.Println("Rows affected:", rowsAffected)

	lastInsertId, err := result.LastInsertId()

	if err != nil {
		log.Println(err)
	}
	log.Println("Last inserted id:", lastInsertId)

}

var db *sql.DB

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
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
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
