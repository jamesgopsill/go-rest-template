package main

import (
	"net/http"

	"jamesgopsill/go-rest-template/internal/db"
	"jamesgopsill/go-rest-template/internal/middleware"
	"jamesgopsill/go-rest-template/internal/user"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Info().Msg("Starting App")
	r := initialiseApp("tmp/test.db", gin.ReleaseMode)
	r.Run("localhost:3000")
}

func initialiseApp(dbPath string, mode string) *gin.Engine {
	db.Initialise(dbPath)

	gin.SetMode(mode)

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowHeaders = []string{"Authorization", "Content-Type"}

	r.Use(cors.New(config))

	r.GET("/ping", pong)
	r.POST("/user/register", user.Register)
	r.POST("/user/login", user.Login)
	r.POST("/user/update", middleware.AuthenticateMiddleware(db.USER_SCOPE), user.Update)

	return r
}

func pong(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"error": nil,
		"data":  "pong",
	})
}
