package middleware

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func JwtGuard() gin.HandlerFunc {
	return func(c *gin.Context) {
		signature := []byte(os.Getenv("JWT_SECRET"))
		tokenStr := c.Request.Header.Get("Access-Token")
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return signature, nil
		})

		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			userId := claims["userId"]
			userName := claims["userName"]
			c.Set("userId", userId)
			c.Set("userName", userName)
		}
		c.Next()
	}
}
