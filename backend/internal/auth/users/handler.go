package users

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"main/internal/auth/password"
	passwd "main/internal/auth/password"
	"main/internal/pg"

	"net/http"
)

type RegisterRequest struct {
	Login           string `json:"login"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"passwordConfirm"`
}

func RegisterHandler(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, "Invalid JSON")
		return
	}

	login := req.Login
	email := req.Email

	salt, _ := password.GenerateSalt(32)
	password1 := req.Password
	password2 := req.PasswordConfirm
	fmt.Println(login, salt, password1, password2)

	hash1 := password.HashPassword(password1, salt)
	hash2 := password.HashPassword(password2, salt)

	if !bytes.Equal(hash1, hash2) {
		fmt.Println("Register: passwords do not match")
		c.String(http.StatusBadRequest, "Invalid JSON")
		return
	}

	var exists bool
	err := pg.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", email).Scan(&exists)
	if err != nil {
		fmt.Println(err)
	}
	if !exists {
		pg.InsertInDB(login, email, hash1, salt)
		fmt.Println("Register: insertion succesful!")
		c.String(http.StatusCreated, "Register succesful")
	} else {
		fmt.Println("Register: there is already user with this email, aborting:", email)
		c.String(http.StatusConflict, "Register failed, user already exists")
	}

}

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func AuthorizeUser(c *gin.Context) {
	var req AuthRequest

	fmt.Println("assdasd")
	// Декодируем JSON из тела запроса
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, "Invalid JSON")
		return
	}
	username := req.Email

	password := req.Password
	fmt.Println(username)
	var userExists bool
	err := pg.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)", username).
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
	err = pg.DB.QueryRow("SELECT password_hash FROM users WHERE username = $1", username).
		Scan(&correctHash)
	if err != nil {
		fmt.Println(err)
	}

	var salt []byte
	err = pg.DB.QueryRow("SELECT salt FROM users WHERE username = $1", username).
		Scan(&salt)
	if err != nil {
		fmt.Println(err)
	}

	if bytes.Equal(passwd.HashPassword(password, salt), correctHash) {
		fmt.Println("Authorization: success!")
		c.String(http.StatusOK, "authorize success")
	} else {
		fmt.Println("Authorization: passwords do not match, aborting")
		c.String(http.StatusUnauthorized, "authorize failure")
	}
}
