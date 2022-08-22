package user

import (
	"jamesgopsill/go-rest-template/internal/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	//user := c.MustGet(gin.AuthUserKey).(string)
	var user db.User
	db.Connection.First(&user, 1)
	c.JSON(http.StatusOK, gin.H{
		"error": nil,
		"data":  user,
	})
}
