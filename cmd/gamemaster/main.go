package main

import (
	"github.com/gin-gonic/gin"
	"github.com/l12u/gamemaster/internal/router"
)

func main() {
	println("Hello World!")

	router.SetupProvider()

	gin.DisableConsoleColor()
	r := gin.Default()

	r.POST("/games", router.PostGame)
	r.DELETE("/games/:id", router.DeleteGame)
	r.GET("/games", router.GetAllGames)
	r.GET("/games/:id", router.GetGame)
	r.POST("/games/:id/players", router.PostPlayerToGame)
	r.DELETE("/games/:id/players/:pId", router.DeletePlayerFromGame)
	r.PUT("/games/:id/state", router.PutState)

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP"})
	})

	_ = r.Run()
}
