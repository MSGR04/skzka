package main

import (
	"hw2hard/handlers"
	"hw2hard/store"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type IDGenerator interface{ New() string }
type UUIDGen struct{}

func (UUIDGen) New() string { return uuid.NewString() }

func main() {

	str := store.NewInMemoryTaskStore(UUIDGen{})
	users := store.NewInMemoryUserStore()

	r := gin.Default()

	handlers.RegisterTaskRoutes(r, str)
	handlers.RegisterUserRoutes(r, users)

	if err := r.Run(":8000"); err != nil {
		log.Fatal(err)
	}
}
