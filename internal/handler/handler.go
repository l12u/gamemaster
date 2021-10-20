package handler

import (
	"github.com/l12u/gamemaster/internal/storage"
	"github.com/l12u/gamemaster/pkg/env"
	"github.com/l12u/gamemaster/pkg/valid"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/l12u/gamemaster/internal/model"
)

var provider storage.Provider

func SetupProvider() {
	enableRedis := env.BoolOrDefault("ENABLE_REDIS_STORAGE", false)
	if enableRedis {
		provider = storage.NewRedisProvider(
			env.StringOrDefault("REDIS_ADDRESS", "localhost:6379"),
			env.StringOrDefault("REDIS_PASSWORD", ""),
			env.IntOrDefault("REDIS_DATABASE", 0),
		)
		err := provider.(*storage.RedisProvider).Connect()
		if err != nil {
			panic(err)
		}
	} else {
		provider = storage.NewLocalProvider()
	}
}

// PostGame inserts a game into the registry
func PostGame(c *gin.Context) {
	var g model.Game

	err := c.BindJSON(&g)
	if err != nil {
		_ = c.Error(err)
		return
	}

	ct := time.Now().UnixMilli()
	g.CreatedAt = ct
	g.UpdatedAt = ct

	if g.Players == nil {
		g.Players = make([]model.Player, 0)
	}
	if g.Roles == nil {
		g.Roles = make(map[string]string)
	}
	g.State = model.StateLobby
	g.Id = uuid.NewString()

	err = provider.PutGame(&g)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": g.Id})
}

// GetAllGames simply returns all games
func GetAllGames(c *gin.Context) {
	games, err := provider.GetAllGames()
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, games)
}

// GetGame gets a specific game
func GetGame(c *gin.Context) {
	id := c.Param("id")

	if ok, _ := provider.HasGame(id); !ok {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "a game with this id does not exist"})
		return
	}

	g, err := provider.GetGame(id)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, g)
}

// PostPlayerToGame adds a player to a specific game
func PostPlayerToGame(c *gin.Context) {
	id := c.Param("id")
	var player model.Player

	err := c.BindJSON(&player)
	if err != nil {
		return
	}

	if ok, _ := provider.HasGame(id); !ok {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "a game with this id does not exist"})
		return
	}

	g, err := provider.GetGame(id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	ok := false
	for _, p := range g.Players {
		if p.Id == player.Id {
			ok = true
			break
		}
	}

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "already a player with this id in the game"})
		return
	}

	g.Players = append(g.Players, player)
	err = provider.PutGame(g)
	if err != nil {
		_ = c.Error(err)
		return
	}
}

// DeletePlayerFromGame deletes a player from a specific game
func DeletePlayerFromGame(c *gin.Context) {
	id := c.Param("id")
	pId := c.Param("pId")

	if ok, _ := provider.HasGame(id); !ok {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "a game with this id does not exist"})
		return
	}

	g, err := provider.GetGame(id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	var pi int
	ok := false
	for i, p := range g.Players {
		if p.Id == pId {
			pi = i
			ok = true
			break
		}
	}

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "no player with this id in the game"})
		return
	}

	g.Players = append(g.Players[:pi], g.Players[pi+1:]...)
	err = provider.PutGame(g)
	if err != nil {
		_ = c.Error(err)
		return
	}
}

// PutState changes the state of a game
func PutState(c *gin.Context) {
	id := c.Param("id")

	var req model.UpdateStateRequest
	if c.BindJSON(&req) != nil {
		return
	}

	if ok, _ := provider.HasGame(id); !ok {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "a game with this id does not exist"})
		return
	}

	if !model.IsSupportedState(req.NewState) {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "not supported game state given"})
		return
	}

	g, err := provider.GetGame(id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	g.State = req.NewState
	err = provider.PutGame(g)
	if err != nil {
		_ = c.Error(err)
		return
	}
}

// DeleteGame deletes a specific game
func DeleteGame(c *gin.Context) {
	id := c.Param("id")

	if ok, _ := provider.HasGame(id); !ok {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "a game with this id does not exist"})
		return
	}

	err := provider.DeleteGame(id)
	if err != nil {
		_ = c.Error(err)
		return
	}
}

// PostBoard registers a new board with its
// respecting URL.
// Here we can have multiple boards with the
// same type.
func PostBoard(c *gin.Context) {
	var b model.Board

	err := c.BindJSON(&b)
	if err != nil {
		_ = c.Error(err)
		return
	}

	if !valid.ValidateURL(b.URL) {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "the url is not valid"})
		return
	}
	if !valid.ValidateSlug(b.Type) {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "wrong format for type"})
		return
	}

	ct := time.Now().UnixMilli()
	b.RegisteredAt = ct
	b.Id = uuid.NewString()

	err = provider.PutBoard(&b)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": b.Id})
}

// GetBoards returns all boards of a specific type
func GetBoards(c *gin.Context) {
	t := c.Param("type")

	if !valid.ValidateSlug(t) {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "wrong format for type"})
		return
	}

	boards, err := provider.GetBoards(t)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, boards)
}

// GetAllBoards returns all registered boards
func GetAllBoards(c *gin.Context) {
	boards, err := provider.GetAllBoards()
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, boards)
}

// DeleteBoard deletes an already registered board.
func DeleteBoard(c *gin.Context) {
	id := c.Param("id")

	if ok, _ := provider.HasBoard(id); !ok {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "a board with this id does not exist"})
		return
	}

	err := provider.DeleteBoard(id)
	if err != nil {
		_ = c.Error(err)
		return
	}
}
