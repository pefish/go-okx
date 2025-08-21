package main

import (
	"context"
	"fmt"
	"log"

	i_logger "github.com/pefish/go-interface/i-logger"
	okex "github.com/pefish/go-okx"
	"github.com/pefish/go-okx/api"
	okx_requests_account "github.com/pefish/go-okx/requests/rest/account"
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
		":",
		okex.NormalServer,
	)
	if err != nil {
		return err
	}
	getBillsRes, err := client.Rest.Account.GetBills(okx_requests_account.GetBills{
		InstType: okex.SwapInstrument,
		Type:     okex.BillTradeType,
		Begin:    1711710001063,
	}, true)
	if err != nil {
		return err
	}
	if getBillsRes.Code != 0 {
		return errors.New(fmt.Sprintf("GetBills error. err: %s", getBillsRes.Msg))
	}

	for _, p := range getBillsRes.Bills {
		fmt.Printf(
			"symbol: %s, pnl: %f, fee: %f, PosBal: %f, PosBalChg: %f, SubType: %d\n",
			p.InstID,
			p.Pnl,
			p.Fee,
			p.PosBal,
			p.PosBalChg,
			p.SubType,
		)
	}

	return nil
}
