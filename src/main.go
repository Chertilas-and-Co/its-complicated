package main

import (
	"database/sql"
	"fmt"
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

	password := r.FormValue("password")
	passwordConfirm := r.FormValue("passwordConfirm")
	fmt.Println(password, " ", passwordConfirm)

	if password != passwordConfirm {
		http.Error(w, "passwords do not match", http.StatusBadRequest)
		return
	}

	w.Write([]byte("Registration success!"))
}

func authorize(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Form parse error", http.StatusBadRequest)
		return
	}

	password := r.FormValue("password")

	w.Write([]byte(password))
}

var db *sql.DB

func main() {
	// Загрузить .env (ошибка обработки роли не играет, можно логировать)
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

	router := httprouter.New()
	router.ServeFiles("/static/*filepath", http.Dir("static"))

	router.GET("/auth", handleLogger(getAuthPage))
	router.GET("/register", handleLogger(getRegPage))
	router.NotFound = http.HandlerFunc(handlerFuncLogger(notFoundHandler))
	router.GET("/hello/:name", handleLogger(hello))
	router.POST("/register", handleLogger(registerHandler))

	fmt.Println("Server is listening at :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
