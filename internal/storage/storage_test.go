package storage

import (
	"github.com/alicebob/miniredis/v2"
	"github.com/l12u/gamemaster/internal/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

// ==================================
// LocalProvider
// ==================================

func TestLocalProvider_PutGame(t *testing.T) {
	p := NewLocalProvider()

	testPutGameMultipleDifferentId(p, t)
	testPutGameMultipleSameId(p, t)
	testPutGameNull(p, t)
}

func TestLocalProvider_DeleteGame(t *testing.T) {
	p := NewLocalProvider()

	testDeleteGameNotPresent(p, t)
	testDeleteGameSameId(p, t)
	testDeleteGameMultipleIds(p, t)
}

func TestLocalProvider_ClearGames(t *testing.T) {
	p := NewLocalProvider()

	a := assert.New(t)
	err := p.ClearGames()
	a.NoError(err)
	a.True(getProviderSize(p) == 0)

	err = p.ClearGames()
	a.NoError(err)
	a.True(getProviderSize(p) == 0)
}

// ==================================
// RedisProvider
// ==================================

func TestRedisProvider_PutGame(t *testing.T) {
	m := newMockedRedisProvider()
	defer m.redis.Close()

	p := m.Provider

	testPutGameMultipleDifferentId(p, t)
	testPutGameMultipleSameId(p, t)
	testPutGameNull(p, t)
}

func TestRedisProvider_DeleteGame(t *testing.T) {
	m := newMockedRedisProvider()
	defer m.redis.Close()

	p := m.Provider

	testDeleteGameNotPresent(p, t)
	testDeleteGameSameId(p, t)
	testDeleteGameMultipleIds(p, t)
}

func TestRedisProvider_ClearGames(t *testing.T) {
	m := newMockedRedisProvider()
	defer m.redis.Close()

	p := m.Provider

	a := assert.New(t)
	err := p.ClearGames()
	a.NoError(err)
	a.True(getProviderSize(p) == 0)

	err = p.ClearGames()
	a.NoError(err)
	a.True(getProviderSize(p) == 0)
}

// ==================================
// PutGame
// ==================================

func testPutGameMultipleDifferentId(p Provider, t *testing.T) {
	a := assert.New(t)
	err := p.ClearGames()
	a.NoError(err)
	a.True(getProviderSize(p) == 0)

	games := []model.Game{
		{
			Id: "SomeId",
		},
		{
			Id: "123456",
		},
		{
			Id: "some_id",
		},
		{
			Id: "some id",
		},
	}

	for _, game := range games {
		err = p.PutGame(&game)
		a.NoError(err)
		a.True(p.HasGame(game.Id))
	}
	a.True(getProviderSize(p) == len(games))
}

func testPutGameMultipleSameId(p Provider, t *testing.T) {
	a := assert.New(t)
	err := p.ClearGames()
	a.NoError(err)
	a.True(getProviderSize(p) == 0)

	games := []model.Game{
		{
			Id: "SomeId",
		},
		{
			Id: "SomeId",
		},
		{
			Id: "SomeId",
		},
	}

	for _, game := range games {
		err = p.PutGame(&game)
		a.NoError(err)
		a.True(p.HasGame(game.Id))
	}
	a.True(getProviderSize(p) == 1)
}

func testPutGameNull(p Provider, t *testing.T) {
	a := assert.New(t)
	err := p.ClearGames()
	a.NoError(err)
	a.True(getProviderSize(p) == 0)

	err = p.PutGame(nil)
	a.Error(err)
	a.True(getProviderSize(p) == 0)
}

// ==================================
// DeleteGame
// ==================================

func testDeleteGameNotPresent(p Provider, t *testing.T) {
	a := assert.New(t)
	err := p.ClearGames()
	a.NoError(err)
	a.True(getProviderSize(p) == 0)

	prepareProvider(p, []string{"foo"})

	// DeleteGame just makes sure, that the game is deleted.
	// no error should be returned, only if something went
	// completely wrong.
	err = p.DeleteGame("bar")
	a.NoError(err)
	a.True(p.HasGame("foo"))
}

func testDeleteGameSameId(p Provider, t *testing.T) {
	a := assert.New(t)
	err := p.ClearGames()
	a.NoError(err)
	a.True(getProviderSize(p) == 0)

	prepareProvider(p, []string{"foo", "bar"})

	err = p.DeleteGame("bar")
	a.NoError(err)
	a.False(p.HasGame("bar"))

	// Even if it is executed multiple times, it should be
	// idempotent.
	err = p.DeleteGame("bar")
	a.NoError(err)
	a.False(p.HasGame("bar"))
}

func testDeleteGameMultipleIds(p Provider, t *testing.T) {
	a := assert.New(t)
	err := p.ClearGames()
	a.NoError(err)
	a.True(getProviderSize(p) == 0)

	gameKeys := []string{"foo", "bar", "fuz", "baz"}
	prepareProvider(p, gameKeys)
	toDelete := []string{"foo", "fuz"}

	for _, s := range toDelete {
		err = p.DeleteGame(s)
		a.NoError(err)
		a.False(p.HasGame(s))
	}
	a.True(getProviderSize(p) == len(gameKeys)-len(toDelete))
}

func getProviderSize(p Provider) int {
	m, _ := p.GetAllGames()
	return len(m)
}

func prepareProvider(p Provider, ids []string) {
	for _, id := range ids {
		p.PutGame(&model.Game{Id: id})
	}
}

type redisMock struct {
	Provider Provider
	redis    *miniredis.Miniredis
}

func newMockedRedisProvider() *redisMock {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	p := NewRedisProvider(s.Addr(), "", 1)
	err = p.Connect()
	if err != nil {
		panic(err)
	}

	return &redisMock{
		Provider: p,
		redis:    s,
	}
}
