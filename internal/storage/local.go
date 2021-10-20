package storage

import (
	"github.com/l12u/gamemaster/internal/model"
)

type LocalProvider struct {
	games  model.GameMap
	boards model.BoardMap
}

func NewLocalProvider() *LocalProvider {
	return &LocalProvider{
		games:  make(model.GameMap),
		boards: make(model.BoardMap),
	}
}

func (l *LocalProvider) PutGame(g *model.Game) error {
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

func (l *LocalProvider) PutBoard(b *model.Board) error {
	l.boards[b.Id] = b
	return nil
}

func (l *LocalProvider) GetBoard(id string) (*model.Board, error) {
	return l.boards[id], nil
}

func (l *LocalProvider) GetBoards(t string) (model.BoardMap, error) {
	bm := make(model.BoardMap)
	for _, board := range l.boards {
		if board.Type == t {
			bm[board.Id] = board
		}
	}
	return bm, nil
}

func (l *LocalProvider) GetAllBoards() (model.BoardMap, error) {
	return l.boards, nil
}

func (l *LocalProvider) DeleteBoard(id string) error {
	delete(l.boards, id)
	return nil
}

func (l *LocalProvider) HasBoard(id string) (bool, error) {
	_, ok := l.boards[id]
	return ok, nil
}
