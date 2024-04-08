package main

import (
	"context"
	"fmt"
	"log"

	okex "github.com/pefish/go-okx"
	"github.com/pefish/go-okx/api"
	"github.com/pefish/go-okx/requests/rest/trade"
	"github.com/pkg/errors"
)

func main() {
	err := do()
	if err != nil {
		log.Fatal(err)
	}
}

func do() error {
	symbol := "BTC-USDT-SWAP"

	client, err := api.NewClient(
		context.Background(),
		"9c5760c6-ff0c-4e24-9bca-e60dc989bf46",
		"",
		":",
		okex.NormalServer,
	)
	if err != nil {
		return err
	}

	placeOrderRes, err := client.Rest.Trade.PlaceOrder([]trade.PlaceOrder{
		{
			InstID:  symbol,
			TdMode:  okex.TradeCrossMode,
			Side:    okex.OrderBuy,
			OrdType: okex.OrderMarket,
			Sz:      1,
		},
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
