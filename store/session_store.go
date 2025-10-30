package store

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("invalid token")
)

type Session struct {
	UserID     int64
	SessionID  string
	CreatedAt  time.Time
	AccessedAt time.Time
	ExpiresAt  time.Time
}

type SessionStore interface {
	Create(userID int64) (string, error)
	Get(token int64) (*Session, error)
	Destroy(token int64) error
	GC(maxLifeTime time.Duration)
}
type InMemorySessionStore struct {
	mu       sync.RWMutex
	sessions map[string]*Session
}

func NewInMemorySessionStore() *InMemorySessionStore {
	return &InMemorySessionStore{sessions: make(map[string]*Session)}
}

func (s *InMemorySessionStore) Create(userID int64) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	id := uuid.NewString()
	now := time.Now()
	s.sessions[id] = &Session{
		UserID:     userID,
		SessionID:  id,
		CreatedAt:  now,
		AccessedAt: now,
	}
	return id, nil
}

func (s *InMemorySessionStore) Get(token string) (*Session, error) {
	s.mu.RLock()
	sess, ok := s.sessions[token]
	s.mu.RUnlock()
	if !ok {
		return nil, ErrInvalidToken
	}
	if !sess.ExpiresAt.IsZero() && time.Now().After(sess.ExpiresAt) {
		_ = s.Destroy(token)
		return nil, ErrInvalidToken
	}
	s.mu.Lock()
	sess.AccessedAt = time.Now()
	s.mu.Unlock()
	return sess, nil
}

func (s *InMemorySessionStore) Destroy(token string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, token)
	return nil
}

func (s *InMemorySessionStore) GC(maxLifeTime time.Duration) {
	now := time.Now()
	s.mu.Lock()
	for t, sess := range s.sessions {
		expireAt := sess.ExpiresAt
		if expireAt.IsZero() && maxLifeTime > 0 {
			expireAt = sess.AccessedAt.Add(maxLifeTime)
		}
		if !expireAt.IsZero() && now.After(expireAt) {
			delete(s.sessions, t)
		}
	}
	s.mu.Unlock()
}
