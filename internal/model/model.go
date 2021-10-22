package model

import "strings"

const (
	StateLobby   = "lobby"
	StateRunning = "running"

	RoleHost   = "host"
	RolePlayer = "player"
)

var supportedStates = []string{StateLobby, StateRunning}
var supportedRoles = []string{RoleHost, RolePlayer}

var EmptyGameMap = make(GameMap)

type Player struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

type Game struct {
	Id        string      `json:"id"`
	Type      string      `json:"type"`
	Players   []*Player   `json:"players"`
	State     string      `json:"state"`
	CreatedAt int64       `json:"createdAt"`
	UpdatedAt int64       `json:"updatedAt"`
	GameData  interface{} `json:"gameData"`
}

type GameMap map[string]*Game

type UpdateStateRequest struct {
	NewState string `json:"state"`
}

func (g GameMap) AsSlice() []*Game {
	sl := make([]*Game, len(g))
	for _, game := range g {
		sl = append(sl, game)
	}
	return sl
}

func (g *Game) GetPlayer(pid string) *Player {
	for _, player := range g.Players {
		if player.Id == pid {
			return player
		}
	}

	return nil
}

func IsSupportedState(s string) bool {
	var l = strings.ToLower(s)

	for _, ss := range supportedStates {
		if l == ss {
			return true
		}
	}
	return false
}

func IsSupportedRole(s string) bool {
	var l = strings.ToLower(s)

	for _, ss := range supportedRoles {
		if l == ss {
			return true
		}
	}
	return false
}
