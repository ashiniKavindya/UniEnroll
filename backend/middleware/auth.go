package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v4"

	"backend/models"
)

// AuthMiddleware validates JWT token from Authorization header
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			secret = "devsecret"
		}

		token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			c.Abort()
			return
		}

		// Extract user ID and role from claims
		userID, _ := (*claims)["sub"].(string)
		role, _ := (*claims)["role"].(string)

		// Set user info in context for use in handlers
		c.Set("userID", userID)
		c.Set("role", role)

		// Enforce password change if required (allow change-password endpoint)
		if c.Request.URL.Path != "/api/change-password" {
			f, ok := c.MustGet("firestore").(*firestore.Client)
			if ok && f != nil {
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				doc, err := f.Collection("users").Doc(userID).Get(ctx)
				if err == nil {
					var u models.User
					if err := doc.DataTo(&u); err == nil && u.ForcePasswordChange {
						c.JSON(http.StatusForbidden, gin.H{"error": "password change required"})
						c.Abort()
						return
					}
				}
			}
		}

		c.Next()
	}
}

// RequireRole checks if user has specific role
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user role not found"})
			c.Abort()
			return
		}

		userRole, ok := role.(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user role"})
			c.Abort()
			return
		}

		allowed := false
		for _, r := range allowedRoles {
			if userRole == r {
				allowed = true
				break
			}
		}

		if !allowed {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}
