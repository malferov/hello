package main

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func performRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestHealthCheck(t *testing.T) {
	router := gin.New()
	router.GET("/hc", healthCheck)

	w := performRequest(router, "GET", "/hc")

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "OK", w.Body.String())
}

func TestGetVersion(t *testing.T) {
	router := gin.New()
	router.GET("/version", getVersion)

	w := performRequest(router, "GET", "/version")

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "welcome")
}

func TestPutUser(t *testing.T) {
	router := gin.New()
	router.PUT("/hello/:username", putUser)

	w := performRequest(router, "PUT", "/hello/john")

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Equal(t, "", w.Body.String())
}

//assert.JSONEq(t, `{}`, w.Body.JSON())
