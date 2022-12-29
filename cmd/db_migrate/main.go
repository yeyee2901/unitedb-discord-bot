package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/yeyee2901/unitedb-discord-bot/pkg/config"
	"github.com/yeyee2901/unitedb-discord-bot/pkg/models"
)

func main() {
	cfg := config.LoadConfig()
	db := mustInitDB(cfg)

	InsertBattleItems(db, cfg)
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

func getUniteDbResource(urlEndpoint string, resp any) {
	httpResp, err := http.Get(urlEndpoint)
	if err != nil {
		log.Fatal(err)
	}

	if httpResp.StatusCode != http.StatusOK {
		log.Fatalf("%s returned HTTP %d", urlEndpoint, httpResp.StatusCode)
	}

	if err := json.NewDecoder(httpResp.Body).Decode(resp); err != nil {
		log.Fatal(err)
	}
}

func InsertBattleItems(db *sqlx.DB, cfg *config.AppConfig) {
	tx := db.MustBegin()

	// create table
	b, err := os.ReadFile("schema/pokemon_held_items.sql")
	if err != nil {
		log.Fatal(err.Error())
	}

	tx.MustExec(string(b))

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		log.Fatal(err.Error())
	}

	// fetch resources from the web API
	battleItems := []models.BattleItem{}
	getUniteDbResource(cfg.UniteDB.BaseURL+cfg.UniteDB.Endpoints.BattleItems, &battleItems)

	// insert the data
	q := `
        INSERT INTO pokemon_battle_items
            (name, description, tier, cooldown, trainer_level)
        VALUES
            (:name, :description, :tier, :cooldown, :trainer_level)
    `

	tx = db.MustBegin()
	if _, err := tx.NamedExec(q, battleItems); err != nil {
		tx.Rollback()
		log.Fatal(err.Error())
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		log.Fatal(err.Error())
	}
}
