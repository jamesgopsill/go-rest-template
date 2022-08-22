package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func AuthenticateMiddleware(scopes string) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		log.Info().Msg("Auth Fired")
		log.Info().Msg(scopes)
		if scopes == "test" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Not Authorised",
				"data":  nil,
			})
		}
		c.Set(gin.AuthUserKey, "world")
	}
	return gin.HandlerFunc(fn)
}
