package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func AuthMiddleware() gin.HandlerFunc {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev-secret-key"
	}

	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			respondUnauthorized(c, "Missing Authorization header")
			return
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			respondUnauthorized(c, "Invalid Authorization header format")
			return
		}

		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			respondUnauthorized(c, "Invalid or expired token")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			respondUnauthorized(c, "Invalid token claims")
			return
		}

		userIDRaw, ok := claims["user_id"]
		if !ok {
			respondUnauthorized(c, "Invalid token payload")
			return
		}

		userIDFloat, ok := userIDRaw.(float64)
		if !ok {
			respondUnauthorized(c, "Invalid token payload")
			return
		}

		c.Set("user_id", uint(userIDFloat))
		c.Next()
	}
}

func respondUnauthorized(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"code":    40102,
		"message": message,
	})
}
