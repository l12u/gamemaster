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

var EmptyGameMap = make(GameMap, 0)
var EmptyBoardMap = make(BoardMap, 0)

type Player struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Game struct {
	Id        string            `json:"id"`
	Type      string            `json:"type"`
	Players   []Player          `json:"players"`
	Roles     map[string]string `json:"roles"`
	State     string            `json:"state"`
	CreatedAt int64             `json:"createdAt"`
	UpdatedAt int64             `json:"updatedAt"`
	GameData  interface{}       `json:"gameData"`
}

type GameMap map[string]*Game

type UpdateStateRequest struct {
	NewState string `json:"state"`
}

type Board struct {
	Id           string `json:"id"`
	Type         string `json:"type"`
	URL          string `json:"url"`
	RegisteredAt int64  `json:"registeredAt"`
}

type BoardMap map[string]*Board

func (g GameMap) AsSlice() []*Game {
	sl := make([]*Game, len(g))
	for _, game := range g {
		sl = append(sl, game)
	}
	return sl
}

func (b BoardMap) AsSlice() []*Board {
	sl := make([]*Board, len(b))
	for _, board := range b {
		sl = append(sl, board)
	}
	return sl
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
