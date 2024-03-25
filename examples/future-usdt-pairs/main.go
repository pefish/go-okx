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
	res, err := client.Rest.PublicData.GetInstruments(public.GetInstruments{
		InstType: "SWAP",
	})
	if err != nil {
		return err
	}
	for _, instrument := range res.Instruments {
		if instrument.CtType != okex.ContractLinearType { // 不是正向合约（U 本位合约）。
			continue
		}
		if instrument.State != "live" { // 不在交易中
			continue
		}
		if instrument.SettleCcy != "USDT" {
			continue
		}
		if instrument.InstID != "SHIB-USDT-SWAP" {
			continue
		}
		fmt.Printf(
			`
产品id: %s, 
上线时间: %s, 
最大杠杆倍数: %f, 
下单价格精度: %f, 
下单数量精度: %f, 
最小下单数量: %f, 
市价单最大委托张数: %f
每张面额：%f
Uly: %s
InstFamily: %s
`,
			instrument.InstID,
			instrument.ListTime.String(),
			instrument.Lever,
			instrument.TickSz,
			instrument.LotSz,
			instrument.MinSz,
			instrument.MaxMktSz,
			instrument.CtVal,
			instrument.Uly,
			instrument.InstFamily,
		)
	}

	return nil
}
