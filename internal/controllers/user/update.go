package user

import (
	"jamesgopsill/go-rest-template/internal/db"
	"jamesgopsill/go-rest-template/internal/db/entities"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/utils"
)

type updateRequest struct {
	ID    string `json:"id" binding:"required"`
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

func Update(c *gin.Context) {

	var body updateRequest
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
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Auth pass-through problem.",
			"data":  nil,
		})
	}

	if claims.ID != body.ID && utils.Contains(claims.Scopes, db.ADMIN_SCOPE) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Need adminstrative rights.",
			"data":  nil,
		})
	}

	var user entities.User
	res := db.Connection.First(&user, "id=?", claims.ID)
	if res.Error != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "Account does not exist",
			"data":  nil,
		})
		return
	}

	db.Connection.Model(&user).Update("Name", body.Name)
	db.Connection.Model(&user).Update("Email", body.Email)

	c.JSON(http.StatusOK, gin.H{
		"error": nil,
		"data":  "success",
	})
}
