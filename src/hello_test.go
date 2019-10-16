package main

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func performRequest(r http.Handler, method, path string, body io.Reader) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, body)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestHealthCheck(t *testing.T) {
	router := gin.New()
	router.GET("/hc", healthCheck)
	// normal request
	w := performRequest(router, "GET", "/hc", nil)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "OK", w.Body.String())
}

func TestGetVersion(t *testing.T) {
	router := gin.New()
	router.GET("/version", getVersion)
	// normal request
	w := performRequest(router, "GET", "/version", nil)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "welcome")
}

func TestPutUser(t *testing.T) {
	router := gin.New()
	router.PUT("/hello/:username", putUser)

	testcase := "normal request"
	payload := []byte(`{"dateOfBirth":"2006-01-02"}`)
	w := performRequest(router, "PUT", "/hello/john", bytes.NewBuffer(payload))
	assert.Equal(t, http.StatusNoContent, w.Code, testcase)
	assert.Empty(t, w.Body.String(), testcase)

	testcase = "malformed username"
	w = performRequest(router, "PUT", "/hello/john7", nil)
	assert.Equal(t, http.StatusBadRequest, w.Code, testcase)
	assert.Contains(t, w.Body.String(), "error", testcase)

	testcase = "future date"
	payload = []byte(`{"dateOfBirth":"2020-01-02"}`)
	w = performRequest(router, "PUT", "/hello/john", bytes.NewBuffer(payload))
	assert.Equal(t, http.StatusBadRequest, w.Code, testcase)
	assert.Contains(t, w.Body.String(), "error", testcase)

	testcase = "malformed date key"
	payload = []byte(`{"date_of_birth":"2006-01-02"}`)
	w = performRequest(router, "PUT", "/hello/john", bytes.NewBuffer(payload))
	assert.Equal(t, http.StatusBadRequest, w.Code, testcase)
	assert.Contains(t, w.Body.String(), "error", testcase)

	testcase = "malformed date value"
	payload = []byte(`{"dateOfBirth":"1234-56-78"}`)
	w = performRequest(router, "PUT", "/hello/john", bytes.NewBuffer(payload))
	assert.Equal(t, http.StatusBadRequest, w.Code, testcase)
	assert.Contains(t, w.Body.String(), "error", testcase)
}

func TestGetUser(t *testing.T) {
	router := gin.New()
	router.GET("/hello/:username", getUser)

	testcase := "normal request"
	w := performRequest(router, "GET", "/hello/john", nil)
	assert.Equal(t, http.StatusOK, w.Code, testcase)
	assert.Contains(t, w.Body.String(), "message", testcase)

	testcase = "malformed username or user not found"
	w = performRequest(router, "GET", "/hello/john7", nil)
	assert.Equal(t, http.StatusNotFound, w.Code, testcase)
	assert.JSONEq(t, `{"message": "user not found"}`, w.Body.String(), testcase)
}
