package unitedb

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/yeyee2901/unitedb-discord-bot/pkg/datasource"
)

type UniteDbService struct {
	*datasource.DataSource
}

func NewUniteDbService(ds *datasource.DataSource) *UniteDbService {
	return &UniteDbService{ds}
}

func (u *UniteDbService) RefreshBattleItems() error {
	// fetch resources from the web API
	battleItems := []BattleItem{}
	err := getUniteDbResource(u.Config.UniteDB.BaseURL+u.Config.UniteDB.Endpoints.BattleItems, &battleItems)
	if err != nil {
		log.Error().Err(err).Msg("unite-db.getUniteDbResource")
		return err
	}

	return nil
}

func getUniteDbResource(urlEndpoint string, resp any) error {
	httpResp, err := http.Get(urlEndpoint)
	if err != nil {
		return err
	}

	if httpResp.StatusCode != http.StatusOK {
		err := fmt.Errorf("%s returned HTTP %s", urlEndpoint, httpResp.Status)
		return err
	}

	if err := json.NewDecoder(httpResp.Body).Decode(resp); err != nil {
		return err
	}

	return nil
}
