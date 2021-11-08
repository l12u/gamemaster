package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Board struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}

type BoardConfig struct {
	Boards []*Board `json:"boards"`
}

func Empty() *BoardConfig {
	return &BoardConfig{Boards: []*Board{}}
}

func (b *BoardConfig) GetBoard(t string) *Board {
	for _, board := range b.Boards {
		if board.Type == t {
			return board
		}
	}
	return nil
}

func FromFile(path string) (*BoardConfig, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var cfg BoardConfig
	err = json.Unmarshal(bytes, &cfg)
	if err != nil {
		return nil, err
	}

	if cfg.Boards == nil {
		cfg.Boards = make([]*Board, 0)
	}

	return &cfg, nil
}
