package api

import (
	"context"
	"net/http"

	jsonrpc "github.com/filecoin-project/go-jsonrpc"
	lotusapi "github.com/filecoin-project/lotus/api"
)

type LotusClient struct {
	lotusapi.FullNodeStruct
	Closer jsonrpc.ClientCloser
}

func NewLotusClient(ctx context.Context, addr, authToken string) (*LotusClient, error) {
	var headers http.Header
	if authToken != "" {
		headers = http.Header{"Authorization": []string{"Bearer " + authToken}}
	}
	var api lotusapi.FullNodeStruct
	closer, err := jsonrpc.NewMergeClient(ctx, addr, "Filecoin",
		[]interface{}{&api.Internal, &api.CommonStruct.Internal}, headers)
	if err != nil {
		return nil, err
	}

	return &LotusClient{
		api,
		closer,
	}, nil
}

func (c *LotusClient) Close() {
	if c.Closer != nil {
		c.Closer()
	}
}
