package kraken

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/aopoltorzhicky/go_kraken/rest"
	"github.com/primexz/KrakenDCA/config"
	log "github.com/sirupsen/logrus"
)

type KrakenApi struct {
	api *rest.Kraken

	log *log.Entry
}

type KrakenSpread struct {
	Error  []interface{}          `json:"error"`
	Result map[string]interface{} `json:"result"`
	Last   int                    `json:"last"`
}

func NewKrakenApi() *KrakenApi {
	return &KrakenApi{
		api: rest.New(config.KrakenPublicKey, config.KrakenPrivateKey),
		log: log.WithFields(log.Fields{
			"prefix": "kraken",
		}),
	}
}

func (k *KrakenApi) GetCurrentBtcFiatPrice() (float64, error) {
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
