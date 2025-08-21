package main

import (
	"context"
	"fmt"
	"log"

	i_logger "github.com/pefish/go-interface/i-logger"
	okex "github.com/pefish/go-okx"
	"github.com/pefish/go-okx/api"
	"github.com/pefish/go-okx/requests/rest/tradedata"
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
		&i_logger.DefaultLogger,
		"YOUR-API-KEY",
		"YOUR-SECRET-KEY",
		"YOUR-PASS-PHRASE",
		okex.NormalServer,
	)
	if err != nil {
		return err
	}
	res, err := client.Rest.TradeData.GetHoldVolLongShortRatio(tradedata.GetHoldVolRatio{
		Ccy:    "BTC",
		Period: "5m",
		Num:    1,
	})
	if err != nil {
		return err
	}

	for _, ratio := range res.Ratios {
		fmt.Printf(
			`
time: %s
long: %f
short: %f
`,
			&ratio.TS,
			ratio.LongRatio,
			ratio.ShortRatio,
		)
	}

	return nil
}
