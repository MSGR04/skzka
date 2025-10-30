package store

import (
	"errors"
	"sync"
)

var (
	ErrConflict = errors.New("conflict")
	ErrUnauth   = errors.New("unauthorized")
)

type User struct {
	ID       int64  `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserStore interface {
	Register(username, password string) (int64, error)
	Login(username, password string) (string, error)
	GetUserByToken(token string) (*User, error)
}

type InMemoryUserStore struct {
	mu           sync.RWMutex
	nextID       int64
	usersByLogin map[string]*User
	usersByID    map[int64]*User
	sessions     SessionStore
}

func NewInMemoryUserStore(sess SessionStore) *InMemoryUserStore {
	return &InMemoryUserStore{
		usersByLogin: make(map[string]*User),
		usersByID:    make(map[int64]*User),
		sessions:     sess,
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
	return u.sessions.Create(user.ID)
}

func (u *InMemoryUserStore) GetUserByToken(token string) (*User, error) {
	sess, err := u.sessions.Get(token)
	if err != nil {
		return nil, ErrUnauth
	}
	u.mu.RLock()
	defer u.mu.RUnlock()
	user := u.usersByID[sess.UserID]
	if user == nil {
		return nil, ErrUnauth
	}
	return user, nil
}
