package main

import (
	"fmt"
	"os"
	"os/signal"
	"selfbot/config"
	"selfbot/discord"
	"syscall"

	"github.com/rs/zerolog"
)

func main() {
	var logger = zerolog.New(os.Stderr).With().Timestamp().Logger()

	var cfg = config.Config{}
	if err := cfg.Load(); err != nil {
		panic(err)
	}

	bot, err := discord.NewBot(logger, cfg.Discord)
	if err != nil {
		panic(err)
	}

	if err := bot.Session.Open(); err != nil {
		panic(err)
	}

	fmt.Println("Waiting for interrupt.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	if err := bot.Close(); err != nil {
		panic(err)
	}
}
