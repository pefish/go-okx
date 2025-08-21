package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	i_logger "github.com/pefish/go-interface/i-logger"
	okex "github.com/pefish/go-okx"
	"github.com/pefish/go-okx/api"
	"github.com/pkg/errors"
)

var toAddress = ""
var amount = 0.1

func main() {
	envMap, _ := godotenv.Read("./.env")
	for k, v := range envMap {
		os.Setenv(k, v)
	}

	err := do()
	if err != nil {
		log.Fatalf("%+v", err)
	}
}

func do() error {
	client, err := api.NewClient(
		context.Background(),
		&i_logger.DefaultLogger,
		os.Getenv("API_KEY"),
		os.Getenv("API_SECRET"),
		os.Getenv("API_PASS"),
		okex.NormalServer,
	)
	if err != nil {
		return err
	}

	getCurrenciesRes, err := client.Rest.Funding.GetCurrencies()
	if err != nil {
		return err
	}

	if getCurrenciesRes.Code != 0 {
		return errors.Errorf("getCurrencies failed. err: %s, code: %d", getCurrenciesRes.Msg, getCurrenciesRes.Code)
	}

	for _, currencyInfo := range getCurrenciesRes.Currencies {
		if currencyInfo.Ccy == "SOL" {
			fmt.Printf(`
<CCY: %s>
<Chain: %s>
<Name: %s>
`,
				currencyInfo.Ccy,
				currencyInfo.Chain,
				currencyInfo.Name,
			)
		}

	}

	return nil
}
