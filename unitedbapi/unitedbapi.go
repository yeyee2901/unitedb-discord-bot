package unitedbapi

import (
	"context"

	unitepb "github.com/yeyee2901/unitedb-api-proto/gen/go/unitedb/v1"
	"github.com/yeyee2901/unitedb-discord-bot/config"
	"google.golang.org/grpc"
)

type UniteDBAPI interface {
	// GetBattleItem is used to fetch battle item from the server
	GetBattleItem(context.Context, *unitepb.GetBattleItemRequest) (*unitepb.GetBattleItemResponse, error)
}

type grpcUniteDBAPI struct {
	conn   *grpc.ClientConn
	config *config.APIMeta
}

// NewGrpcUniteDBAPI constructs new instance of untiedb API & initiate the grpc connection to the API
func NewGrpcUniteDBAPI(cfg *config.APIMeta, opts ...grpc.DialOption) (UniteDBAPI, error) {
	conn, err := grpc.Dial(cfg.Host, opts...)
	if err != nil {
		return nil, err
	}

	return &grpcUniteDBAPI{
		conn:   conn,
		config: cfg,
	}, nil
}

// GetBattleItem implements UniteDBAPI
func (udb *grpcUniteDBAPI) GetBattleItem(c context.Context, req *unitepb.GetBattleItemRequest) (*unitepb.GetBattleItemResponse, error) {
	client := unitepb.NewUniteDBServiceClient(udb.conn)
	return client.GetBattleItem(c, req)
}
