package server1

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var users []User
var lastID int
var lastUpdateTime time.Time
var mutex sync.Mutex

func getUsers(c *gin.Context) {
	mutex.Lock()
	defer mutex.Unlock()
	c.JSON(http.StatusOK, users)
}

func createUser(c *gin.Context) {
	var newUser User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mutex.Lock()
	lastID++
	newUser.ID = lastID
	users = append(users, newUser)
	lastUpdateTime = time.Now()
	mutex.Unlock()

	c.JSON(http.StatusCreated, newUser)
}

func longPolling(c *gin.Context) {
	prevUpdateTime := c.Query("timestamp")

	for {
		mutex.Lock()
		if lastUpdateTime.String() != prevUpdateTime {
			c.JSON(http.StatusOK, users)
			mutex.Unlock()
			return
		}
		mutex.Unlock()
		time.Sleep(2 * time.Second)
	}
}