package main

import (
	"context"
	"fmt"
	okex "github.com/pefish/go-okx"
	"github.com/pefish/go-okx/api"
	"github.com/pefish/go-okx/requests/rest/public"
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
	res, err := client.Rest.PublicData.GetFundingRate(public.GetFundingRate{
		InstID: "ZRX-USDT-SWAP",
	})
	if err != nil {
		return err
	}
	for _, fundingRate := range res.FundingRates {
		fmt.Printf(
			`
产品id: %s
资金费率: %f
`,
			fundingRate.InstID,
			fundingRate.FundingRate,
		)
	}

	return nil
}
