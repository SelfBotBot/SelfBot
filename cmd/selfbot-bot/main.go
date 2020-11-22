package main

import (
	"fmt"
	"os"
	"os/signal"
	"selfbot/config"
	"selfbot/discord"
	"selfbot/sound"
	"selfbot/sound/stores/filesystem"
	"selfbot/sound/stores/owo"
	"selfbot/web"
	"syscall"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/go-redis/redis/v8"

	"github.com/rs/zerolog"
)

func main() {
	var logger = zerolog.New(os.Stderr).With().Timestamp().Logger()

	var cfg config.Config
	if err := cfg.Load(); err != nil {
		panic(err)
	}

	var rdb *redis.Client
	if cfg.Redis.Enabled {
		rdb = redis.NewClient(&redis.Options{
			Addr:     cfg.Redis.Address,
			Password: cfg.Redis.Password,
			DB:       cfg.Redis.Database,
		})
	}

	var gDB *gorm.DB
	var soundStore sound.Store
	var err error
	if cfg.MySQL.Enabled {
		gDB, err = gorm.Open(mysql.Open(cfg.MySQL.URI), &gorm.Config{})
		if err != nil {
			panic(err)
		}

		soundStore, err = owo.NewStore(gDB, cfg)
		if err != nil {
			panic(err)
		}
	} else {
		soundStore, err = filesystem.New("./")
		if err != nil {
			panic(err)
		}
	}

	bot, err := discord.NewBot(logger, cfg.Discord, soundStore)
	if err != nil {
		panic(err)
	}

	webServer, err := web.NewServer(logger, rdb, cfg, bot.VoiceManager)
	if err != nil {
		panic(err)
	}

	if err := bot.Session.Open(); err != nil {
		panic(err)
	}

	go func() {
		if err := webServer.Start(cfg.Web.ListenAddress, true); err != nil {
			logger.Error().Err(err).Msg("Start web failed.")
		}
	}()

	fmt.Println("Waiting for interrupt.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	if err := bot.Close(); err != nil {
		panic(err)
	}

	if err := webServer.Close(); err != nil {
		panic(err)
	}
}
