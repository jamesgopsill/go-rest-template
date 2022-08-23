package user

import (
	"io"
	"jamesgopsill/go-rest-template/internal/config"
	"jamesgopsill/go-rest-template/internal/db"
	"jamesgopsill/go-rest-template/internal/db/entities"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func UploadThumbnail(c *gin.Context) {

	thumbnail, header, err := c.Request.FormFile("thumbnail")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"data":  nil,
		})
		return
	}

	claims, ok := c.MustGet(gin.AuthUserKey).(*MyCustomClaims)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Auth pass-through problem.",
			"data":  nil,
		})
		return
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

	els := strings.Split(header.Filename, ".")
	filepath := config.UserThumbnailDir + "/" + claims.ID + "." + els[len(els)-1]

	if _, err := os.Stat(config.UserThumbnailDir); os.IsNotExist(err) {
		// create the dir
		if err := os.Mkdir(config.UserThumbnailDir, os.ModePerm); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
				"data":  nil,
			})
		}
	}

	out, err := os.Create(filepath)
	if err != nil {
		log.Info().Msg(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
			"data":  nil,
		})
		return
	}
	defer out.Close()
	_, err = io.Copy(out, thumbnail)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
			"data":  nil,
		})
		return
	}

	// update database
	db.Connection.Model(&user).Update("Thumbnail", claims.ID+"."+els[len(els)-1])

	// write and return result
	c.JSON(http.StatusOK, gin.H{
		"error": nil,
		"data":  "success",
	})

}
