package main

import (
	"github.com/gin-gonic/gin"
	"github.com/l12u/gamemaster/internal/handler"
	"github.com/l12u/gamemaster/internal/middleware"
	"k8s.io/klog"
	"time"
)

func main() {
	klog.Infoln("Hello World!")

	handler.SetupProvider()
	handler.SetupBoardConfigAuto()

	gin.DisableConsoleColor()
	r := gin.New()
	r.Use(middleware.Logger(3 * time.Second))
	r.Use(gin.Recovery())

	r.POST("/games", handler.PostGame)
	r.GET("/games", handler.GetAllGames)
	r.GET("/games/:id", handler.GetGame)
	r.DELETE("/games/:id", handler.DeleteGame)
	r.POST("/games/:id/players", handler.PostPlayerToGame)
	r.PUT("/games/:id/players/:pId", handler.PutPlayerToGame)
	r.DELETE("/games/:id/players/:pId", handler.DeletePlayerFromGame)
	r.PUT("/games/:id/state", handler.PutState)

	r.GET("/boards", handler.GetAllBoards)
	r.GET("/boards/:type", handler.GetBoard)

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP"})
	})

	_ = r.Run()
}
