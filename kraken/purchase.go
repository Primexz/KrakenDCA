package kraken

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/primexz/KrakenDCA/config"
	"github.com/primexz/KrakenDCA/notification"
)

func (k *KrakenApi) BuyBtc() {
	if config.LimitOrderMode {
		for i := 0; i < config.LimitOrderRetryCount; i++ {
			fiatPrice, err := k.GetCurrentBtcFiatPrice()
			if err != nil {
				k.log.Error("Failed to get current btc price", err)
			}

			if k.placeLimitOrder(fiatPrice) {
				k.log.Info("Successfully bought limit btc")
				return
			}

			k.log.Warn("Retrying to place limit order")
		}

		notification.SendPushNotification("Failed to buy btc", fmt.Sprintf("Failed to buy btc after %d retries", config.LimitOrderRetryCount))
		k.log.Errorf("Failed to buy btc after %d retries", config.LimitOrderRetryCount)
	} else {
		if !k.placeMarketOrder() {
			notification.SendPushNotification("Failed to buy btc", "Failed to buy btc")
			k.log.Error("Failed to buy btc")
		}
	}
}

// placeMarketOrder places a market order on kraken
func (k *KrakenApi) placeMarketOrder() bool {
	if _, err := k.api.AddOrder("xbt"+strings.ToLower(config.Currency), "buy", "market", config.KrakenOrderSize, nil); err != nil {
		k.log.Error("Failed to buy btc", err.Error())
		return false
	}

	return true
}

// placeLimitOrder places a limit order on kraken
// returns true if order was successfully placed
// returns false if order was not placed and should be retried
func (k *KrakenApi) placeLimitOrder(fiatPrice float64) bool {
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

	response, err := k.api.AddOrder("xbt"+strings.ToLower(config.Currency), "buy", "limit", config.KrakenOrderSize, args)
	if err != nil {
		k.log.Error("Failed to buy btc", err)
		return false
	}

	transactionId := response.TransactionIds[0]

	for {
		orderInfo, err := k.api.QueryOrders(true, "", transactionId)
		if err != nil {
			k.log.Error("Failed to get order status", err.Error())
			return false
		}

		order, ok := orderInfo[transactionId]
		if ok {
			orderStatus := order.Status
			k.log.Info("Order status:", orderStatus)

			if orderStatus == "closed" {
				k.log.Info("Order successfully executed")
				break
			}

			if orderStatus == "canceled" && order.Reason == "User requested" {
				k.log.Info("Order canceled by user")
				return true
			}

			if orderStatus == "canceled" && order.Reason == "Post only order" {
				k.log.Info("Order canceled by kraken due to post only order, retrying with new order")
				return false
			}

			if orderStatus == "canceled" {
				k.log.Info("Unknown reason for order cancelation.")
				return true
			}

			if orderStatus == "expired" {
				k.log.Info("Order expired, retrying with new order")
				return false
			}
		} else {
			k.log.Error("Failed to query order status")
		}

		//wait on pending, open
		time.Sleep(5 * time.Second)
	}

	return true
}
