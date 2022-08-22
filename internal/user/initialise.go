package user

import (
	"os"
)

var mySigningKey []byte
var issuer string

func Initialise() {
	secret := os.Getenv("GO_REST_JWT_SECRET")
	if secret == "" {
		panic("No secret.")
	}
	mySigningKey = []byte(secret)

	issuer = os.Getenv("GO_REST_JWT_ISSUER")
	if issuer == "" {
		panic("No issuer.")
	}
}
