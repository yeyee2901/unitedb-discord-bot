package datasource

import (
	"github.com/go-sql-driver/mysql"
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

func MustInitDB(cfg *config.AppConfig) *sqlx.DB {
	// will panic on failure
	dbConfig := mysql.Config{
		User:                 cfg.DB.User,
		Passwd:               cfg.DB.Password,
		Net:                  "tcp",
		Addr:                 cfg.DB.Host,
		DBName:               cfg.DB.Database,
		AllowNativePasswords: true,
		CheckConnLiveness:    true,
	}

	d := sqlx.MustConnect("mysql", dbConfig.FormatDSN())

	if err := d.Ping(); err != nil {
		panic(err)
	}

	return d
}
