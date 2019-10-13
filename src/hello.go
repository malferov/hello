package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	version = "0.2.0"
)

type user struct {
	username string    `json:"userName"`
	birthday time.Time `json:"dateOfBirth"`
}

func main() {
	if len(os.Args) < 3 {
		log.Fatal("please specify port and build arguments")
	}
	port := os.Args[1]
	router := gin.Default()
	router.GET("/hc", healthCheck)
	router.GET("/version", getVersion)
	router.PUT("/hello/:username", putUser)
	router.Run(":" + port)
}

func healthCheck(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}

func getVersion(c *gin.Context) {
	build := "n/a"
	if len(os.Args) == 3 {
		build = os.Args[2]
	}
	body := gin.H{
		"data":    "welcome",
		"version": version,
		"build":   build,
		"lang":    "golang",
	}
	hostname, err := os.Hostname()
	if err == nil {
		body["hostname"] = hostname
		c.JSON(http.StatusOK, body)
	} else {
		body["error"] = err.Error()
		c.JSON(http.StatusInternalServerError, body)
	}
}

func putUser(c *gin.Context) {
	c.Status(http.StatusNoContent)
}
