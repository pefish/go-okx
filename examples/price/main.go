package main

import (
	"context"
	"fmt"
	okex "github.com/pefish/go-okx"
	"github.com/pefish/go-okx/api"
	"github.com/pefish/go-okx/requests/rest/market"
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
