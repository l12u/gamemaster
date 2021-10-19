package main

import (
	"github.com/gin-gonic/gin"
	"github.com/l12u/gamemaster/internal/routes"
)

/*

ENABLE_REDIS_STORAGE

REDIS_HOST
REDIS_PASSWORD
REDIS_DATABASE

*/

func main() {
	println("Hello World!")

	gin.DisableConsoleColor()
	r := gin.Default()

	r.POST("/games", routes.PostGame)
	r.DELETE("/games/:id", routes.DeleteGame)
	r.GET("/games", routes.GetAllGames)
	r.GET("/games/:id", routes.GetGame)
	r.POST("/games/:id/players", routes.PostPlayerToGame)
	r.DELETE("/games/:id/players/:pId", routes.DeletePlayerFromGame)
	r.PUT("/games/:id/state", routes.PutState)

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP"})
	})

	_ = r.Run()
}
