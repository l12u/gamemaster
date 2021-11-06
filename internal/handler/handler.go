package handler

import (
	"fmt"
	"github.com/l12u/gamemaster/internal/config"
	"github.com/l12u/gamemaster/internal/errcode"
	"github.com/l12u/gamemaster/internal/storage"
	"github.com/l12u/gamemaster/pkg/env"
	"github.com/l12u/gamemaster/pkg/valid"
	"k8s.io/klog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/l12u/gamemaster/internal/model"
)

var provider storage.Provider
var cfg *config.BoardConfig

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

	c, err := config.FromFile(env.StringOrDefault("BOARD_CONFIG", "/etc/gamemaster/boards.json"))
	if err != nil {
		// give a warning
		klog.Warningln("no board config found, that is not good :(")
		c = config.Empty()
	}
	cfg = c
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
		g.Players = make([]*model.Player, 0)
	} else {
		for _, player := range g.Players {
			if !model.IsSupportedRole(player.Role) {
				errcode.D(c, errcode.InvalidRole,
					fmt.Sprintf("invalid role given for player %s", player.Name),
				)
			}
		}
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

	c.JSON(http.StatusOK, games.AsSlice())
}

// GetGame gets a specific game
func GetGame(c *gin.Context) {
	id := c.Param("id")

	if ok, _ := provider.HasGame(id); !ok {
		errcode.D(c, errcode.GameDoesNotExist, "a game with this id does not exist")
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
		_ = c.Error(err)
		return
	}

	if ok, _ := provider.HasGame(id); !ok {
		errcode.D(c, errcode.GameDoesNotExist, "a game with this id does not exist")
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
		errcode.D(c, errcode.PlayerAlreadyExist, "already a player with this id in the game")
		return
	}

	g.Players = append(g.Players, &player)
	err = provider.PutGame(g)
	if err != nil {
		_ = c.Error(err)
		return
	}
}

// PutPlayerToGame updates a player from a specific game
func PutPlayerToGame(c *gin.Context) {
	id := c.Param("id")
	pId := c.Param("pId")

	var data *model.Player
	err := c.BindJSON(&data)
	if err != nil {
		_ = c.Error(err)
		return
	}

	if ok, _ := provider.HasGame(id); !ok {
		errcode.D(c, errcode.GameDoesNotExist, "a game with this id does not exist")
		return
	}

	g, err := provider.GetGame(id)
	if err != nil {
		_ = c.Error(err)
		return
	}

	p := g.GetPlayer(pId)
	if p == nil {
		errcode.D(c, errcode.PlayerDoesNotExist, "a player with that id does not exist in that game")
		return
	}

	newRole := data.Role
	newName := data.Name

	if newRole != "" {
		if !model.IsSupportedRole(newRole) {
			errcode.D(c, errcode.InvalidRole, "invalid role")
			return
		}
		p.Role = newRole
	}
	if newName != "" {
		p.Name = newName
	}
}

// DeletePlayerFromGame deletes a player from a specific game
func DeletePlayerFromGame(c *gin.Context) {
	id := c.Param("id")
	pId := c.Param("pId")

	if ok, _ := provider.HasGame(id); !ok {
		errcode.D(c, errcode.GameDoesNotExist, "a game with this id does not exist")
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
		errcode.D(c, errcode.PlayerDoesNotExist, "no player with this id in the game")
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
		errcode.D(c, errcode.GameDoesNotExist, "a game with this id does not exist")
		return
	}

	if !model.IsSupportedState(req.NewState) {
		errcode.D(c, errcode.InvalidState, "not supported game state given")
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
		errcode.D(c, errcode.GameDoesNotExist, "a game with this id does not exist")
		return
	}

	err := provider.DeleteGame(id)
	if err != nil {
		_ = c.Error(err)
		return
	}
}

// GetBoard returns all boards of a specific type
func GetBoard(c *gin.Context) {
	t := c.Param("type")

	if !valid.ValidateSlug(t) {
		errcode.D(c, errcode.InvalidType, "wrong format for type")
		return
	}

	var board *config.Board
	for _, b := range cfg.Boards {
		if b.Type == t {
			board = b
		}
	}

	if board == nil {
		errcode.D(c, errcode.BoardDoesNotExist, "a board with this type does not exist")
		return
	}

	c.JSON(http.StatusOK, board)
}

// GetAllBoards returns all registered boards
func GetAllBoards(c *gin.Context) {
	c.JSON(http.StatusOK, cfg.Boards)
}
