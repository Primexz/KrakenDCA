package main

import (
	"log"
	"runtime"

	"github.com/primexz/KrakenDCA/bot"
	"github.com/primexz/KrakenDCA/config"
)

func main() {
	log.Printf("Kraken DCA 🐙 %s, commit %s, built at %s (%s [%s, %s])", version, commit, date, runtime.Version(), runtime.GOOS, runtime.GOARCH)

	config.LoadConfiguration()
	bot.StartBot()
}
