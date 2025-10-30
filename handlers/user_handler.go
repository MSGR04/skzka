package handlers

import (
	"errors"
	"hw2hard/store"
	"net/http"

	"github.com/gin-gonic/gin"
)

var ErrConflict = errors.New("conflict")

type RegisterRequest struct {
	Username string `json:"username" example:"user_123"`
	Password string `json:"password" example:"password228"`
}

type LoginRequest struct {
	Username string `json:"username" example:"user_123"`
	Password string `json:"password" example:"password228"`
}

func RegisterUserRoutes(r *gin.Engine, users store.UserStore) {
	r.POST("/register", registerHandler(users))
	r.POST("/login", loginHandler(users))
}

func registerHandler(users store.UserStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RegisterRequest
		if err := c.BindJSON(&req); err != nil || req.Username == "" || req.Password == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
			return
		}
		_, err := users.Register(req.Username, req.Password)
		if err != nil {
			if errors.Is(err, store.ErrConflict) {
				c.Status(http.StatusCreated)
				return
			}
			if err == store.ErrConflict {
				c.JSON(http.StatusConflict, gin.H{"error": "user exists"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register"})
			return
		}
		c.Status(http.StatusCreated)
	}
}
func loginHandler(users store.UserStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginRequest
		if err := c.BindJSON(&req); err != nil || req.Username == "" || req.Password == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
			return
		}
		token, err := users.Login(req.Username, req.Password)
		if err != nil {
			c.Header("WWW-Authenticate", `Bearer realm="api", error="invalid_credentials"`)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": token})
	}
}
