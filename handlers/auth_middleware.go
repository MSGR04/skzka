package handlers

import (
	"net/http"
	"strings"

	"hw2hard/store"

	"github.com/gin-gonic/gin"
)

const userCtxKey = "current_user"

func AuthMiddleware(users store.UserStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if !strings.HasPrefix(h, "Bearer ") {
			c.Header("WWW-Authenticate", `"Bearer realm="api"`)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		token := strings.TrimSpace(strings.TrimPrefix(h, "Bearer"))
		user, err := users.GetUserByToken(token)
		if err != nil {
			c.Header("WWW-Authenticate", `Bearer realm="api", error="invalid_token"`)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		c.Set(userCtxKey, user)
		c.Next()
	}
}

func CurrentUser(c *gin.Context) *store.User {
	if v, ok := c.Get(userCtxKey); ok {
		if u, ok := v.(*store.User); ok {
			return u
		}
	}
	return nil
}
