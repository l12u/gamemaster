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

	return &cfg, nil
}
