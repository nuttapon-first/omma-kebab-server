package login

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/nuttapon-first/omma-kebab-server/modules/auth"
	"github.com/nuttapon-first/omma-kebab-server/modules/model"
	"github.com/nuttapon-first/omma-kebab-server/store"
	"golang.org/x/crypto/bcrypt"
)

// ////////////////////////////////////////////////////////////////////
// SPI

type LoginHandler struct {
	store store.Storer
}

func NewLoginHandler(store store.Storer) *LoginHandler {
	return &LoginHandler{store: store}
}

// ////////////////////////////////////////////////////////////////////

func (h *LoginHandler) Login(c *gin.Context) {
	login := &model.LoginUser{}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		c.JSON(http.StatusNotFound, map[string]interface{}{
			"success": 1,
			"error":   "Jwt secret not found",
		})
		return
	}

	if err := c.ShouldBindJSON(login); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": 2,
			"error":   err.Error(),
		})
		return
	}

	user := &model.User{}
	err := h.store.Table("Users").Where(map[string]interface{}{"user_name": login.UserName}).First(user).Error
	if err != nil {
		if err.Error() == "record not found" {
			c.JSON(http.StatusNotFound, map[string]interface{}{
				"success": 1,
				"error":   "User not found",
			})
			return
		} else {
			c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"success": 5,
				"error":   err.Error(),
			})
			return
		}
	}

	userCredential := &model.UserCredential{}
	err = h.store.Table("UserCredentials").Where(map[string]interface{}{"user_id": user.ID}).First(userCredential).Error
	if err != nil {
		if err.Error() == "record not found" {
			c.JSON(http.StatusNotFound, map[string]interface{}{
				"success": 1,
				"error":   "User Credential not found",
			})
			return
		} else {
			c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"success": 5,
				"error":   err.Error(),
			})
			return
		}
	}

	// Comparing the password with the hash
	err = bcrypt.CompareHashAndPassword([]byte(userCredential.Credential), []byte(login.UserPassword))
	if err != nil {
		c.JSON(http.StatusForbidden, map[string]interface{}{
			"success": 3,
			"error":   "Invalid password",
		})
		return
	}

	userInfo := auth.UserInfo{
		UserID:   user.ID,
		UserName: user.UserFullName,
	}

	token, err := auth.GenerateAccessToken(secret, userInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": 5,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"success": 0,
		"token":   token,
	})
}
