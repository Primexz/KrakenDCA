package main

import (
	"fmt"
	"runtime"

	"github.com/primexz/KrakenDCA/bot"
	"github.com/primexz/KrakenDCA/config"
	"github.com/primexz/KrakenDCA/logger"
)

var (
	log *logger.Logger
)

func init() {
	log = logger.NewLogger("main")
}

func main() {
	log.Info(fmt.Sprintf("Kraken DCA 🐙 %s, commit %s, built at %s (%s [%s, %s])", version, commit, date, runtime.Version(), runtime.GOOS, runtime.GOARCH))

	config.LoadConfiguration()

	if config.ExperimentalMakerFee {
		log.Warn("Experimental maker fee is enabled. This feature is not recommended for production use.")
	}

	bot.NewBot().StartBot()

	select {}
}
