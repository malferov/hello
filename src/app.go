package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
)

const (
	version = "1.8"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatal("please specify port and build arguments")
	}
	port := os.Args[1]
	r := gin.Default()
	body := gin.H{
		"data":    "welcome",
		"version": version,
		"build":   os.Args[2],
		"lang":    "golang",
	}
	r.GET("/", func(c *gin.Context) {
		hostname, err := os.Hostname()
		if err == nil {
			body["hostname"] = hostname
		} else {
			body["error"] = err.Error()
		}
		c.JSON(http.StatusOK, body)
	})
	r.Run(":" + port)
}
