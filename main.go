package main

import (
	"log"
	"time"

	"hw2hard/handlers"
	"hw2hard/store"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// интерфейс генератора ID для задач
type IDGenerator interface{ New() string }

// реализация генерации UUID
type UUIDGen struct{}

func (UUIDGen) New() string { return uuid.NewString() }

func main() {
	// === 1. Инициализируем хранилища ===

	// Хранилище сессий (аналог session manager)
	sessionStore := store.NewInMemorySessionStore()

	// Запускаем фоновый GC для очистки протухших сессий
	go func() {
		ticker := time.NewTicker(30 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			sessionStore.GC(24 * time.Hour)
		}
	}()

	// Хранилище пользователей, использующее SessionStore
	userStore := store.NewInMemoryUserStore(sessionStore)

	// Хранилище задач (ID генерируем через UUID)
	taskStore := store.NewInMemoryTaskStore(UUIDGen{})

	// === 2. Создаём HTTP-сервер ===
	r := gin.Default()

	// публичные роуты
	handlers.RegisterUserRoutes(r, userStore)

	// защищённые роуты (требуют Authorization: Bearer ...)
	handlers.RegisterTaskRoutes(r, taskStore, userStore)

	// === 3. Запускаем сервер ===
	log.Println("Server running on http://localhost:8000")
	if err := r.Run(":8000"); err != nil {
		log.Fatal(err)
	}
}
