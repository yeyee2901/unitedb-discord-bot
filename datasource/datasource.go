package datasource

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/yeyee2901/unitedb-discord-bot/config"
)

// RedisStore is interface abstraction for communicating with redis
type RedisStore interface {
	// GetWithContext get using provided context
	GetWithContext(c context.Context, key string) (string, error)

	// SaveWithContext save using provided context
	SaveWithContext(c context.Context, key string, data string, exp time.Duration) error

	// Get retrieves the string data identified by the key. This method uses
	// context background internally. If you need to use external context, use
	// GetWithContext instead.
	Get(key string) (string, error)

	// Save retrieves the string data identified by the key. This method uses
	// context background internally. If you need to use external context, use
	// SaveWithContext instead.
	Save(key string, data string, exp time.Duration) error
}

type redisStoreImpl struct {
	client *redis.Client
}

// Get retrieves the string data identified by the key. This method uses
// context background internally. If you need to use external context, use
// GetWithContext instead.
func (re *redisStoreImpl) Get(key string) (string, error) {
	return re.GetWithContext(context.Background(), key)
}

// Save retrieves the string data identified by the key. This method uses
// context background internally. If you need to use external context, use
// SaveWithContext instead.
func (re *redisStoreImpl) Save(key string, data string, exp time.Duration) error {
	return re.SaveWithContext(context.Background(), key, data, exp)
}

// Get retrieves the string data identified by the key
func (r *redisStoreImpl) GetWithContext(c context.Context, key string) (string, error) {
	return r.client.Get(c, key).Result()
}

// Save saves the string data to the redis client using the provided key
// & expiration time
func (r *redisStoreImpl) SaveWithContext(c context.Context, key string, data string, exp time.Duration) error {
	return r.client.Set(c, key, data, exp).Err()
}

// NewRedisStore constructs a new RedisStore interface for communicating with
// the redis
func NewRedisStore(cfg *config.RedisMeta) RedisStore {
	return &redisStoreImpl{
		client: MustConnectRedisNormalMode(cfg),
	}
}
