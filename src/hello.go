package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"
)

const (
	version = "0.2.0"
)

type DateOfBirth struct {
	Value CustomTime `json:"dateOfBirth" binding:"required"`
}

// not 100% sure about this implementation
// for `high load` setup we need shared db connection, rather than global Rdb var

var Rdb = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
})

func main() {
	if len(os.Args) < 3 {
		log.Fatal("please specify port and build arguments")
	}
	port := os.Args[1]
	router := gin.Default()
	router.GET("/hc", healthCheck)
	router.GET("/version", getVersion)
	router.PUT("/hello/:username", putUser)
	router.GET("/hello/:username", getUser)

	pong, err := Rdb.Ping().Result()
	log.Println(pong, err)

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
	if err != nil {
		body["error"] = err.Error()
		c.JSON(http.StatusInternalServerError, body)
	} else {
		body["hostname"] = hostname
		c.JSON(http.StatusOK, body)
	}
}

func putUser(c *gin.Context) {
	username := c.Param("username")
	// validate username
	var alpha = regexp.MustCompile(`^[[:alpha:]]+$`).MatchString
	if !alpha(username) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username should contain only letters"})
	} else {
		var birthday DateOfBirth
		// validate json payload
		c.Header("Content-Type", "application/json") // BindJSON set text/plain if error
		err := c.BindJSON(&birthday)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "cannot find date of birth",
				"details": err.Error(),
			})
		} else {
			// validate future date
			now := time.Now()
			if now.Before(time.Time(birthday.Value)) {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "date of birth must be a date before the today date",
				})
			} else {
				textual := time.Time(birthday.Value).Format(CustomFormat)
				// post data to db
				err := Rdb.Set(username, textual, 0).Err()
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error":   "storage not ready",
						"details": err,
					})
				} else {
					c.Status(http.StatusNoContent)
					log.Printf("%s, %s", username, textual)
				}
			}
		}
	}
}

func getUser(c *gin.Context) {
	username := c.Param("username")
	// read data from db
	v, err := Rdb.Get(username).Result()
	if err == redis.Nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "user not found",
		})
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "storage not ready",
			"details": err,
		})
	} else {
		birthday, err := time.Parse(CustomFormat, v)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "cannot parse value from storage",
				"details": err,
			})
		} else {
			//now := time.Now()
			n := 1
			var msg string
			if n == 0 {
				msg = fmt.Sprintf("Hello, %s! Happy birthday!", username)
			} else {
				msg = fmt.Sprintf("Hello, %s! Your birthday is in %d day(s), %v", username, n, birthday)
			}
			c.JSON(http.StatusOK, gin.H{
				"message": msg,
			})
		}
	}
}
