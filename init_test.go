package main

import (
	"jamesgopsill/go-rest-template/internal/controllers/user"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

var r *gin.Engine
var invalidSignedString string

const SECRET = "test"

func init() {
	dbPath := "data/test.db"
	issuer := "www.test.com"
	UserThumbnailDir := "data/user-thumbnails"
	if _, err := os.Stat(dbPath); err == nil {
		err := os.Remove(dbPath)
		if err != nil {
			panic("Error")
		}
	}
	os.Setenv("GO_REST_JWT_SECRET", SECRET)
	os.Setenv("GO_REST_DB_PATH", dbPath)
	os.Setenv("GO_REST_JWT_ISSUER", issuer)
	os.Setenv("GO_REST_USER_THUMBNAIL_DIR", UserThumbnailDir)
	os.RemoveAll(UserThumbnailDir)

	var invalidScopes []string
	invalidClaims := user.MyCustomClaims{
		Name:   "a",
		Email:  "b",
		Scopes: invalidScopes,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Unix() + 24*60*60,
			Issuer:    "www.test.com",
		},
	}

	invalidToken := jwt.NewWithClaims(jwt.SigningMethodHS256, invalidClaims)
	invalidSignedString, _ = invalidToken.SignedString(SECRET)

	r = initialiseApp(dbPath, gin.ReleaseMode)
}
