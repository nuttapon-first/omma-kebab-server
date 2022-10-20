package middleware

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/nuttapon-first/omma-kebab-server/pkg/utils"
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
			userRole := claims["userRole"]
			fmt.Println(claims)
			c.Set("userId", userId)
			c.Set("userName", userName)
			c.Set("userRole", userRole)
		}
		c.Next()
	}
}

func Authorization(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		rolesMap := map[string]int{
			"user":    1, // 2^0 // 001
			"manager": 2, // 2^1 // 010
			"admin":   4, // 2^2 // 100
		}

		var sum int
		for i := range roles {
			sum += rolesMap[roles[i]]
		}

		if sum == 0 && len(roles) == 0 {
			c.Next()
		} else if sum == 0 && len(roles) != 0 {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Invalid role authorization",
			})
			return
		}

		roleBinary := utils.BinaryConvertor(sum, len(rolesMap))

		userRole, _ := c.Get("userRole")
		userRoleNumber := rolesMap[userRole.(string)]
		userRoleBinary := utils.BinaryConvertor(userRoleNumber, len(rolesMap))

		for i := 0; i < len(rolesMap); i++ {
			if roleBinary[i]&userRoleBinary[i] == 1 {
				c.Next()
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": "Permission denied",
		})
	}
}
