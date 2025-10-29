package store

import (
	"errors"
	"sync"

	"github.com/google/uuid"
)

var (
	ErrConflict     = errors.New("conflict")
	ErrUnauth       = errors.New("unauthorized")
	ErrInvalidToken = errors.New("invalid token")
)

type User struct {
	ID       int64  `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Session struct {
	UserID    int64  `json:"user_id"`
	SessionID string `json:"session_id"`
}

type UserStore interface {
	Register(username, password string) error
	Login(username, password string) (string, error)
	GetUserByToken(token string) (*User, error)
}

type InMemoryUserStore struct {
	mu           sync.RWMutex
	nextID       int64
	usersByLogin map[string]*User
	usersByID    map[int64]*User
	sessions     map[string]*Session
}

func NewInMemoryUserStore() *InMemoryUserStore {
	return &InMemoryUserStore{
		usersByLogin: make(map[string]*User),
		usersByID:    make(map[int64]*User),
		sessions:     make(map[string]*Session),
	}
}

func (u *InMemoryUserStore) Register(username, password string) (int64, error) {
	u.mu.Lock()
	defer u.mu.Unlock()

	if _, exists := u.usersByLogin[username]; exists {
		return 0, ErrConflict
	}

	u.nextID++
	user := &User{
		ID:       u.nextID,
		Login:    username,
		Password: password,
	}
	u.usersByLogin[username] = user
	u.usersByID[user.ID] = user
	return user.ID, nil
}

func (u *InMemoryUserStore) Login(username, password string) (string, error) {
	u.mu.Lock()
	defer u.mu.Unlock()

	user, ok := u.usersByLogin[username]
	if !ok || user.Password != password {
		return "", ErrUnauth
	}
	token := uuid.NewString()
	u.sessions[token] = &Session{
		UserID:    user.ID,
		SessionID: token,
	}
	return token, nil
}

func (u *InMemoryUserStore) GetUserByToken(token string) (*User, error) {
	u.mu.RLock()
	defer u.mu.RUnlock()

	sess, ok := u.sessions[token]
	if !ok {
		return nil, ErrInvalidToken
	}
	user, ok := u.usersByID[sess.UserID]
	if !ok {
		return nil, ErrInvalidToken
	}
	return user, nil
}
