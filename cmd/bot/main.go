package main

import (
	"context"
	"os"

	"github.com/yeyee2901/unitedb-discord-bot/pkg/config"
	"github.com/yeyee2901/unitedb-discord-bot/pkg/datasource"
	"github.com/yeyee2901/unitedb-discord-bot/pkg/debug"
	"github.com/yeyee2901/unitedb-discord-bot/pkg/discord"

	"github.com/go-redis/redis/v8"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func main() {
	// INIT: datasource
	cfg := config.LoadConfig()
	ds := datasource.NewDataSource(cfg, mustInitDB(cfg), mustInitRedis(cfg))

	// INIT: discord
	var (
		clientId     = mustReadFile_String(cfg.Discord.ClientIdFile)
		clientSecret = mustReadFile_String(cfg.Discord.ClientSecretFile)
		token        = mustReadFile_String(cfg.Discord.TokenFile)
	)

	dcBot := discord.NewDiscordBotService(clientId, clientSecret, token, ds)
	debug.DumpStruct(dcBot)
}

// will panic on failure
func mustInitRedis(cfg *config.AppConfig) *redis.Client {
	r := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    cfg.Redis.Host + ":" + cfg.Redis.Port,
	})

	if err := r.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}

	return r
}

// will panic on failure
func mustInitDB(cfg *config.AppConfig) *sqlx.DB {
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

// will panic on failure or if the file is empty
func mustReadFile_String(filepath string) string {
	if b, err := os.ReadFile(filepath); err != nil {
		panic(err)
	} else {
		if len(string(b)) == 0 {
			panic(err)
		}
		return string(b)
	}
}
