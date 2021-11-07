package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/l12u/gamemaster/internal/model"
)

const (
	// We use a seperate set in the Redis for each game
	// so that we can iterate over them if needed.
	gamesIterListKey = "games"

	// Each game is stored with this specific key format
	// into the Redis, where the unique id is appended
	// to the end.
	gameKeyFormat = "game:%s"
)

type RedisProvider struct {
	options *redis.Options
	client  *redis.Client
}

var ctx = context.Background()

func NewRedisProvider(addr string, pw string, db int) *RedisProvider {
	return &RedisProvider{
		options: &redis.Options{
			Addr:     addr,
			Password: pw,
			DB:       db,
		},
	}
}

func (r *RedisProvider) Connect() error {
	rdb := redis.NewClient(r.options)
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return err
	}

	r.client = rdb
	return nil
}

func (r *RedisProvider) PutGame(g *model.Game) error {
	if g == nil {
		return errors.New("game can not be nil")
	}

	m, err := json.Marshal(g)
	if err != nil {
		return err
	}
	str := string(m)

	err = r.client.Set(ctx, getKeyFromGame(g.Id), str, 0).Err()
	if err != nil {
		return err
	}

	// push to iter list
	err = r.client.SAdd(ctx, gamesIterListKey, g.Id).Err()
	return err
}

func (r *RedisProvider) GetGame(id string) (*model.Game, error) {
	str, err := r.client.Get(ctx, getKeyFromGame(id)).Result()
	if err != nil {
		return nil, err
	}

	var g model.Game
	err = json.Unmarshal([]byte(str), &g)
	if err != nil {
		return nil, err
	}
	return &g, nil
}

func (r *RedisProvider) GetAllGames() (model.GameMap, error) {
	gameIds, err := r.client.SMembers(ctx, gamesIterListKey).Result()
	if err != nil {
		return nil, err
	}
	if len(gameIds) == 0 {
		return model.EmptyGameMap, nil
	}
	gameKeys := getKeysFromGameIds(gameIds)

	res, err := r.client.MGet(ctx, gameKeys...).Result()
	if err != nil {
		return nil, err
	}

	gameMap := make(model.GameMap, len(res))
	for _, r := range res {
		str, ok := r.(string)
		if !ok {
			continue
		}

		var g model.Game
		err = json.Unmarshal([]byte(str), &g)
		if err != nil {
			return nil, err
		}
		gameMap[g.Id] = &g
	}

	return gameMap, nil
}

func (r *RedisProvider) DeleteGame(id string) error {
	err := r.client.Del(ctx, getKeyFromGame(id)).Err()
	if err != nil {
		return err
	}

	// remove from iter list
	err = r.client.SRem(ctx, gamesIterListKey, id).Err()
	return err
}

func (r *RedisProvider) ClearGames() error {
	gameIds, err := r.client.SMembers(ctx, gamesIterListKey).Result()
	if err != nil {
		return err
	}
	if len(gameIds) == 0 {
		// we don't have to do anything
		return nil
	}
	gameKeys := getKeysFromGameIds(gameIds)

	// delete all single game entries from Redis
	err = r.client.Del(ctx, gameKeys...).Err()
	if err != nil {
		return err
	}

	// delete iter list
	err = r.client.Del(ctx, gamesIterListKey).Err()
	return err
}

func (r *RedisProvider) HasGame(id string) (bool, error) {
	code, err := r.client.Exists(ctx, getKeyFromGame(id)).Result()
	return code == 1, err
}

func getKeyFromGame(id string) string {
	return fmt.Sprintf(gameKeyFormat, id)
}

func getKeysFromGameIds(ids []string) []string {
	gameKeys := make([]string, len(ids))
	for _, gid := range ids {
		gameKeys = append(gameKeys, getKeyFromGame(gid))
	}
	return gameKeys
}
