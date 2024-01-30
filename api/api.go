package api

import (
	"context"
	go_logger "github.com/pefish/go-logger"
	"github.com/pefish/go-okx"
	"github.com/pefish/go-okx/api/rest"
	"github.com/pefish/go-okx/api/ws"
)

// Client is the main api wrapper of okex
type Client struct {
	Rest   *rest.ClientRest
	Ws     *ws.ClientWs
	ctx    context.Context
	logger go_logger.InterfaceLogger
}

// NewClient returns a pointer to a fresh Client
func NewClient(ctx context.Context, apiKey, secretKey, passphrase string, destination okex.Destination) (*Client, error) {
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

	logger := go_logger.Logger

	r := rest.NewClient(apiKey, secretKey, passphrase, restURL, destination).SetLogger(logger)
	c := ws.NewClient(ctx, apiKey, secretKey, passphrase, map[bool]okex.BaseURL{true: wsPriURL, false: wsPubURL}).SetLogger(logger)

	return &Client{
		Rest:   r,
		Ws:     c,
		ctx:    ctx,
		logger: logger,
	}, nil
}

func (c *Client) SetLogger(logger go_logger.InterfaceLogger) *Client {
	c.logger = logger
	c.Rest = c.Rest.SetLogger(logger)
	c.Ws = c.Ws.SetLogger(logger)
	return c
}
