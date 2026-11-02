package kraken

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

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
}

// publicClient is used for unauthenticated requests, with a timeout so a hung request can't stall the bot
var publicClient = &http.Client{Timeout: 15 * time.Second}

func NewKrakenApi() *KrakenApi {
	return &KrakenApi{
		api: rest.New(config.C.KrakenPublicKey, config.C.KrakenPrivateKey),
		log: log.WithFields(log.Fields{
			"prefix": "kraken",
		}),
	}
}

func (k *KrakenApi) GetCurrentBtcFiatPrice() (float64, error) {
	cryptoName := config.C.CryptoPrefix + "XBT" + config.C.FiatPrefix + config.C.Currency

	resp, err := publicClient.Get("https://api.kraken.com/0/public/Spread?pair=" + cryptoName)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("kraken spread request returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	return parseSpreadAskPrice(body, cryptoName)
}

// parseSpreadAskPrice returns the ask price of the latest entry in a kraken Spread response.
// Each entry is [time, bid, ask]. Any unexpected shape returns an error instead of panicking.
func parseSpreadAskPrice(body []byte, pair string) (float64, error) {
	var result KrakenSpread
	if err := json.Unmarshal(body, &result); err != nil {
		return 0, err
	}

	if len(result.Error) > 0 {
		return 0, fmt.Errorf("kraken returned an error for pair %s: %v", pair, result.Error)
	}

	prices, ok := result.Result[pair].([]interface{})
	if !ok || len(prices) == 0 {
		return 0, fmt.Errorf("no spread data for pair %s", pair)
	}

	latestPrices, ok := prices[len(prices)-1].([]interface{})
	if !ok || len(latestPrices) == 0 {
		return 0, fmt.Errorf("unexpected spread entry for pair %s: %v", pair, prices[len(prices)-1])
	}

	currentPrice, ok := latestPrices[len(latestPrices)-1].(string)
	if !ok {
		return 0, fmt.Errorf("unexpected price value for pair %s: %v", pair, latestPrices[len(latestPrices)-1])
	}

	parsedPrice, err := strconv.ParseFloat(currentPrice, 64)
	if err != nil {
		return 0, err
	}

	if parsedPrice <= 0 {
		return 0, fmt.Errorf("invalid price %v for pair %s", parsedPrice, pair)
	}

	return parsedPrice, nil
}
