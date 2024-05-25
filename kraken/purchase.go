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
	if config.C.LimitOrderMode {
		limitAdjustment := 0.1

		for i := 0; i < config.C.LimitOrderRetryCount; i++ {
			fiatPrice, err := k.GetCurrentBtcFiatPrice()
			if err != nil {
				k.log.Error("Failed to get current btc price", err)
			}

			success, action := k.placeLimitOrder(fiatPrice, limitAdjustment)
			if success {
				k.log.Info("Successfully bought limit btc")
				return
			}

			if action == INCREASE_LIMIT_ADJUSTMENT {
				limitAdjustment += 0.1
				k.log.WithField("limitAdjustment", limitAdjustment).Warn("Increasing limit adjustment")
			}

			k.log.Warn("Retrying to place limit order")
		}

		notification.SendPushNotification("Failed to buy btc", fmt.Sprintf("Failed to buy btc after %d retries", config.C.LimitOrderRetryCount))
		k.log.Errorf("Failed to buy btc after %d retries", config.C.LimitOrderRetryCount)
	} else {
		if !k.placeMarketOrder() {
			notification.SendPushNotification("Failed to buy btc", "Failed to buy btc")
			k.log.Error("Failed to buy btc")
		}
	}
}

// placeMarketOrder places a market order on kraken
func (k *KrakenApi) placeMarketOrder() bool {
	if _, err := k.api.AddOrder("xbt"+strings.ToLower(config.C.Currency), "buy", "market", config.C.KrakenOrderSize, nil); err != nil {
		k.log.Error("Failed to buy btc", err.Error())
		return false
	}

	return true
}

// placeLimitOrder places a limit order on kraken
// returns true if order was successfully placed
// returns false if order was not placed and should be retried
func (k *KrakenApi) placeLimitOrder(fiatPrice float64, limitAdjustment float64) (bool, LimitPurchaseAction) {
	priceRound, _ := strconv.ParseFloat(fmt.Sprintf("%.1f", fiatPrice), 64)

	args := map[string]interface{}{
		// if set to true, no order will be submitted
		"validate": false,

		//price can only be specified up to 1 decimals
		"price":       priceRound - limitAdjustment,
		"oflags":      "post",
		"timeinforce": "GTD",
		"expiretm":    "+240", // close order after 4 minutes
	}

	response, err := k.api.AddOrder("xbt"+strings.ToLower(config.C.Currency), "buy", "limit", config.C.KrakenOrderSize, args)
	if err != nil {
		k.log.Error("Failed to buy btc", err)
		return false, NONE
	}

	transactionId := response.TransactionIds[0]

	for {
		orderInfo, err := k.api.QueryOrders(true, "", transactionId)
		if err != nil {
			k.log.Error("Failed to get order status", err.Error())
			return false, NONE
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
				return true, NONE
			}

			if orderStatus == "canceled" && order.Reason == "Post only order" {
				k.log.Info("Order canceled by kraken due to post only order, retrying with new order")

				// This happens when the price drops too fast, return the action type to increase the limit adjustment
				return false, INCREASE_LIMIT_ADJUSTMENT
			}

			if orderStatus == "canceled" {
				k.log.Info("Unknown reason for order cancelation.")
				return true, NONE
			}

			if orderStatus == "expired" {
				k.log.Info("Order expired, retrying with new order")
				return false, NONE
			}
		} else {
			k.log.Error("Failed to query order status")
		}

		//wait on pending, open
		time.Sleep(5 * time.Second)
	}

	return true, NONE
}
