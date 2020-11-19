package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"selfbot/config"
	"selfbot/discord"
	"strings"
	"syscall"
)

func main() {

	var cfg = config.Config{}
	if err := cfg.Load(); err != nil {
		panic(err)
	}

	bot, err := discord.NewBot(cfg.Discord)
	if err != nil {
		panic(err)
	}

	files, err := ioutil.ReadDir("./")
	if err != nil {
		panic(err)
	}

	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".dca") {
			if err := bot.VoiceManager.LoadSound(f.Name(), f.Name()[0:len(f.Name())-4]); err != nil {
				fmt.Println("Unable to load "+f.Name()+", ", err)
				continue
			}
		}
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
