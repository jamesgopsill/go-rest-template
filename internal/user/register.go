package user

import (
	"jamesgopsill/go-rest-template/internal/db"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type registerRequest struct {
	Name            string `json:"name" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
	ConfirmEmail    string `json:"confirmEmail" binding:"required,email"`
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirmPassword" binding:"required"`
}

func Register(c *gin.Context) {

	var body registerRequest
	var err error

	if err = c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"data":  nil,
		})
		return
	}

	// TODO: Register the user.
	if body.Email != body.ConfirmEmail {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "Emails do not match",
			"data":  nil,
		})
		return
	}

	if body.Password != body.ConfirmPassword {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "Emails do not match",
			"data":  nil,
		})
		return
	}

	// Check if the user already exists
	var user db.User
	res := db.Connection.First(&user, "email=?", body.Email)
	if res.Error == nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "Account already exists",
			"data":  nil,
		})
		return
	}

	// log.Info().Msg("user does not exist")

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.MinCost)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "Issue creating password",
			"data":  nil,
		})
		return
	}

	var scopes []string
	scopes = append(scopes, db.USER_SCOPE)

	db.Connection.Create(&db.User{
		ID:           uuid.New().String(),
		Name:         body.Name,
		Email:        body.Email,
		PasswordHash: string(hash),
		Scopes:       scopes,
		Thumbnail:    "",
	})

	c.JSON(http.StatusOK, gin.H{
		"error": nil,
		"data":  "success",
	})
}
