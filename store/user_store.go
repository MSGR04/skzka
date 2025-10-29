package store

import (
	"errors"
	"sync"

	"github.com/google/uuid"
)

var (
	ErrConflict = errors.New("conflict")
	ErrUnauth   = errors.New("unauthorized")
)

type UserStore interface {
	Register(username, password string) error
	Login(username, password string) (string, error)
}

type InMemoryUserStore struct {
	mu       sync.RWMutex
	users    map[string]string
	sessions map[string]string
}

func NewInMemoryUserStore() *InMemoryUserStore {
	return &InMemoryUserStore{
		users:    make(map[string]string),
		sessions: make(map[string]string),
	}
}

func (u *InMemoryUserStore) Register(username, password string) error {
	u.mu.Lock()
	defer u.mu.Unlock()
	if _, ok := u.users[username]; ok {
		return ErrConflict
	}
	u.users[username] = password
	return nil
}
func (u *InMemoryUserStore) Login(username, password string) (string, error) {
	u.mu.Lock()
	defer u.mu.Unlock()
	pw, ok := u.users[username]
	if !ok || pw != password {
		return "", ErrUnauth
	}
	token := uuid.NewString()
	u.sessions[token] = username
	return token, nil
}
