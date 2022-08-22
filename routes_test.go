package main

import (
	"bytes"
	"encoding/json"
	"jamesgopsill/go-rest-template/internal/user"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang-jwt/jwt"
	"github.com/rs/zerolog/log"
)

var r *gin.Engine
var invalidSignedString string

const SECRET = "test"

// var validSignedString string

func init() {
	dbPath := "tmp/test.db"
	if _, err := os.Stat(dbPath); err == nil {
		err := os.Remove(dbPath)
		if err != nil {
			panic("Error")
		}
	}
	os.Setenv("GO_REST_JWT_SECRET", SECRET)
	os.Setenv("GO_REST_JWT_ISSUER", "www.test.com")

	var invalidScopes []string
	invalidClaims := user.MyCustomClaims{
		Name:   "a",
		Email:  "b",
		Scopes: invalidScopes,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: 15000,
			Issuer:    "www.test.com",
		},
	}

	invalidToken := jwt.NewWithClaims(jwt.SigningMethodHS256, invalidClaims)
	invalidSignedString, _ = invalidToken.SignedString(SECRET)

	/*
		var validScopes []string
		validScopes = append(validScopes, db.ADMIN_SCOPE)
		validScopes = append(validScopes, db.USER_SCOPE)
		claims := user.MyCustomClaims{
			Name:   "a",
			Email:  "b",
			Scopes: validScopes,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: 15000,
				Issuer:    "www.test.com",
			},
		}

		validToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		validSignedString, _ = validToken.SignedString("test")
	*/

	r = initialiseApp(dbPath, gin.ReleaseMode)
}

func TestPing(t *testing.T) {
	req, _ := http.NewRequest("GET", "/ping", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	log.Info().Msg(w.Body.String())
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRegister(t *testing.T) {
	mockRequest := `{
		"name": "test",
		"email" : "test@test.com",
		"confirmEmail" : "test@test.com",
		"password" : "test",
		"confirmPassword" : "test"
	}`
	req, _ := http.NewRequest("POST", "/user/register", bytes.NewBufferString(mockRequest))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRegisterAccountExists(t *testing.T) {
	mockRequest := `{
		"name": "test",
		"email" : "test@test.com",
		"confirmEmail" : "test@test.com",
		"password" : "test",
		"confirmPassword" : "test"
	}`
	req, _ := http.NewRequest("POST", "/user/register", bytes.NewBufferString(mockRequest))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnprocessableEntity {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestLogin(t *testing.T) {
	mockRequest := `{
		"password": "test",
		"email": "test@test.com"
	}`
	req, _ := http.NewRequest("POST", "/user/login", bytes.NewBufferString(mockRequest))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthMiddlewareInvalidToken(t *testing.T) {
	mockRequest := `{}`
	req, _ := http.NewRequest("POST", "/user/update", bytes.NewBufferString(mockRequest))
	req.Header.Set("Authorization", "Bearer "+invalidSignedString)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

type loginResponse struct {
	Error string
	Data  string
}

func TestUpdateUser(t *testing.T) {
	var mockRequest string
	mockRequest = `{
		"password": "test",
		"email": "test@test.com"
	}`
	req, _ := http.NewRequest("POST", "/user/login", bytes.NewBufferString(mockRequest))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusOK, w.Code)

	var response loginResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		panic("error")
	}
	els := strings.Split(response.Data, " ")

	token, _ := jwt.ParseWithClaims(els[1], &user.MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET), nil
	})
	claims, ok := token.Claims.(*user.MyCustomClaims)
	if !ok {
		panic("Error in claims")
	}

	mockRequest = `{
		"id": "` + claims.ID + `",
		"name": "updated name",
		"email": "test@test.com"
	}`
	req, _ = http.NewRequest("POST", "/user/update", bytes.NewBufferString(mockRequest))
	req.Header.Set("Authorization", response.Data)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRefreshToken(t *testing.T) {
	mockRequest := `{
		"password": "test",
		"email": "test@test.com"
	}`
	req, _ := http.NewRequest("POST", "/user/login", bytes.NewBufferString(mockRequest))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusOK, w.Code)

	var response loginResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		panic("error")
	}

	req, _ = http.NewRequest("POST", "/user/refresh-token", bytes.NewBufferString(mockRequest))
	req.Header.Set("Authorization", response.Data)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusOK, w.Code)
}
