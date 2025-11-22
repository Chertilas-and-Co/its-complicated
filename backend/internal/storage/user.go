package storage

import (
	"errors"
	"its-complicated/internal/models"
	"sync"
)

var (
	ErrUserExists = errors.New("user already exists")
)

// UserStorage provides an in-memory storage for users.
type UserStorage struct {
	mu      sync.Mutex
	users   map[string]models.User
	counter int
}

// NewUserStorage creates and returns a new UserStorage.
func NewUserStorage() *UserStorage {
	return &UserStorage{
		users:   make(map[string]models.User),
		counter: 1,
	}
}

// CreateUser adds a new user to the storage.
// It returns an error if a user with the same username already exists.
func (s *UserStorage) CreateUser(username, password string) (*models.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.users[username]; exists {
		return nil, ErrUserExists
	}

	user := models.User{
		ID:       s.counter,
		Username: username,
		Password: password, // In a real app, hash the password!
	}

	s.users[username] = user
	s.counter++

	return &user, nil
}
