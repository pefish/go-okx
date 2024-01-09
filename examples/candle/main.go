package main

import (
	"context"
	"fmt"
	okex "github.com/pefish/go-okx"
	"github.com/pefish/go-okx/api"
	"github.com/pefish/go-okx/events"
	"github.com/pefish/go-okx/events/public"
	"github.com/pefish/go-okx/requests/rest/market"
	public2 "github.com/pefish/go-okx/requests/ws/public"
	"log"
)

func main() {
	//err := do()
	err := doWs()
	if err != nil {
		log.Fatal(err)
	}
}

func do() error {
	client, err := api.NewClient(
		context.Background(),
		"YOUR-API-KEY",
		"YOUR-SECRET-KEY",
		"YOUR-PASS-PHRASE",
		okex.NormalServer,
	)
	if err != nil {
		return err
	}
	res, err := client.Rest.Market.GetCandlesticks(market.GetCandlesticks{
		InstID: "XLM-USDT-SWAP",
		Limit:  2,
		Bar:    "5m",
	})
	if err != nil {
		return err
	}
	for _, candle := range res.Candles {
		fmt.Printf(
			`
high: %f
open: %f
low: %f
close: %f
finished: %#v
`,
			candle.H,
			candle.O,
			candle.L,
			candle.C,
			candle.Confirm,
		)
	}

	return nil
}

func doWs() error {
	client, err := api.NewClient(
		context.Background(),
		"YOUR-API-KEY",
		"YOUR-SECRET-KEY",
		"YOUR-PASS-PHRASE",
		okex.CandleWsServer,
	)
	if err != nil {
		return err
	}

	log.Println("Starting")
	errChan := make(chan *events.Error)
	subChan := make(chan *events.Subscribe)
	uSubChan := make(chan *events.Unsubscribe)
	logChan := make(chan *events.Login)
	sucChan := make(chan *events.Success)
	client.Ws.SetChannels(errChan, subChan, uSubChan, logChan, sucChan)

	obCh := make(chan *public.Candlesticks)
	err = client.Ws.Public.Candlesticks([]public2.Candlesticks{
		{
			InstID:  "BTC-USDT-SWAP",
			Channel: "candle5m",
		},
		{
			InstID:  "ETH-USDT-SWAP",
			Channel: "candle5m",
		},
	}, obCh)
	if err != nil {
		return err
	}

	for {
		select {
		case <-logChan:
			log.Print("[Authorized]")
		case success := <-sucChan:
			log.Printf("[SUCCESS]\t%+v", success)
		case sub := <-subChan:
			channel, _ := sub.Arg.Get("channel")
			log.Printf("[Subscribed]\t%s", channel)
		case uSub := <-uSubChan:
			channel, _ := uSub.Arg.Get("channel")
			log.Printf("[Unsubscribed]\t%s", channel)
		case err := <-client.Ws.ErrChan:
			log.Printf("[Error]\t%+v", err)
			for _, datum := range err.Data {
				log.Printf("[Error]\t\t%+v", datum)
			}
		case i := <-obCh:
			symbol, _ := i.Arg.Get("instId")
			for _, p := range i.Candles {
				fmt.Printf(
					`
symbol: %s
high: %f
open: %f
low: %f
close: %f
time: %s
`,
					symbol,
					p.H,
					p.O,
					p.L,
					p.C,
					p.TS.String(),
				)
			}
		case b := <-client.Ws.DoneChan:
			log.Printf("[End]:\t%v", b)
			return nil
		}
	}

	return nil
}
