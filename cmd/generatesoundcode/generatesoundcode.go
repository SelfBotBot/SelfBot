package main

import (
	"fmt"
	"selfbot/discord/voice"
)

func main() {
	data, err := voice.LoadSound("welcome.dca")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v\n", data)
}
