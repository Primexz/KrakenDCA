package kraken

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aopoltorzhicky/go_kraken/rest"
	"github.com/primexz/KrakenDCA/config"
	"github.com/primexz/KrakenDCA/notification"
)

func (k *KrakenApi) BuyBtc(_retry_do_not_use int) {
	if _retry_do_not_use > 5 {
		k.log.Error("Failed to buy btc after 5 retries, stop recursion")
		return
	}

	currency := config.Currency

	fiatPrice, err := k.GetCurrentBtcFiatPrice()
	if err != nil {
		k.log.Error("Failed to get current btc price", err.Error())
	}

	var (
		response  rest.AddOrderResponse
		krakenErr error
	)

	if config.ExperimentalMakerFee {
		priceRound, _ := strconv.ParseFloat(fmt.Sprintf("%.1f", fiatPrice), 64)
		args := map[string]interface{}{
			// if set to true, no order will be submitted
			"validate": false,

			//price can only be specified up to 1 decimals
			"price":       priceRound - 0.1,
			"oflags":      "post",
			"timeinforce": "GTD",
			"expiretm":    "+240", // close order after 4 minutes
		}

		response, krakenErr = k.api.AddOrder("xbt"+strings.ToLower(currency), "buy", "limit", config.KrakenOrderSize, args)
		if krakenErr != nil {
			k.log.Error("Failed to buy btc", krakenErr.Error())
			return
		}

		transactionId := response.TransactionIds[0]

		for {
			orderInfo, err := k.api.QueryOrders(true, "", transactionId)
			if err != nil {
				k.log.Error("Failed to get order status", err.Error())
				return
			}

			order, ok := orderInfo[transactionId]
			if ok {
				orderStatus := order.Status
				k.log.Info("Order status:", orderStatus)

				if orderStatus == "closed" {
					k.log.Info("Order successfully executed")
					break // don't return to send notification and log
				}

				if orderStatus == "canceled" && order.Reason == "User requested" {
					k.log.Info("Order canceled by user")
					return
				}

				if orderStatus == "canceled" && order.Reason == "Post only order" {
					k.log.Info("Order canceled by kraken due to post only order, retrying with new order")
					k.BuyBtc(_retry_do_not_use + 1)
					return
				}

				if orderStatus == "canceled" {
					k.log.Info("Unknown reason for order cancelation.")
					return
				}

				if orderStatus == "expired" {
					k.log.Info("Order expired, retrying with new order")
					k.BuyBtc(_retry_do_not_use + 1)
					return
				}
			} else {
				k.log.Error("Failed to query order status")
			}

			//wait on pending, open
			time.Sleep(5 * time.Second)
		}

	} else {
		response, krakenErr = k.api.AddOrder("xbt"+strings.ToLower(currency), "buy", "market", config.KrakenOrderSize, nil)
		if krakenErr != nil {
			k.log.Error("Failed to buy btc", krakenErr.Error())
			return
		}
	}

	notification.SendPushNotification("BTC bought", fmt.Sprintf("Description: %s\nPrice: %f %s", response.Description.Info, fiatPrice, currency))

	k.log.Info("Successfully bought btc ->", response.Description.Info, response.Description.Price)
}
