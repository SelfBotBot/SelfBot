package filesystem

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"selfbot/sound"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
)

var _ sound.SoundStore = new(Store)

type Store struct {
	soundsFolder string
	soundNames   map[uuid.UUID]string
}

func New(soundsFolder string) (*Store, error) {
	ret := &Store{
		soundsFolder: soundsFolder,
		soundNames:   make(map[uuid.UUID]string),
	}

	files, err := ioutil.ReadDir(ret.soundsFolder)
	if err != nil {
		return nil, fmt.Errorf("new filesystem sound store: readdir: %w", err)
	}

	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".dca") {
			ret.soundNames[uuid.Must(uuid.NewV4())] = f.Name()[0 : len(f.Name())-4]
		}
	}

	return ret, nil
}

func (s *Store) SaveSound(sound *sound.Sound) (soundID uuid.UUID, err error) {
	return uuid.Nil, errors.New("unimplemented: saving sounds for filesystem isn't supported")
}

func (s *Store) LoadSound(soundID uuid.UUID) (soundInfo sound.Sound, err error) {
	soundData, err := readSoundFile(s.soundNames[soundID] + ".dca")
	if err != nil {
		return sound.Sound{}, fmt.Errorf("load sound: %w", err)
	}

	return sound.Sound{
		ID:        soundID,
		Name:      s.soundNames[soundID],
		UserID:    "416717866411360258",
		CreatedAt: time.Now(),
		Data:      soundData,
	}, nil
}

func (s *Store) ListSounds(listOptions sound.ListOptions) (listResponse sound.ListResponse, err error) {
	var keys []uuid.UUID
	for k := range s.soundNames {
		keys = append(keys, k)
	}
	return sound.ListResponse{
		SoundIDs: keys,
	}, nil
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
