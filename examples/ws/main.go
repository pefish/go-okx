package main

import (
	"context"
	"fmt"
	okex "github.com/pefish/go-okx"
	"github.com/pefish/go-okx/api"
	"github.com/pefish/go-okx/events"
	"github.com/pefish/go-okx/events/public"
	"log"
)

func main() {
	err := do()
	if err != nil {
		log.Fatal(err)
	}
}

func do() error {
	handleClient, err := api.NewClient(
		context.Background(),
		"YOUR-API-KEY",
		"YOUR-SECRET-KEY",
		"YOUR-PASS-PHRASE",
		okex.CandleWsServer,
	)
	if err != nil {
		return err
	}

	client, err := api.NewClient(
		context.Background(),
		"YOUR-API-KEY",
		"YOUR-SECRET-KEY",
		"YOUR-PASS-PHRASE",
		okex.AwsServer,
	)
	if err != nil {
		return err
	}

	log.Println("Starting")

	handleClient.Ws.SubscribeChan = make(chan *events.Subscribe)
	handleClient.Ws.UnsubscribeCh = make(chan *events.Unsubscribe)
	handleClient.Ws.ErrChan = make(chan *events.Error)
	handleClient.Ws.Public.CandlesticksCh = make(chan *public.Candlesticks)
	//err = handleClient.Ws.Public.Subscribe(false, []map[string]string{
	//	map[string]string{
	//		"instId":  "BTC-USDT-SWAP",
	//		"channel": "candle5m",
	//	},
	//})
	//if err != nil {
	//	return err
	//}

	client.Ws.SubscribeChan = make(chan *events.Subscribe)
	client.Ws.UnsubscribeCh = make(chan *events.Unsubscribe)
	client.Ws.ErrChan = make(chan *events.Error)
	client.Ws.Public.FundingRateCh = make(chan *public.FundingRate)
	client.Ws.Public.OpenInterestCh = make(chan *public.OpenInterest)
	err = client.Ws.Public.Subscribe(false, []map[string]string{
		map[string]string{
			"channel": "funding-rate",
			"instId":  "BTC-USDT-SWAP",
		},
		map[string]string{
			"channel": "open-interest",
			"instId":  "BTC-USDT-SWAP",
		},
		map[string]string{
			"channel": "funding-rate",
			"instId":  "ETH-USDT-SWAP",
		},
		map[string]string{
			"channel": "open-interest",
			"instId":  "ETH-USDT-SWAP",
		},
	})
	if err != nil {
		return err
	}

	for {
		select {
		case sub := <-handleClient.Ws.SubscribeChan:
			channel, _ := sub.Arg.Get("channel")
			log.Printf("[Subscribed]\t%s", channel)
		case uSub := <-handleClient.Ws.UnsubscribeCh:
			channel, _ := uSub.Arg.Get("channel")
			log.Printf("[Unsubscribed]\t%s", channel)
		case err := <-handleClient.Ws.ErrChan:
			log.Printf("[Error]\t%+v", err)
			for _, datum := range err.Data {
				log.Printf("[Error]\t\t%+v", datum)
			}
		case sub := <-client.Ws.SubscribeChan:
			channel, _ := sub.Arg.Get("channel")
			log.Printf("[Subscribed]\t%s", channel)
		case uSub := <-client.Ws.UnsubscribeCh:
			channel, _ := uSub.Arg.Get("channel")
			log.Printf("[Unsubscribed]\t%s", channel)
		case err := <-client.Ws.ErrChan:
			log.Printf("[Error]\t%+v", err)
			for _, datum := range err.Data {
				log.Printf("[Error]\t\t%+v", datum)
			}
		case o := <-client.Ws.Public.OpenInterestCh:
			for _, p := range o.OpenInterests {
				fmt.Printf(`
symbol: %s
hold vol: %f
`,
					p.InstID,
					p.OiCcy,
				)
			}
		case f := <-client.Ws.Public.FundingRateCh:
			for _, p := range f.Rates {
				fmt.Printf(`
symbol: %s
funding rate: %f
`,
					p.InstID,
					p.FundingRate,
				)
			}
		case i := <-handleClient.Ws.Public.CandlesticksCh:
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
		}
	}

	return nil
}
