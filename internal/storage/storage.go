package storage

import "github.com/l12u/gamemaster/internal/model"

type Provider interface {
	PutGame(g *model.Game) error
	GetGame(id string) (*model.Game, error)
	GetAllGames() (model.GameMap, error)
	DeleteGame(id string) error
	ClearGames() error
	HasGame(id string) (bool, error)

	PutBoard(b *model.Board) error
	GetBoard(id string) (*model.Board, error)
	GetBoards(t string) (model.BoardMap, error)
	GetAllBoards() (model.BoardMap, error)
	DeleteBoard(id string) error
	HasBoard(id string) (bool, error)
}
