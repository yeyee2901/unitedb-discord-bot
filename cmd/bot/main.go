// Package main ...
// nolint
package main

import (
	"bytes"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/yeyee2901/unitedb-discord-bot/config"
	"github.com/yeyee2901/unitedb-discord-bot/datasource"
	"github.com/yeyee2901/unitedb-discord-bot/discord"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	// load config file
	var path string
	if os.Getenv("CONFIG") == "" {
		path = "setting.yaml"
	}
	cfg := config.MustLoadConfig(path)

	// init logger
	var logger zerolog.Logger
	if cfg.Bot.Mode == "production" {
		logger = zerolog.New(&lumberjack.Logger{
			Filename: "log/zerolog.log",
			Compress: true,
		})
		logger = logger.With().Caller().Logger()
		logger = logger.With().Timestamp().Logger()
	} else {
		logger = zerolog.New(os.Stdout)
		logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stdout})
		logger = logger.With().Caller().Logger()
		logger = logger.With().Timestamp().Logger()
	}

	// open the connection
	service, session, err := discord.NewDiscordBotService(
		mustLoadFile(cfg.Bot.ClientIDFile, true),
		mustLoadFile(cfg.Bot.ClientSecretFile, true),
		mustLoadFile(cfg.Bot.TokenFile, true),
		datasource.NewRedisStore(&cfg.Redis),
		&logger,
		cfg,
	)
	if err != nil {
		log.Fatalln(err)
	}

	// register commands
	cmdIDs, err := service.RegisterCommands(session)
	if err != nil {
		log.Fatalln(err)
	}

	// wait indefinitely
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	<-sig

	// cleanup
	err = service.UnregisterCommands(session, cmdIDs...)
	if err != nil {
		log.Fatalln(err)
	}

	service.Close(session)
}

func mustLoadFile(path string, stripNewLine bool) string {
	b, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	if stripNewLine {
		b = bytes.ReplaceAll(b, []byte{'\n'}, []byte{})
	}

	return string(b)
}
