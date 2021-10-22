package storage

import (
	"fmt"
	"github.com/l12u/gamemaster/internal/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TODO test all functions with both local and redis (mocked) storage providers

func TestLocalProvider_PutGame(t *testing.T) {
	p := NewLocalProvider()
	testPutGame(p, t)
}

func testPutGame(p Provider, t *testing.T) {
	a := assert.New(t)

	games0, err := p.GetAllGames()
	a.NoError(err)
	a.True(len(games0) == 0)

	// test when putting in a specific game

	g0 := model.Game{Id: "someId"}
	err = p.PutGame(&g0)
	a.NoError(err)

	games1, err := p.GetAllGames()
	a.NoError(err)
	a.True(len(games1) == 1)

	g1, err := p.GetGame("someId")
	a.NoError(err)
	if a.NotNil(g1) {
		a.True(g1.Id == "someId")
	}

	// test when putting the same twice

	g2 := model.Game{Id: "someId", State: "someOtherState"}
	err = p.PutGame(&g2)
	a.NoError(err)

	games2, err := p.GetAllGames()
	a.NoError(err)
	a.True(len(games2) == 1)

	g3, err := p.GetGame("someId")
	a.NoError(err)
	if a.NotNil(g3) {
		a.True(g3.Id == "someId")
	}
	a.True(g3.State == "someOtherState")
	ok, err := p.HasGame("someId")
	a.NoError(err)
	a.True(ok)

	// test when inserting 10 items, we get 10 more games

	games3, err := p.GetAllGames()
	a.NoError(err)
	prevSize := len(games3)

	for i := 0; i < 10; i++ {
		g4 := model.Game{Id: fmt.Sprintf("someOtherId%d", i)}
		err = p.PutGame(&g4)
		a.NoError(err)
	}

	games4, err := p.GetAllGames()
	a.NoError(err)
	afterSize := len(games4)

	a.True(afterSize-prevSize == 10)
}
