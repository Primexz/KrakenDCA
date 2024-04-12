package kraken

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/aopoltorzhicky/go_kraken/rest"
	"github.com/primexz/KrakenDCA/config"
	"github.com/primexz/KrakenDCA/notification"
)

type KrakenSpread struct {
	Error  []interface{}          `json:"error"`
	Result map[string]interface{} `json:"result"`
	Last   int                    `json:"last"`
}

func GetCurrentBtcFiatPrice() float64 {
	cryptoName := config.CryptoPrefix + "XBT" + config.FiatPrefix + config.Currency

	resp, err := http.Get("https://api.kraken.com/0/public/Spread?pair=" + cryptoName)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var result KrakenSpread
	if err := json.Unmarshal(body, &result); err != nil {
		log.Fatal(result)
	}

	prices := result.Result[cryptoName].([]interface{})
	latestPrices := prices[len(prices)-1].([]interface{})
	currentPrice := latestPrices[len(latestPrices)-1]

	parsedPrice, err := strconv.ParseFloat(currentPrice.(string), 32)
	if err != nil {
		log.Fatal("failed to parse price", parsedPrice)
	}

	return parsedPrice
}

func BuyBtc() {
	currency := config.Currency

	response, err := getApi().AddOrder("xbt"+strings.ToLower(currency), "buy", "market", config.KrakenOrderSize, nil)
	if err != nil {
		log.Println("Failed to buy btc", err.Error())
		return
	}

	notification.SendPushNotification("BTC bought", fmt.Sprintf("Description: %s\nPrice: %f %s", response.Description.Info, GetCurrentBtcFiatPrice(), currency))

	log.Println("Successfully bought btc ->", response.Description.Info, response.Description.Price)
}

func getApi() *rest.Kraken {
	return rest.New(config.KrakenPublicKey, config.KrakenPrivateKey)
}
