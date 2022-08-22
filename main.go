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

	log.Info().Msg("Hello World")

	r := initialiseApp()

	r.Run("localhost:3000")
}

func initialiseApp() *gin.Engine {
	db.Initialise("tmp/test.db")

	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowHeaders = []string{"Authorization", "Content-Type"}

	r.Use(cors.New(config))

	r.GET("/ping", pong)
	r.GET("/user/register", user.Register)
	r.GET("/user/update", middleware.AuthenticateMiddleware("test"), user.Update)
	r.GET("/user/login", middleware.AuthenticateMiddleware(""), user.Login)

	return r
}

func pong(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"error": nil,
		"data":  "pong",
	})
}
