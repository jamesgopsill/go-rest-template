package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func RefreshToken(c *gin.Context) {

	claims, ok := c.MustGet(gin.AuthUserKey).(*MyCustomClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Auth pass-through problem.",
			"data":  nil,
		})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Oops, something has gone wrong",
			"data":  nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error": nil,
		"data":  "Bearer " + ss,
	})

}
