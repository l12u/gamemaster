package main

import (
	"github.com/gin-gonic/gin"
	"github.com/l12u/gamemaster/internal/handler"
)

func main() {
	println("Hello World!")

	handler.SetupProvider()

	gin.DisableConsoleColor()
	r := gin.Default()

	r.POST("/games", handler.PostGame)
	r.DELETE("/games/:id", handler.DeleteGame)
	r.GET("/games", handler.GetAllGames)
	r.GET("/games/:id", handler.GetGame)
	r.POST("/games/:id/players", handler.PostPlayerToGame)
	r.DELETE("/games/:id/players/:pId", handler.DeletePlayerFromGame)
	r.PUT("/games/:id/state", handler.PutState)

	r.POST("/boards", handler.PostBoard)
	r.GET("/boards/:type", handler.GetBoards)
	r.GET("/boards", handler.GetAllBoards)
	r.DELETE("/boards/:id", handler.DeleteBoard)

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP"})
	})

	_ = r.Run()
}
