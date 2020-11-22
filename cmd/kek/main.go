package main

import (
	"fmt"
	"selfbot/config"
	"selfbot/sound"
	"selfbot/sound/stores/filesystem"
	"selfbot/sound/stores/owo"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	var cfg config.Config
	if err := cfg.Load(); err != nil {
		fmt.Println(err)
	}

	var gDB *gorm.DB
	var err error
	gDB, err = gorm.Open(mysql.Open(cfg.MySQL.URI), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	owoStore, err := owo.NewStore(gDB, cfg)
	if err != nil {
		panic(err)
	}
	fsStore, err := filesystem.New("./")
	if err != nil {
		panic(err)
	}

	resp, err := fsStore.ListSounds(sound.ListOptions{})
	if err != nil {
		panic(err)
	}

	for _, v := range resp.SoundIDs {
		s, err := fsStore.LoadSound(v)
		if err != nil {
			fmt.Println(err)
			continue
		}

		s.UserID = "416717866411360258"

		if _, err := owoStore.SaveSound(&s); err != nil {
			fmt.Println("Saving owo sound: ", err)
		}

		time.Sleep(time.Second * 10)
	}

}
