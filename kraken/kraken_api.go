package kraken

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/aopoltorzhicky/go_kraken/rest"
	"github.com/primexz/KrakenDCA/config"
	"github.com/primexz/KrakenDCA/logger"
	"github.com/primexz/KrakenDCA/notification"
)

var log *logger.Logger

func init() {
	log = logger.NewLogger("kraken_api")
}

type KrakenSpread struct {
	Error  []interface{}          `json:"error"`
	Result map[string]interface{} `json:"result"`
	Last   int                    `json:"last"`
}

func GetCurrentBtcFiatPrice() (float64, error) {
	cryptoName := config.CryptoPrefix + "XBT" + config.FiatPrefix + config.Currency

	resp, err := http.Get("https://api.kraken.com/0/public/Spread?pair=" + cryptoName)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var result KrakenSpread
	if err := json.Unmarshal(body, &result); err != nil {
		return 0, err
	}

	prices := result.Result[cryptoName].([]interface{})
	latestPrices := prices[len(prices)-1].([]interface{})
	currentPrice := latestPrices[len(latestPrices)-1]

	parsedPrice, err := strconv.ParseFloat(currentPrice.(string), 32)
	if err != nil {
		return 0, err
	}

	return parsedPrice, nil
}

func BuyBtc(_retry_do_not_use int) {
	if _retry_do_not_use > 3 {
		log.Error("Failed to buy btc after 3 retries, stop recursion")
		return
	}

	currency := config.Currency

	fiatPrice, err := GetCurrentBtcFiatPrice()
	if err != nil {
		log.Error("Failed to get current btc price", err.Error())
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

		log.Debug("Buying btc with price ", args["price"])

		response, krakenErr = getApi().AddOrder("xbt"+strings.ToLower(currency), "buy", "limit", config.KrakenOrderSize, args)
		if krakenErr != nil {
			log.Error("Failed to buy btc", krakenErr.Error())
			return
		}

		transactionId := response.TransactionIds[0]

		for {
			orderInfo, err := getApi().QueryOrders(true, "", transactionId)
			if err != nil {
				log.Error("Failed to get order status", err.Error())
				return
			}

			order, ok := orderInfo[transactionId]
			if !ok {
				log.Error("Failed to query order status")
				return
			}

			orderStatus := order.Status
			log.Info("current order status:", orderStatus)

			if orderStatus == "closed" {
				log.Info("Order successfully executed")
				break // don't return to send notification and log
			}

			if orderStatus == "canceled" && order.Reason == "User requested" {
				log.Info("Order canceled by user")
				return
			}

			if orderStatus == "canceled" {
				log.Info("Unknown reason for order cancelation.")
				return
			}

			if orderStatus == "canceled" && order.Reason == "Post only order" {
				log.Info("Order canceled by kraken due to post only order, retrying with new order")
				BuyBtc(_retry_do_not_use + 1)
				return
			}

			if orderStatus == "expired" {
				log.Info("Order expired, retrying with new order")
				BuyBtc(_retry_do_not_use + 1)
				return
			}

			//wait on pending, open
			time.Sleep(5 * time.Second)
		}

	} else {
		response, krakenErr = getApi().AddOrder("xbt"+strings.ToLower(currency), "buy", "market", config.KrakenOrderSize, nil)
		if krakenErr != nil {
			log.Error("Failed to buy btc", krakenErr.Error())
			return
		}
	}

	notification.SendPushNotification("BTC bought", fmt.Sprintf("Description: %s\nPrice: %f %s", response.Description.Info, fiatPrice, currency))

	log.Info("Successfully bought btc ->", response.Description.Info, response.Description.Price)
}

func getApi() *rest.Kraken {
	return rest.New(config.KrakenPublicKey, config.KrakenPrivateKey)
}
