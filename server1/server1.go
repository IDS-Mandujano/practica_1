package server1

import "github.com/gin-gonic/gin"

func StartServer1() {
    r := gin.Default()
    r.GET("/users", getUsers)
    r.POST("/users", createUser)
    r.GET("/longpoll", longPolling)

    r.GET("/", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "Servidor 1 activo"})
    })

    r.Run(":9090")
}