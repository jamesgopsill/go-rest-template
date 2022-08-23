package main

import (
	"net/http"

	"jamesgopsill/go-rest-template/internal/config"
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
	config.Initalise()
	db.Initialise()

	gin.SetMode(mode)

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowHeaders = []string{"Authorization", "Content-Type"}

	r.Use(cors.New(corsConfig))

	r.GET("/ping", pong)
	r.POST("/user/register", user.Register)
	r.POST("/user/login", user.Login)
	r.POST("/user/update", middleware.Authenticate(db.USER_SCOPE), user.Update)
	r.POST("/user/refresh-token", middleware.Authenticate(db.USER_SCOPE), user.RefreshToken)
	r.POST("/user/update-password", middleware.Authenticate(db.USER_SCOPE), user.UpdatePassword)
	r.POST("/user/upload-thumbnail", middleware.Authenticate(db.USER_SCOPE), user.UploadThumbnail)
	r.Static("/user/thumbnail", config.UserThumbnailDir)

	return r
}

func pong(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"error": nil,
		"data":  "pong",
	})
}
