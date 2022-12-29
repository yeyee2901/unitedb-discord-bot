package main

import (
	"context"

	"github.com/yeyee2901/unitedb-discord-bot/pkg/config"
	"github.com/yeyee2901/unitedb-discord-bot/pkg/datasource"
	"github.com/yeyee2901/unitedb-discord-bot/pkg/debug"

	"github.com/go-redis/redis/v8"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func main() {
	cfg := config.LoadConfig()
	ds := datasource.NewDataSource(cfg, initDB(cfg), initRedis(cfg))
	debug.DumpStruct(ds)
}

func initRedis(cfg *config.AppConfig) *redis.Client {
	r := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    cfg.Redis.Host + ":" + cfg.Redis.Port,
	})

	if err := r.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}

	return r
}

func initDB(cfg *config.AppConfig) *sqlx.DB {
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
