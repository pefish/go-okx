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
	"github.com/pefish/go-okx/requests/rest/funding"
	"github.com/pefish/go-okx/requests/rest/trade"
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

	placeOrderRes, err := client.Rest.Funding.Withdrawal(funding.Withdrawal{
		Ccy:    "SOL",
		Chain:  "sol",
		ToAddr: toAddress,
		Amt:    amount,
		Dest:   okex.WithdrawalDigitalAddressDestination,
	})
	if err != nil {
		return err
	}

	if placeOrderRes.Code != 0 {
		return errors.Errorf("PlaceOrder failed. err: %s, code: %d", placeOrderRes.Msg, placeOrderRes.Code)
	}

	fmt.Println(placeOrderRes.PlaceOrders[0].OrdID)

	getOrderDetailRes, err := client.Rest.Trade.GetOrderDetail(trade.OrderDetails{
		InstID: symbol,
		OrdID:  placeOrderRes.PlaceOrders[0].OrdID,
	})
	if err != nil {
		return err
	}
	if getOrderDetailRes.Code != 0 {
		return errors.Errorf("GetOrderDetail failed. err: %s, code: %d", getOrderDetailRes.Msg, getOrderDetailRes.Code)
	}

	fmt.Printf(`
	AvgPx: %f
	`,
		float64(getOrderDetailRes.Orders[0].AvgPx),
	)

	return nil
}
