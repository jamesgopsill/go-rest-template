package config

import "os"

var DBPath string
var Secret string
var MySigningKey []byte
var Issuer string
var UserThumbnailDir string
var HostWhiteList string

func Initalise() {

	Secret := os.Getenv("GO_REST_JWT_SECRET")
	if Secret == "" {
		panic("No secret.")
	}
	MySigningKey = []byte(Secret)

	Issuer = os.Getenv("GO_REST_JWT_ISSUER")
	if Issuer == "" {
		panic("No issuer.")
	}

	DBPath = os.Getenv("GO_REST_DB_PATH")
	if DBPath == "" {
		panic("No DB Path.")
	}

	UserThumbnailDir = os.Getenv("GO_REST_USER_THUMBNAIL_DIR")
	if UserThumbnailDir == "" {
		panic("No User Thumbnail Dir.")
	}

	HostWhiteList = os.Getenv("GO_REST_HOST_WHITE_LIST")

}
