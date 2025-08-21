package main

import (
	"context"
	"fmt"
	"log"

	i_logger "github.com/pefish/go-interface/i-logger"
	okex "github.com/pefish/go-okx"
	"github.com/pefish/go-okx/api"
	"github.com/pefish/go-okx/requests/rest/public"
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
	res, err := client.Rest.PublicData.GetOpenInterest(public.GetOpenInterest{
		InstType: "SWAP",
		InstID:   "BTC-USDT-SWAP",
	})
	if err != nil {
		return err
	}
	for _, openInterest := range res.OpenInterests {
		fmt.Printf(
			`
产品id: %s
持仓量: %f
`,
			openInterest.InstID,
			openInterest.OiCcy,
		)
	}

	return nil
}
