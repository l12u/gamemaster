package storage

import (
	"github.com/l12u/gamemaster/internal/model"
)

type Provider interface {
	PutGame(g *model.Game) error
	GetGame(id string) (*model.Game, error)
	GetAllGames() (model.GameMap, error)
	DeleteGame(id string) error
	ClearGames() error
	HasGame(id string) (bool, error)
}
