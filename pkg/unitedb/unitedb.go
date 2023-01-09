package unitedb

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

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
	// check if:
	// - database is empty
	// - last update is more than update_interval_days
	shouldRefresh, err := u.shouldRefresh()
	if err != nil {
		log.Error().Err(err).Msg("unite-db.RefreshBattleItems.shouldRefresh")
		return err
	}

	if !shouldRefresh {
		return nil
	}

	// fetch resources from the web API
	battleItems := []BattleItem{}
	if err := u.getUniteDbResource(u.Config.UniteDB.BaseURL+u.Config.UniteDB.Endpoints.BattleItems, &battleItems); err != nil {
		log.Error().Err(err).Msg("unite-db.RefreshBattleItems.fetch")
		return err
	}

	return nil
}

func (*UniteDbService) getUniteDbResource(urlEndpoint string, resp any) error {
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

func (u *UniteDbService) shouldRefresh() (bool, error) {
	// fetch from db
	res, err := u.GetBattleItems("", "")
	if err != nil {
		return false, err
	}

	// check whether data is > 2 weeks old
	lastUpdate := res[0].LastUpdated
	twoWeeks := lastUpdate.AddDate(0, 0, 14)
	if time.Now().After(twoWeeks) {
		return true, nil
	} else {
		return false, nil
	}
}
