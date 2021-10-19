package model

import "strings"

const (
	StateLobby   = "lobby"
	StateRunning = "running"

	RoleHost   = "host"
	RolePlayer = "player"
)

var supportedStates = []string{StateLobby, StateRunning}

// var supportedRoles = []string{RoleHost, RolePlayer}

type Player struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Game struct {
	Id        string            `json:"id"`
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

func IsSupportedState(s string) bool {
	var l = strings.ToLower(s)

	for _, ss := range supportedStates {
		if l == ss {
			return true
		}
	}
	return false
}
