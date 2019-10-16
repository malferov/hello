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
	Addr:     os.Getenv("REDIS_ENDPOINT"),
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
		var msg = "cannot find date of birth"
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   msg,
				"details": err.Error(),
			})
		} else if time.Time(birthday.Value).IsZero() {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   msg,
				"details": "malformed key",
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
			n := daysNextBirthday(birthday)
			var msg string
			if n == 0 {
				msg = fmt.Sprintf("Hello, %s! Happy birthday!", username)
			} else {
				msg = fmt.Sprintf("Hello, %s! Your birthday is in %d day(s)", username, n)
			}
			c.JSON(http.StatusOK, gin.H{
				"message": msg,
			})
		}
	}
}

func daysNextBirthday(b time.Time) int {
	now := time.Now().Truncate(24 * time.Hour) // truncate the time to the day
	next_birthday := time.Date(
		now.Year(),
		b.Month(),
		b.Day(),
		0, 0, 0, 0, // truncate the time to the day
		time.UTC, // strip location
	)
	duration := next_birthday.Sub(now)
	if duration < 0 { // birthday passed in this year
		next_birthday = time.Date(now.Year()+1, b.Month(), b.Day(), 0, 0, 0, 0, time.UTC)
		duration = next_birthday.Sub(now)
	}
	return int(duration.Hours() / 24)
}
