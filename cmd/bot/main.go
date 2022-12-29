package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/yeyee2901/unitedb-discord-bot/pkg/config"
	"github.com/yeyee2901/unitedb-discord-bot/pkg/datasource"
	"github.com/yeyee2901/unitedb-discord-bot/pkg/discord"

	"github.com/go-redis/redis/v8"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	// INIT: datasource
	cfg := config.LoadConfig()
	initLogger(cfg)
	ds := datasource.NewDataSource(cfg, mustInitDB(cfg), mustInitRedis(cfg))

	// cleanup function with recovery to retrieve error
	defer func() {
		if err := recover(); err != nil {
			log.Error().Err(fmt.Errorf("%v", err)).Msg("EXIT.fatal")
		}

		if err := ds.DB.Close(); err != nil {
			log.Error().Err(err).Msg("EXIT.db")
		}

		if err := ds.Redis.Close(); err != nil {
			log.Error().Err(err).Msg("EXIT.redis")
		}

		fmt.Println("Bot exited.")
		log.Info().Msg("EXIT")
	}()

	// INIT: discord
	var (
		clientId     = mustReadFile_String(cfg.Discord.ClientIdFile)
		clientSecret = mustReadFile_String(cfg.Discord.ClientSecretFile)
		token        = mustReadFile_String(cfg.Discord.TokenFile)
	)

	dcSession, err := discordgo.New(fmt.Sprintf("Bot %s", token))
	if err != nil {
		panic(err)
	}

	dcBot := discord.NewDiscordBotService(clientId, clientSecret, token, ds, dcSession)
	dcBot.InitBot()

	// open connection to discord using websocket
	log.Info().Msg("START")
	fmt.Println("Bot Start")
	err = dcBot.Bot.Open()
	if err != nil {
		panic(err)
	}
	defer dcBot.Bot.Close()

	// quit signal
	osQuit := make(chan os.Signal)
	signal.Notify(osQuit, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)

	// wait for quit signal
	sig := <-osQuit
	log.Warn().Str("interrupt", fmt.Sprintf("Received signal %s", sig.String())).Msg("EXIT.interrupt")
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
		return strings.ReplaceAll(string(b), "\n", "")
	}
}

func initLogger(cfg *config.AppConfig) {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zerolog.TimeFieldFormat = "2006-01-02 15:04:05"
	log.Logger = zerolog.New(&lumberjack.Logger{
		Filename:   cfg.Discord.Logfile,
		MaxSize:    100,
		MaxBackups: 3,
		MaxAge:     30,
		Compress:   true,
	})
	log.Logger = log.With().Caller().Logger()
	log.Logger = log.With().Timestamp().Logger()
}
