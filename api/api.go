package api

import (
	"context"

	i_logger "github.com/pefish/go-interface/i-logger"
	okex "github.com/pefish/go-okx"
	"github.com/pefish/go-okx/api/rest"
	"github.com/pefish/go-okx/api/ws"
)

// Client is the main api wrapper of okex
type Client struct {
	Rest   *rest.ClientRest
	Ws     *ws.ClientWs
	ctx    context.Context
	logger i_logger.ILogger
}

// NewClient returns a pointer to a fresh Client
func NewClient(
	ctx context.Context,
	logger i_logger.ILogger,
	apiKey,
	secretKey,
	passphrase string,
	destination okex.Destination,
) (*Client, error) {
	restURL := okex.RestURL
	wsPubURL := okex.PublicWsURL
	wsPriURL := okex.PrivateWsURL
	switch destination {
	case okex.AwsServer:
		restURL = okex.AwsRestURL
		wsPubURL = okex.AwsPublicWsURL
		wsPriURL = okex.AwsPrivateWsURL
	case okex.DemoServer:
		restURL = okex.DemoRestURL
		wsPubURL = okex.DemoPublicWsURL
		wsPriURL = okex.DemoPrivateWsURL
	case okex.CandleWsServer:
		restURL = okex.AwsRestURL
		wsPubURL = okex.HandleWsURL
		wsPriURL = okex.AwsPrivateWsURL
	}

	r := rest.NewClient(logger, apiKey, secretKey, passphrase, restURL, destination)
	c := ws.NewClient(ctx, logger, apiKey, secretKey, passphrase, map[bool]okex.BaseURL{true: wsPriURL, false: wsPubURL})

	return &Client{
		Rest:   r,
		Ws:     c,
		ctx:    ctx,
		logger: logger,
	}, nil
}
