package user

import (
	"jamesgopsill/go-rest-template/internal/db"
	"jamesgopsill/go-rest-template/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type updatePasswordRequest struct {
	ID                 string `json:"id" binding:"required"`
	OldPassword        string `json:"oldPassword" binding:"required"`
	NewPassword        string `json:"newPassword" binding:"required"`
	ConfirmNewPassword string `json:"confirmNewPassword" binding:"required"`
}

func UpdatePassword(c *gin.Context) {

	var body updatePasswordRequest
	var err error

	if err = c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"data":  nil,
		})
		return
	}

	claims, ok := c.MustGet(gin.AuthUserKey).(*MyCustomClaims)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Auth pass-through problem.",
			"data":  nil,
		})
		return
	}

	if body.NewPassword != body.ConfirmNewPassword {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "New password and confirm new password do not match.",
			"data":  nil,
		})
		return
	}

	var user db.User
	res := db.Connection.First(&user, "id=?", body.ID)
	if res.Error != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "Account does not exist",
			"data":  nil,
		})
		return
	}

	if claims.ID != body.ID {

		err = bcrypt.CompareHashAndPassword(
			[]byte(user.PasswordHash),
			[]byte(body.OldPassword),
		)
		if err != nil {
			c.JSON(http.StatusNoContent, gin.H{
				"error": "Wrong password",
				"data":  nil,
			})
			return
		}

	}

	if claims.ID != body.ID || utils.Contains(claims.Scopes, db.ADMIN_SCOPE) {
		hash, err := bcrypt.GenerateFromPassword([]byte(body.NewPassword), bcrypt.MinCost)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"error": "Issue creating password",
				"data":  nil,
			})
			return
		}

		db.Connection.Model(&user).Update("PasswordHash", hash)
		c.JSON(http.StatusOK, gin.H{
			"error": nil,
			"data":  "success",
		})
		return
	}

	c.JSON(http.StatusBadGateway, gin.H{
		"error": "Should never get here.",
		"data":  nil,
	})

}
