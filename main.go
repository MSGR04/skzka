package main

import (
	"hw2hard/handlers"
	"log"
	"hw2hard/store"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type IDGenerator interface{ New() string }
type UUIDGen struct{}

func (UUIDGen) New() string { return uuid.NewString() }

func main() {
	sessionStore := store.NewInMemorySessionStore()

	go func() {
		ticker := time.NewTicker(30 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			sessionStore.GC(24 * time.Hour)
		}
	}()

	str := store.NewInMemoryTaskStore(UUIDGen{})
	users := store.NewInMemoryUserStore()

	r := gin.Default()

	handlers.RegisterTaskRoutes(r, str)
	handlers.RegisterUserRoutes(r, users)

	if err := r.Run(":8000"); err != nil {
		log.Fatal(err)
	}
}
