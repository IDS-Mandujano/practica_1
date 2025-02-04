package server1

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/users", getUsers)
	r.POST("/users", createUser)
	r.GET("/longpoll", longPolling)

	r.Run(":8080")
}