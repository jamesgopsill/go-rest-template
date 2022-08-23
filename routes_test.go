package main

import (
	"bytes"
	"encoding/json"
	"io"
	"jamesgopsill/go-rest-template/internal/controllers/user"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func TestPing(t *testing.T) {
	req, err := http.NewRequest("GET", "/ping", nil)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		log.Info().Msg(w.Body.String())
	}
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
	req, err := http.NewRequest("POST", "/user/register", bytes.NewBufferString(mockRequest))
	assert.NoError(t, err)
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
	req, err := http.NewRequest("POST", "/user/register", bytes.NewBufferString(mockRequest))
	assert.NoError(t, err)
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
	req, err := http.NewRequest("POST", "/user/login", bytes.NewBufferString(mockRequest))
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthMiddlewareInvalidToken(t *testing.T) {
	mockRequest := `{}`
	req, err := http.NewRequest("POST", "/user/update", bytes.NewBufferString(mockRequest))
	assert.NoError(t, err)
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
	req, err := http.NewRequest("POST", "/user/login", bytes.NewBufferString(mockRequest))
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusOK, w.Code)

	var response loginResponse
	err = json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	els := strings.Split(response.Data, " ")

	token, err := jwt.ParseWithClaims(els[1], &user.MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET), nil
	})
	assert.NoError(t, err)
	claims, ok := token.Claims.(*user.MyCustomClaims)
	assert.Equal(t, true, ok)

	mockRequest = `{
		"id": "` + claims.ID + `",
		"name": "updated name",
		"email": "test@test.com"
	}`
	req, err = http.NewRequest("POST", "/user/update", bytes.NewBufferString(mockRequest))
	assert.NoError(t, err)
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
	req, err := http.NewRequest("POST", "/user/login", bytes.NewBufferString(mockRequest))
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusOK, w.Code)

	var response loginResponse
	err = json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)

	req, _ = http.NewRequest("POST", "/user/refresh-token", bytes.NewBufferString(mockRequest))
	req.Header.Set("Authorization", response.Data)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestThumbnail(t *testing.T) {
	mockRequest := `{
		"password": "test",
		"email": "test@test.com"
	}`
	req, err := http.NewRequest("POST", "/user/login", bytes.NewBufferString(mockRequest))
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusOK, w.Code)

	var response loginResponse
	err = json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)

	path := "test-files/thumbnail.png"
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("thumbnail", path)
	assert.NoError(t, err)

	file, err := os.Open(path)
	assert.NoError(t, err)

	_, err = io.Copy(part, file)
	assert.NoError(t, err)
	assert.NoError(t, writer.Close())

	req, err = http.NewRequest("POST", "/user/upload-thumbnail", body)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", response.Data)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusOK, w.Code)

	els := strings.Split(response.Data, " ")
	token, err := jwt.ParseWithClaims(els[1], &user.MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET), nil
	})
	assert.NoError(t, err)
	claims, ok := token.Claims.(*user.MyCustomClaims)
	if !ok {
		panic("Error in claims")
	}

	req, err = http.NewRequest("GET", "/user/thumbnail/"+claims.ID+".png", bytes.NewBufferString(""))
	assert.NoError(t, err)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusOK, w.Code)
}
