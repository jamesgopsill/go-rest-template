package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/rs/zerolog/log"
)

var r = initialiseApp()

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
	req, _ := http.NewRequest("GET", "/user/register", bytes.NewBufferString(mockRequest))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	log.Info().Msg(w.Body.String())
	assert.Equal(t, http.StatusOK, w.Code)
}
