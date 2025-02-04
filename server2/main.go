package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

var replicatedUsers []User
var lastUpdateTime time.Time
var mutex sync.Mutex

func fetchUsersShortPolling() {
	for {
		resp, err := http.Get("http://localhost:8080/users")
		if err == nil {
			body, _ := ioutil.ReadAll(resp.Body)
			var users []User
			if err := json.Unmarshal(body, &users); err == nil {
				mutex.Lock()
				replicatedUsers = users
				lastUpdateTime = time.Now()
				mutex.Unlock()
				fmt.Println("Usuarios replicados:", replicatedUsers)
			}
		}
		time.Sleep(5 * time.Second)
	}
}

func fetchUsersLongPolling() {
	for {
		resp, err := http.Get("http://localhost:8080/longpoll?timestamp=" + lastUpdateTime.String())
		if err == nil {
			body, _ := ioutil.ReadAll(resp.Body)
			var users []User
			if err := json.Unmarshal(body, &users); err == nil {
				mutex.Lock()
				replicatedUsers = users
				lastUpdateTime = time.Now()
				mutex.Unlock()
				fmt.Println("Usuarios replicados (Long Polling):", replicatedUsers)
			}
		}
	}
}

func getReplicatedUsers(c *gin.Context) {
	mutex.Lock()
	defer mutex.Unlock()
	c.JSON(http.StatusOK, replicatedUsers)
}

func main() {
	go fetchUsersShortPolling()
	go fetchUsersLongPolling()

	r := gin.Default()
	r.GET("/replicated-users", getReplicatedUsers)

	r.Run(":8081")
}