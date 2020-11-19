package voice

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"sort"
)

func (m *Manager) ListSounds() []string {
	keys := make([]string, 0, len(m.Sounds))
	for k := range m.Sounds {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	return keys
}

// loadSound attempts to load an encoded sound file from disk.
func (m *Manager) LoadSound(fileName, name string) error {
	data, err := LoadSound(fileName)
	if err != nil {
		return err
	}

	m.Sounds[name] = data
	return nil
}

func LoadSound(fileName string) ([][]byte, error) {

	var ret [][]byte
	file, err := os.Open(fileName)
	if err != nil {
		return ret, err
	}

	var opuslen int16
	for {
		// Read opus frame length from dca file.
		err = binary.Read(file, binary.LittleEndian, &opuslen)

		// If this is the end of the file, just return.
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			err := file.Close()
			if err != nil {
				return ret, err
			}
			return ret, nil
		}

		if err != nil {
			fmt.Println("Error reading from dca file :", err)
			return ret, err
		}

		// Read encoded pcm from dca file.
		InBuf := make([]byte, opuslen)
		err = binary.Read(file, binary.LittleEndian, &InBuf)

		// Should not be any end of file errors
		if err != nil {
			fmt.Println("Error reading from dca file :", err)
			return ret, err
		}

		ret = append(ret, InBuf)
	}
}
