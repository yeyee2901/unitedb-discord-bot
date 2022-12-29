package datasource

import (
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/yeyee2901/unitedb-discord-bot/pkg/config"
	"github.com/yeyee2901/unitedb-discord-bot/pkg/debug"
)

func TestDB_GetBattleItems(t *testing.T) {
	cfg := config.LoadConfig()
	ds := NewDataSource(cfg, MustInitDB(cfg), nil)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// test without filter
	log.Info().Msg("Without filter")
	if res, err := ds.GetBattleItemsByName("", ""); err != nil {
		t.Fatal(err)
	} else {
		debug.DumpStruct(res)
	}

	log.Info().Msg("With name filter")
	if res, err := ds.GetBattleItemsByName("eject", ""); err != nil {
		t.Fatal(err)
	} else {
		debug.DumpStruct(res)
	}

	log.Info().Msg("With tier filter")
	if res, err := ds.GetBattleItemsByName("", "S"); err != nil {
		t.Fatal(err)
	} else {
		debug.DumpStruct(res)
	}

	log.Info().Msg("With both filter")
	if res, err := ds.GetBattleItemsByName("X", "S"); err != nil {
		t.Fatal(err)
	} else {
		debug.DumpStruct(res)
	}
}
