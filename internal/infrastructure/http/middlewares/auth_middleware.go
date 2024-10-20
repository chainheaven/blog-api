package middlewares

import (
	"blog-api/internal/application/user"
	"blog-api/internal/infrastructure/auth"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	bearerSchema         = "Bearer "
	authorizationHeader  = "Authorization"
	userIDKey            = "userID"
	passwordChangedAtKey = "passwordChangedAt"
)

var (
	errMissingAuthHeader = "Authorization header is required"
	errInvalidToken      = "Invalid or expired token"
)

// AuthMiddleware 返回一個 Gin 中間件，用於驗證 JWT 令牌
func AuthMiddleware(jwtService *auth.JWTService, userService *user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("Starting AuthMiddleware")
		authHeader := c.GetHeader(authorizationHeader)

		if authHeader == "" {
			log.Println("Missing Authorization header")
			c.JSON(http.StatusUnauthorized, gin.H{"error": errMissingAuthHeader})
			c.Abort()
			return
		}

		tokenString := extractToken(authHeader)
		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			log.Printf("Token validation error: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": errInvalidToken})
			c.Abort()
			return
		}

		// 獲取用戶當前的資料
		currentUser, err := userService.GetUserProfile(claims.UserID)
		if err != nil {
			log.Printf("Failed to get user profile: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		// 比較 token 中的密碼更改時間與用戶當前的密碼更改時間
		if claims.PasswordChangedAt.Before(currentUser.PasswordChangedAt) {
			log.Printf("Token expired due to password change. Token time: %v, Current time: %v", claims.PasswordChangedAt, currentUser.PasswordChangedAt)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired due to password change"})
			c.Abort()
			return
		}

		c.Set(userIDKey, claims.UserID)
		c.Set(passwordChangedAtKey, claims.PasswordChangedAt)
		log.Printf("User authenticated: %d, Password changed at: %v", claims.UserID, claims.PasswordChangedAt)

		c.Next()
		log.Println("Finished AuthMiddleware")
	}
}

// extractToken 從Header中提取 token
func extractToken(authHeader string) string {
	if strings.HasPrefix(authHeader, bearerSchema) {
		return authHeader[len(bearerSchema):]
	}
	return authHeader
}

// GetUserID 從 Gin 上下文中獲取已認證的用戶 ID
func GetUserID(c *gin.Context) (uint, error) {
	userID, exists := c.Get(userIDKey)
	if !exists {
		return 0, fmt.Errorf("user ID not found in context")
	}

	id, ok := userID.(uint)
	if !ok {
		return 0, fmt.Errorf("user ID is not of type uint")
	}

	return id, nil
}

// GetPasswordChangedAt 從 Gin 上下文中獲取密碼修改時間
func GetPasswordChangedAt(c *gin.Context) (time.Time, error) {
	passwordChangedAt, exists := c.Get(passwordChangedAtKey)
	if !exists {
		return time.Time{}, fmt.Errorf("password changed at not found in context")
	}

	changedAt, ok := passwordChangedAt.(time.Time)
	if !ok {
		return time.Time{}, fmt.Errorf("password changed at is not of type time.Time")
	}

	return changedAt, nil
}
