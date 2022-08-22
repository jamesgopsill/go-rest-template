package user

import (
	"jamesgopsill/go-rest-template/internal/db"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type MyCustomClaims struct {
	ID     string   `json:"id"`
	Name   string   `json:"name"`
	Email  string   `json:"email"`
	Scopes []string `json:"scopes"`
	jwt.StandardClaims
}

func Login(c *gin.Context) {

	var body loginRequest
	var err error

	if err = c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"data":  nil,
		})
		return
	}

	// Make sure some details are present
	if body.Email == "" || body.Password == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Please provide valid login details",
			"data":  nil,
		})
		return
	}

	// Check if the user exists
	var user db.User
	res := db.Connection.First(&user, "email=?", body.Email)
	if res.Error != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "Account does not exist",
			"data":  nil,
		})
		return
	}

	// Check password
	if err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(body.Password),
	); err != nil {
		c.JSON(http.StatusNoContent, gin.H{
			"error": "Wrong password",
			"data":  nil,
		})
		return
	}

	secret := os.Getenv("GO_REST_JWT_SECRET")
	if secret == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Oops, something has gone wrong",
			"data":  nil,
		})
		return
	}
	mySigningKey := []byte(secret)

	issuer := os.Getenv("GO_REST_JWT_ISSUER")
	if issuer == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Oops, something has gone wrong",
			"data":  nil,
		})
		return
	}

	// Create the jwt with the claims
	claims := MyCustomClaims{
		ID:     user.ID,
		Name:   user.Name,
		Email:  user.Email,
		Scopes: user.Scopes,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Unix() + 60*60*24,
			Issuer:    issuer,
		},
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
