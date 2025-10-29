package store

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrNotFound = errors.New("not found")
	ErrNotReady = errors.New("not ready")
)

type IDGenerator interface{ New() string }
type TaskStatus string

const (
	StatusInProgress TaskStatus = "in_progress"
	StatusReady      TaskStatus = "ready"
)

type Task struct {
	ID        string     `json:"id"`
	Status    TaskStatus `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	Result    string     `json:"result"`
}

type TaskStore interface {
	CreateTask() (string, error)
	SetStatus(id string, status TaskStatus) error
	SetResult(id string, result string) error
	GetStatus(id string) (TaskStatus, error)
	GetResult(id string) (string, error)
}

type InMemoryTaskStore struct {
	mu    sync.RWMutex
	tasks map[string]Task
	idgen IDGenerator
}

func NewInMemoryTaskStore(gen IDGenerator) *InMemoryTaskStore {
	return &InMemoryTaskStore{tasks: make(map[string]Task), idgen: gen}
}

func (s *InMemoryTaskStore) CreateTask() (string, error) {
	id := s.idgen.New()
	now := time.Now()
	s.mu.Lock()
	s.tasks[id] = Task{ID: id, Status: StatusInProgress, CreatedAt: now, UpdatedAt: now}
	s.mu.Unlock()
	return id, nil
}
func (s *InMemoryTaskStore) SetStatus(id string, status TaskStatus) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	t, ok := s.tasks[id]
	if !ok {
		return ErrNotFound
	}
	t.Status = status
	t.UpdatedAt = time.Now()
	s.tasks[id] = t
	return nil
}
func (s *InMemoryTaskStore) SetResult(id string, result string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	t, ok := s.tasks[id]
	if !ok {
		return ErrNotFound
	}
	t.Result = result
	t.UpdatedAt = time.Now()
	s.tasks[id] = t
	return nil
}
func (s *InMemoryTaskStore) GetStatus(id string) (TaskStatus, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.tasks[id]
	if !ok {
		return "", ErrNotFound
	}
	return t.Status, nil
}

func (s *InMemoryTaskStore) GetResult(id string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.tasks[id]
	if !ok {
		return "", ErrNotFound
	}
	if t.Status != StatusReady {
		return "", ErrNotReady
	}
	return t.Result, nil
}
