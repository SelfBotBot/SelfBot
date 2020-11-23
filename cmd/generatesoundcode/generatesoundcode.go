package main

import (
	"fmt"
	"os"
	"selfbot/sound"
)

func main() {
	data, err := readSoundFile("welcome.dca")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v\n", data)
}

func readSoundFile(fileName string) ([][]byte, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("read sound file: open file: %w", err)
	}

	defer file.Close()
	ret, err := sound.DataRead(file)
	if err != nil {
		return nil, fmt.Errorf("read sound file: %w", err)
	}

	return ret, nil
}
