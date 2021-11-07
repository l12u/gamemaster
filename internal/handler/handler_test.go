package handler

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/l12u/gamemaster/internal/config"
	"github.com/l12u/gamemaster/internal/model"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TODO use gin test context to test http requests
//
// DeleteGame
// - delete existent game
// - delete non-existent game

func TestPostGame(t *testing.T) {
	SetupBoardConfig(&config.BoardConfig{Boards: []*config.Board{
		{Type: "hangman"},
		{Type: "chess"},
	}})
	SetupProvider()
	gin.SetMode(gin.TestMode)

	testPostGameValid(t)
	testPostGameInvalid(t)
}

func TestDeleteGame(t *testing.T) {
	SetupBoardConfig(&config.BoardConfig{Boards: []*config.Board{
		{Type: "hangman"},
	}})
	SetupProvider()
	gin.SetMode(gin.TestMode)

	a := assert.New(t)

	// prepare multiple games to POST
	games := []*model.Game{
		createGame("hangman", "alice", "bob", "cicero", "dennis"),
	}
	bodies := getBodies(a, games)

	existentId := ""
	for _, body := range bodies {
		w := postGameRequest(body)
		a.EqualValues(http.StatusOK, w.Code)

		// check if the given id now exists again

		resp, err := getJsonResponse(w)
		a.NoError(err)

		id := resp["id"].(string)
		existentId = id

		w = getGameRequest(id)
		a.EqualValues(http.StatusOK, w.Code)
	}

	w := deleteGameRequest(existentId)
	a.EqualValues(http.StatusOK, w.Code)

	w = deleteGameRequest(existentId)
	a.False(w.Code == http.StatusOK)
}

func testPostGameValid(t *testing.T) {
	a := assert.New(t)

	// prepare multiple games to POST
	games := []*model.Game{
		createGame("hangman", "alice", "bob", "cicero", "dennis"),
		createGame("chess", "benedikt", "fabian"),
	}
	bodies := getBodies(a, games)

	ids := make(map[string]bool, len(bodies))

	for _, body := range bodies {
		w := postGameRequest(body)
		a.EqualValues(http.StatusOK, w.Code)

		// check if the given id now exists again

		resp, err := getJsonResponse(w)
		a.NoError(err)

		id := resp["id"].(string)
		ids[id] = true

		w = getGameRequest(id)
		a.EqualValues(http.StatusOK, w.Code)
	}

	// check if there were any duplicates
	a.True(len(ids) == len(games))
}

func testPostGameInvalid(t *testing.T) {
	a := assert.New(t)

	// prepare multiple games to POST
	games := []*model.Game{
		createGame("some_type", "foo", "bar"),
		createGame("chess", "foo", "foo", "bar"),
		createGameNoHost("chess", "foo", "bar", "baz"),
		createGameAllHost("chess", "foo", "bar", "baz"),
		createGameCustomHostRole("chess", "invalid_role", "foo", "bar", "baz"),
	}
	bodies := getBodies(a, games)

	for i, body := range bodies {
		w := postGameRequest(body)

		a.True(w.Code != http.StatusOK, "Result of test %v: %s", i, w.Body.String())
	}
}

func deleteGameRequest(gameId string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest(http.MethodDelete, "", nil)
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: gameId},
	}

	DeleteGame(c)

	return w
}

func postGameRequest(body io.Reader) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest(http.MethodPost, "", body)
	PostGame(c)

	return w
}

func getGameRequest(gameId string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest(http.MethodGet, "", nil)
	c.Params = gin.Params{
		gin.Param{Key: "id", Value: gameId},
	}

	GetGame(c)

	return w
}

func createGame(t string, playerIds ...string) *model.Game {
	players := make([]*model.Player, len(playerIds))
	for i, id := range playerIds {
		role := model.RolePlayer
		if i == 0 {
			role = model.RoleHost
		}

		players[i] = &model.Player{
			Id:   id,
			Name: "",
			Role: role,
		}
	}
	return &model.Game{
		Type:    t,
		Players: players,
	}
}

func createGameNoHost(t string, playerIds ...string) *model.Game {
	players := make([]*model.Player, len(playerIds))
	for i, id := range playerIds {
		players[i] = &model.Player{
			Id:   id,
			Name: "",
			Role: model.RolePlayer,
		}
	}
	return &model.Game{
		Type:    t,
		Players: players,
	}
}

func createGameAllHost(t string, playerIds ...string) *model.Game {
	players := make([]*model.Player, len(playerIds))
	for i, id := range playerIds {
		players[i] = &model.Player{
			Id:   id,
			Name: "",
			Role: model.RoleHost,
		}
	}
	return &model.Game{
		Type:    t,
		Players: players,
	}
}

func createGameCustomHostRole(t string, role string, playerIds ...string) *model.Game {
	players := make([]*model.Player, len(playerIds))
	for i, id := range playerIds {
		players[i] = &model.Player{
			Id:   id,
			Name: "",
			Role: role,
		}
	}
	return &model.Game{
		Type:    t,
		Players: players,
	}
}

func getJsonResponse(w *httptest.ResponseRecorder) (gin.H, error) {
	var got gin.H
	err := json.Unmarshal(w.Body.Bytes(), &got)
	if err != nil {
		return nil, err
	}
	return got, nil
}

func getBodies(a *assert.Assertions, games []*model.Game) []io.Reader {
	bodies := make([]io.Reader, len(games))
	for i, game := range games {
		j, err := json.Marshal(game)
		a.NoError(err)

		bodies[i] = io.NopCloser(bytes.NewBuffer(j))
	}
	return bodies
}
