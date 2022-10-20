package auth

import (
	"time"

	"github.com/golang-jwt/jwt"
)

type UserInfo struct {
	UserID   int    `json:"userId"`
	UserName string `json:"userName"`
	UserRole string `json:"userRole"`
}

type CustomJWTClaims struct {
	*jwt.StandardClaims
	UserInfo
}

func GenerateAccessToken(signature string, userDetails UserInfo) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &CustomJWTClaims{
		&jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
			// Audience:  "Nuttapon",
		},
		UserInfo{userDetails.UserID, userDetails.UserName, userDetails.UserRole},

		// UserInfo{"1234", "Nuttapon", "SuperAdmin"},
	})

	jwtToken, err := token.SignedString([]byte(signature))
	if err != nil {
		return "", err
	}

	return jwtToken, nil
}
