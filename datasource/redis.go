package datasource

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/yeyee2901/unitedb-discord-bot/config"
)

// MustConnectRedisNormalMode connects to normal redis instance (not sentinel
// mode). panics on error
func MustConnectRedisNormalMode(cfg *config.RedisMeta) *redis.Client {
	r := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    cfg.Host,
	})

	if err := r.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}

	return r
}
