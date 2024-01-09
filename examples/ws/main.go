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

	client.Ws.SubscribeChan = make(chan *events.Subscribe)
	client.Ws.UnsubscribeCh = make(chan *events.Unsubscribe)
	client.Ws.ErrChan = make(chan *events.Error)
	client.Ws.Public.FundingRateCh = make(chan *public.FundingRate)
	client.Ws.Public.OpenInterestCh = make(chan *public.OpenInterest)
	client.Ws.Public.MarkPriceCh = make(chan *public.MarkPrice)
	client.Ws.Public.LiquidationOrdersCh = make(chan *public.LiquidationOrders)
	err = client.Ws.Public.Subscribe(false, []map[string]string{
		//map[string]string{
		//	"channel": "funding-rate",
		//	"instId":  "BTC-USDT-SWAP",
		//},
		//map[string]string{
		//	"channel": "open-interest",
		//	"instId":  "BTC-USDT-SWAP",
		//},
		//map[string]string{
		//	"channel": "funding-rate",
		//	"instId":  "ETH-USDT-SWAP",
		//},
		//map[string]string{
		//	"channel": "open-interest",
		//	"instId":  "ETH-USDT-SWAP",
		//},
		map[string]string{
			"channel": "mark-price",
			"instId":  "BTC-USDT",
		},
		map[string]string{
			"channel":  "liquidation-orders",
			"instType": "SWAP",
		},
	})
	if err != nil {
		return err
	}

	for {
		select {
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
time: %s
`,
					p.InstID,
					p.OiCcy,
					p.TS.String(),
				)
			}
		case f := <-client.Ws.Public.FundingRateCh:
			for _, p := range f.Rates {
				fmt.Printf(`
symbol: %s
funding rate: %f
time: %s
`,
					p.InstID,
					p.FundingRate,
					p.FundingTime.String(),
				)
			}
		case f := <-client.Ws.Public.MarkPriceCh:
			for _, p := range f.Prices {
				fmt.Printf(`
symbol: %s
price: %f
time: %s
`,
					p.InstID,
					p.MarkPx,
					p.TS.String(),
				)
			}
			//err := client.Ws.Public.Unsubscribe(false, []map[string]string{
			//	map[string]string{
			//		"channel": "mark-price",
			//		"instId":  "BTC-USDT",
			//	},
			//})
			//if err != nil {
			//	return err
			//}
		case f := <-client.Ws.Public.LiquidationOrdersCh:
			for _, p := range f.LiquidationOrders {
				fmt.Printf("symbol: %s\n", p.InstID)
				for _, d := range p.Details {
					fmt.Printf("\tside: %s, price: %f, quantity: %f, quantityU: %f, time: %s\n", d.PosSide, d.BkPx, d.Sz, d.BkPx*d.Sz, d.TS.String())
				}
			}
		}
	}

	return nil
}
