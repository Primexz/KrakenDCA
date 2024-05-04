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
	bot.NewBot().StartBot()

	select {}
}
