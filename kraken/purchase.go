package kraken

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/primexz/KrakenDCA/config"
	"github.com/primexz/KrakenDCA/notification"
	log "github.com/sirupsen/logrus"
)

var (
	assetPair = "xbt" + strings.ToLower(config.C.Currency)
)

const (
	// how often a failing order status query is retried before we try to cancel the order
	maxOrderQueryErrors = 12
)

func (k *KrakenApi) BuyBtc() {
	if config.C.LimitOrderMode {
		limitAdjustment := 0.1

		for i := 0; i < config.C.LimitOrderRetryCount; i++ {
			fiatPrice, err := k.GetCurrentBtcFiatPrice()
			if err != nil {
				k.log.Error("Failed to get current btc price", err)
				time.Sleep(5 * time.Second)
				continue
			}

			success, action := k.placeLimitOrder(fiatPrice, limitAdjustment)
			if success {
				return
			}

			switch action {
			case INCREASE_LIMIT_ADJUSTMENT:
				limitAdjustment += 0.1
				k.log.WithField("limitAdjustment", limitAdjustment).Warn("Increasing limit adjustment")
			case STOP:
				notification.SendPushNotification("Failed to buy btc", "Stopped buying btc, check your open orders on kraken")
				k.log.Error("Stopped buying btc, check your open orders on kraken")
				return
			}

			k.log.Warn("Retrying to place limit order")
		}

		if config.C.LimitOrderMarketFallback {
			k.log.Warnf("Post-only order not filled after %d retries, falling back to a market order (taker fee)", config.C.LimitOrderRetryCount)
			if k.placeMarketOrder() {
				return
			}
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

// ValidatePostOnlyOrder asks kraken to validate a post-only limit order without submitting it.
// This checks the API key permissions, the pair, the order size and the price format.
func (k *KrakenApi) ValidatePostOnlyOrder() error {
	fiatPrice, err := k.GetCurrentBtcFiatPrice()
	if err != nil {
		return fmt.Errorf("failed to get current btc price: %w", err)
	}

	args := limitOrderArgs(fiatPrice, 0.1)
	args["validate"] = true

	response, err := k.api.AddOrder(assetPair, "buy", "limit", config.C.KrakenOrderSize, args)
	if err != nil {
		return err
	}

	k.log.WithField("order", response.Description.Info).Info("Post-only order validated by kraken")
	return nil
}

// placeMarketOrder places a market order on kraken
func (k *KrakenApi) placeMarketOrder() bool {
	if _, err := k.api.AddOrder(assetPair, "buy", "market", config.C.KrakenOrderSize, nil); err != nil {
		k.log.Error("Failed to buy btc", err.Error())
		return false
	}

	return true
}

// limitOrderArgs builds the arguments for a post-only limit buy order just below the ask price
func limitOrderArgs(askPrice float64, limitAdjustment float64) map[string]interface{} {
	return map[string]interface{}{
		// if set to true, no order will be submitted
		"validate": false,

		"price":       formatLimitPrice(askPrice, limitAdjustment),
		"oflags":      "post",
		"timeinforce": "GTD",
		"expiretm":    "+240", // close order after 4 minutes
	}
}

// formatLimitPrice returns ask price minus the adjustment with exactly 1 decimal,
// as kraken rejects prices with more decimals than the pair allows
func formatLimitPrice(askPrice float64, limitAdjustment float64) string {
	tenths := math.Round(askPrice*10) - math.Round(limitAdjustment*10)
	return strconv.FormatFloat(tenths/10, 'f', 1, 64)
}

// placeLimitOrder places a post-only limit order on kraken
// returns true if no further order should be placed
// returns false if order was not filled and should be retried
func (k *KrakenApi) placeLimitOrder(fiatPrice float64, limitAdjustment float64) (bool, LimitPurchaseAction) {
	response, err := k.api.AddOrder(assetPair, "buy", "limit", config.C.KrakenOrderSize, limitOrderArgs(fiatPrice, limitAdjustment))
	if err != nil {
		k.log.Error("Failed to buy btc", err)
		return false, NONE
	}

	transactionId := response.TransactionIds[0]
	queryErrors := 0

	for {
		//wait on pending, open
		time.Sleep(5 * time.Second)

		orderInfo, err := k.api.QueryOrders(true, "", transactionId)
		order, ok := orderInfo[transactionId]
		if err != nil || !ok {
			queryErrors++
			k.log.WithField("attempt", queryErrors).Error("Failed to get order status", err)

			// the order may still be open, never place a new one before it is gone
			if queryErrors >= maxOrderQueryErrors {
				return k.cancelUnknownOrder(transactionId)
			}
			continue
		}
		queryErrors = 0

		orderStatus := order.Status
		k.log.Info("Order status:", orderStatus)

		switch {
		case orderStatus == "closed":
			k.logFilledOrder(order.Cost, order.Fee, order.VolumeExecuted)
			return true, NONE

		case order.VolumeExecuted > 0 && (orderStatus == "canceled" || orderStatus == "expired"):
			// the rest is below kraken's minimum order size, so it can't be bought with a new order
			k.log.Warn("Order was only partially filled, not placing a new order")
			k.logFilledOrder(order.Cost, order.Fee, order.VolumeExecuted)
			return true, NONE

		case orderStatus == "canceled" && order.Reason == "User requested":
			k.log.Info("Order canceled by user")
			return true, NONE

		case orderStatus == "canceled" && order.Reason == "Post only order":
			k.log.Info("Order canceled by kraken due to post only order, retrying with new order")

			// This happens when the price drops too fast, return the action type to increase the limit adjustment
			return false, INCREASE_LIMIT_ADJUSTMENT

		case orderStatus == "canceled":
			k.log.WithField("reason", order.Reason).Error("Unknown reason for order cancelation.")
			return false, STOP

		case orderStatus == "expired":
			k.log.Info("Order expired, retrying with new order")
			return false, NONE
		}
	}
}

// cancelUnknownOrder cancels an order whose status could not be queried.
// If the cancel fails too, the order may still fill, so no new order is placed.
func (k *KrakenApi) cancelUnknownOrder(transactionId string) (bool, LimitPurchaseAction) {
	if _, err := k.api.Cancel(transactionId); err != nil {
		k.log.WithField("txid", transactionId).Error("Failed to cancel order with unknown status", err)
		return false, STOP
	}

	k.log.Warn("Canceled order with unknown status, retrying with new order")
	return false, NONE
}

func (k *KrakenApi) logFilledOrder(cost float64, fee float64, volume float64) {
	feePercent := 0.0
	if cost > 0 {
		feePercent = fee / cost * 100
	}

	k.log.WithFields(log.Fields{
		"volume":     volume,
		"cost":       cost,
		"fee":        fee,
		"feePercent": fmt.Sprintf("%.2f%%", feePercent),
	}).Info("Successfully bought limit btc")
}
