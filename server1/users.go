package server1

import (
	"net/http"
	"sync"
	"time"
	"github.com/gin-gonic/gin"
	"fmt"
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

func deleteUser(c *gin.Context) {
	id := c.Param("id")

	mutex.Lock()
	defer mutex.Unlock()

	var userID int
	_, err := fmt.Sscanf(id, "%d", &userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	for i, user := range users {
		if user.ID == userID {
			users = append(users[:i], users[i+1:]...)
			lastUpdateTime = time.Now()
			c.JSON(http.StatusOK, gin.H{"message": "Usuario eliminado"})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
}

func updateUser(c *gin.Context) {
	id := c.Param("id")
	var updatedUser User

	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	var userID int
	_, err := fmt.Sscanf(id, "%d", &userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	for i, user := range users {
		if user.ID == userID {
			users[i].Name = updatedUser.Name
			users[i].Email = updatedUser.Email
			lastUpdateTime = time.Now()
			c.JSON(http.StatusOK, users[i])
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
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