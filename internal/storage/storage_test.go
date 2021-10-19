package storage

import (
	"fmt"
	"github.com/l12u/gamemaster/internal/model"
	"testing"
)

// TODO test all functions with both local and redis (mocked) storage providers

func TestLocalProvider_PutGame(t *testing.T) {
	p := NewLocalProvider()
	testPutGame(p, t)
}

func testPutGame(p Provider, t *testing.T) {
	games0, err := p.GetAllGames()
	if err != nil {
		t.Errorf("failed to execute GetAllGames(): %v", err)
	}
	if len(games0) != 0 {
		t.Errorf("expected to be empty, but contains %d games", len(games0))
	}

	// test when putting in a specific game

	g0 := model.Game{Id: "someId"}
	err = p.PutGame(&g0)
	if err != nil {
		t.Errorf("failed to execute PutGame(%v): %v", g0, err)
	}
	games1, err := p.GetAllGames()
	if err != nil {
		t.Errorf("failed to execute GetAllGames(): %v", err)
	}
	if len(games1) != 1 {
		t.Errorf("expected to contain 1 game, but contains %d games", len(games1))
	}

	g1, err := p.GetGame("someId")
	if err != nil {
		t.Errorf("failed to execute GetGame(%s): %v", "someId", err)
	}
	if g1 == nil {
		t.Errorf("got nil when getting game with id %s", "someId")
	}
	if g1.Id != "someId" {
		t.Errorf("expected to get game with id = %s, but got %s", "someId", g1.Id)
	}

	// test when putting the same twice

	g2 := model.Game{Id: "someId", State: "someOtherState"}
	err = p.PutGame(&g2)
	if err != nil {
		t.Errorf("failed to execute PutGame(%v): %v", g2, err)
	}
	games2, err := p.GetAllGames()
	if err != nil {
		t.Errorf("failed to execute GetAllGames(): %v", err)
	}
	if len(games2) != 1 {
		t.Errorf("expected to still contain 1 game, but contains %d games", len(games2))
	}

	g3, err := p.GetGame("someId")
	if err != nil {
		t.Errorf("failed to execute GetGame(%s): %v", "someId", err)
	}
	if g3 == nil {
		t.Errorf("got nil when getting game with id %s", "someId")
	}
	if g3.Id != "someId" {
		t.Errorf("expected to get game with id = %s, but got %s", "someId", g3.Id)
	}
	if g3.State != "someOtherState" {
		t.Errorf("expected to get game with state = %s, but got %s", "someOtherState", g3.State)
	}
	if ok, _ := p.HasGame("someId"); !ok {
		t.Errorf("expected HasGame() to return true")
	}

	// test when inserting 10 items, we get 10 more games

	games3, err := p.GetAllGames()
	if err != nil {
		t.Errorf("failed to execute GetAllGames(): %v", err)
	}
	prevSize := len(games3)

	for i := 0; i < 10; i++ {
		g4 := model.Game{Id: fmt.Sprintf("someOtherId%d", i)}
		err = p.PutGame(&g4)
		if err != nil {
			t.Errorf("failed to execute PutGame(%v): %v", g4, err)
		}
	}

	games4, err := p.GetAllGames()
	if err != nil {
		t.Errorf("failed to execute GetAllGames(): %v", err)
	}
	afterSize := len(games4)

	if afterSize-prevSize != 10 {
		t.Errorf("expected to get %d more entries, but we got %d", 10, afterSize-prevSize)
	}
}
