package storage

import (
	"errors"
	"github.com/l12u/gamemaster/internal/model"
)

type LocalProvider struct {
	games model.GameMap
}

func NewLocalProvider() *LocalProvider {
	return &LocalProvider{
		games: make(model.GameMap),
	}
}

func (l *LocalProvider) PutGame(g *model.Game) error {
	if g == nil {
		return errors.New("game can not be nil")
	}

	l.games[g.Id] = g
	return nil
}

func (l *LocalProvider) GetGame(id string) (*model.Game, error) {
	return l.games[id], nil
}

func (l *LocalProvider) GetAllGames() (model.GameMap, error) {
	return l.games, nil
}

func (l *LocalProvider) DeleteGame(id string) error {
	delete(l.games, id)
	return nil
}

func (l *LocalProvider) ClearGames() error {
	l.games = make(model.GameMap)
	return nil
}

func (l *LocalProvider) HasGame(id string) (bool, error) {
	_, ok := l.games[id]
	return ok, nil
}
