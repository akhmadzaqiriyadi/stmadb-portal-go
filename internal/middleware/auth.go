// internal/middleware/auth.go
package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"

	// PERBAIKAN 1: Tambahkan import ini
	"github.com/akhmadzaqiriyadi/stmadb-portal-go/prisma/db"
)

// Authenticate adalah middleware untuk memvalidasi token JWT
func Authenticate(dbClient *db.PrismaClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be 'Bearer <token>'"})
			return
		}

		tokenString := parts[1]
		jwtSecret := viper.GetString("JWT_SECRET")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				// PERBAIKAN 3: Gunakan errors.New
				return nil, errors.New(fmt.Sprintf("Unexpected signing method: %v", token.Header["alg"]))
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		userIdFloat, ok := claims["userId"].(float64)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
			return
		}
		userId := int(userIdFloat)

		user, err := dbClient.User.FindUnique(
			// Konversi userId ke types.BigInt
			db.User.ID.Equals(db.BigInt(userId)),
		).Exec(context.Background())

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			return
		}

		c.Set("user", user)
		c.Next()
	}
}

// Authorize adalah middleware untuk memeriksa peran pengguna.
func Authorize(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userCtx, exists := c.Get("user")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
			return
		}

		user := userCtx.(*db.UserModel)
		isAllowed := false
		for _, role := range allowedRoles {
			if string(user.Role) == role {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "You do not have permission to access this resource"})
			return
		}

		c.Next()
	}
}