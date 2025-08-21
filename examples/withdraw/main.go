package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	i_logger "github.com/pefish/go-interface/i-logger"
	okex "github.com/pefish/go-okx"
	"github.com/pefish/go-okx/api"
	"github.com/pefish/go-okx/requests/rest/funding"
	go_random "github.com/pefish/go-random"
	"github.com/pkg/errors"
)

var toAddress = "5BnsHy3CV2SjefwMPQ4pwQPVmigxA8R7gUZypRNsZqxp:"
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

	// 29bdc6c7-e3e6-4694-a2a2-08f6cf73939e
	withdrawID := go_random.MustRandomStringFromDic("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz", 32)
	withdrawalRes, err := client.Rest.Funding.Withdrawal(funding.Withdrawal{
		Ccy:      "SOL",
		Chain:    "SOL-Solana",
		ToAddr:   toAddress,
		Amt:      amount,
		Type:     okex.WithdrawalDigitalAddressDestination,
		ClientID: withdrawID,
	})
	if err != nil {
		return err
	}

	if withdrawalRes.Code != 0 {
		return errors.Errorf("withdrawal failed. err: %s, code: %d", withdrawalRes.Msg, withdrawalRes.Code)
	}

	fmt.Printf("<WdID: %d>\n", int64(withdrawalRes.Withdrawals[0].WdID))

	timer := time.NewTimer(0)
watchWithdraw:
	for {
		select {
		case <-timer.C:
			getWithdrawalHistoryRes, err := client.Rest.Funding.GetWithdrawalHistory(funding.GetWithdrawalHistory{
				ClientID: withdrawID,
			})
			if err != nil {
				return err
			}
			if getWithdrawalHistoryRes.Code != 0 {
				return errors.Errorf("getWithdrawalHistory failed. err: %s, code: %d", getWithdrawalHistoryRes.Msg, getWithdrawalHistoryRes.Code)
			}
			switch getWithdrawalHistoryRes.WithdrawalHistories[0].State {
			case -2:
				return errors.New("提现 Canceled")
			case -1:
				return errors.New("提现 Failed")
			case 2:
				break watchWithdraw
			default:
				fmt.Printf("继续监听提币进程...")
				timer.Reset(3 * time.Second)
				continue
			}

		}
	}
	fmt.Printf("提币成功\n")

	return nil
}
