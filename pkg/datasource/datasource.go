package datasource

import (
	"github.com/yeyee2901/unitedb-discord-bot/pkg/config"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
)

type DataSource struct {
	Config *config.AppConfig
	DB     *sqlx.DB
	Redis  *redis.Client
}

func NewDataSource(cfg *config.AppConfig, db *sqlx.DB, r *redis.Client) *DataSource {
	return &DataSource{cfg, db, r}
}
