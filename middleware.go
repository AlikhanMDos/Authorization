package main

import (
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware функция middleware для аутентификации через JWT.
//
//	@Summary		Middleware для проверки авторизации через JWT
//	@Description	Middleware проверяет наличие и валидность JWT в заголовке Authorization.
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Security		BearerToken
//	@Param			Authorization	header		string				true	"Bearer JWT token"
//	@Success		200				{string}	string				"Authorized"
//	@Failure		401				{object}	map[string]string	"Unauthorized"
//	@Router			/auth [get]
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		// Check if token is blacklisted
		mutex.Lock()
		_, blacklisted := tokenBlacklist[tokenString]
		mutex.Unlock()
		if blacklisted {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is blacklisted (logged out)"})
			c.Abort()
			return
		}

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("phone", claims.Phone)
		c.Next()
	}
}
