package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"
	"todos_api/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

/*
AuthMiddleware is a Gin middleware that protects routes by validating JWT tokens.

This middleware performs authentication by:

  1. Reading the Authorization header
  2. Extracting the Bearer token
  3. Verifying the token signature using the server's JWT secret
  4. Validating token expiration
  5. Extracting the user_id claim
  6. Storing user_id in Gin context for downstream handlers

If authentication fails at any step, the request is rejected with HTTP 401.

Parameters:
  cfg - Application configuration containing JWTSecret used for token verification

Returns:
  gin.HandlerFunc - Middleware function compatible with Gin router

Usage example:

  router.GET("/todos",
      AuthMiddleware(cfg),
      handlers.GetAllTodosHandler(pool),
  )

Authorization Header Format:

  Authorization: Bearer <JWT_TOKEN>

Example:

  Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

Context values set:

  "user_id" - ID of authenticated user

Downstream handlers can retrieve it using:

  userID := c.Get("user_id")

Authentication Flow:

  Client Login → Receive JWT
  Client Request → Send JWT in Authorization header
  Middleware → Validate JWT
  Middleware → Attach user_id to context
  Handler → Execute authorized logic
*/
func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		if tokenString == "" || tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(cfg.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)

		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token Claims"})
			c.Abort()
			return
		}

		userID, ok := claims["user_id"].(string)

		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token Claims"})
			c.Abort()
			return
		}

		if exp, ok := claims["exp"].(float64); ok {
			expirationTime := time.Unix(int64(exp), 0)

			if time.Now().After(expirationTime) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has expired"})
				c.Abort()
				return
			}
		}

		c.Set("user_id", userID)
		c.Next()
	}
}
