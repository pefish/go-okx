package main

import (
	"context"
	"fmt"
	"log"

	i_logger "github.com/pefish/go-interface/i-logger"
	okex "github.com/pefish/go-okx"
	"github.com/pefish/go-okx/api"
	"github.com/pefish/go-okx/requests/rest/account"
	"github.com/pkg/errors"
)

func main() {
	err := do()
	if err != nil {
		log.Fatal(err)
	}
}

func do() error {
	// symbol := "SHIB-USDT-SWAP"

	client, err := api.NewClient(
		context.Background(),
		&i_logger.DefaultLogger,
		"9c5760c6-ff0c-4e24-9bca-e60dc989bf46",
		"",
		"",
		okex.NormalServer,
	)
	if err != nil {
		return err
	}
	getPositionsRes, err := client.Rest.Account.GetPositions(account.GetPositions{
		// InstID:   []string{symbol},
		InstType: okex.SwapInstrument,
	})
	if err != nil {
		return err
	}
	if getPositionsRes.Code != 0 {
		return errors.Errorf("GetPositions failed. err: %s", getPositionsRes.Msg)
	}
	for _, p := range getPositionsRes.Positions {
		fmt.Printf("symbol: %s, pos: %f\n", p.InstID, p.Pos)
	}

	return nil
}
