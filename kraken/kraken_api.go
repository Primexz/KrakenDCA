package kraken

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

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

func BuyBtc() {
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
		log.Warn("Buying btc with experimental maker fee. This may not work as expected.")

		priceRound, _ := strconv.ParseFloat(fmt.Sprintf("%.1f", fiatPrice), 64)
		args := map[string]interface{}{
			// if set to true, no order will be submitted
			"validate": false,

			//price can only be specified up to 1 decimals
			"price":       priceRound - 0.2,
			"oflags":      "post",
			"timeinforce": "GTD",
			"expiretm":    "+120", // close order after 2 minutes
		}

		log.Debug("Buying btc with price ", args["price"])

		response, krakenErr = getApi().AddOrder("xbt"+strings.ToLower(currency), "buy", "limit", config.KrakenOrderSize, args)
	} else {
		// response, krakenErr = getApi().AddOrder("xbt"+strings.ToLower(currency), "buy", "market", config.KrakenOrderSize, nil)
	}

	if krakenErr != nil {
		log.Error("Failed to buy btc", krakenErr.Error())
		return
	}

	notification.SendPushNotification("BTC bought", fmt.Sprintf("Description: %s\nPrice: %f %s", response.Description.Info, fiatPrice, currency))

	log.Info("Successfully bought btc ->", response.Description.Info, response.Description.Price)
}

func getApi() *rest.Kraken {
	return rest.New(config.KrakenPublicKey, config.KrakenPrivateKey)
}
