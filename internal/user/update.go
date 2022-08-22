package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Update(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"error": nil,
		"data":  "success",
	})
}
