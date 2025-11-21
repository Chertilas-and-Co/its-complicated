package handler

import (
	"bytes"
	"fmt"
	"main/backend/internal/pg"
	"main/backend/internal/auth/password"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
)

func registerHandler(
	w http.ResponseWriter,
	r *http.Request,
	_ httprouter.Params,
) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Form parse error", http.StatusBadRequest)
		return
	}

	username := r.FormValue("login")

	salt, _ := password.generateSalt(32)
	password1 := r.FormValue("password")
	password2 := r.FormValue("passwordConfirm")

	hash1 := password.hashPassword(password1, salt)
	hash2 := password.hashPassword(password2, salt)

	if !bytes.Equal(hash1, hash2) {
		fmt.Println("Register: passwords do not match")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var exists bool
	err := pg.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)", username).
		Scan(&exists)
	if err != nil {
		fmt.Println(err)
	}
	if !exists {
		pg.insertInDB(username, hash1, salt)
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
	err := pg.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)", username).
		Scan(&userExists)
	if err != nil {
		fmt.Println(err)
	}

	if !userExists {
		fmt.Println(
			"Authorization: there is no user with such username:",
			username,
		)
		c.String(http.StatusBadRequest, "There is no user")
		return
	}
	var correctHash []byte
	err = pg.db.QueryRow("SELECT password_hash FROM users WHERE username = $1", username).
		Scan(&correctHash)
	if err != nil {
		fmt.Println(err)
	}

	var salt []byte
	err = pg.db.QueryRow("SELECT salt FROM users WHERE username = $1", username).
		Scan(&salt)
	if err != nil {
		fmt.Println(err)
	}

	if bytes.Equal(password.hashPassword(password, salt), correctHash) {
		fmt.Println("Authorization: success!")
		c.String(http.StatusOK, "authorize success")
	} else {
		fmt.Println("Authorization: passwords do not match, aborting")
		c.String(http.StatusUnauthorized, "authorize failure")
	}
}
