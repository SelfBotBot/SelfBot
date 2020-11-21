package filesystem

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"selfbot/sounds"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
)

var _ sounds.SoundStore = new(Store)

type Store struct {
	soundsFolder string
	soundNames   map[uuid.UUID]string
}

func New(soundsFolder string) (*Store, error) {
	ret := &Store{
		soundsFolder: soundsFolder,
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

func (s *Store) SaveSound(sound *sounds.Sound) (soundID uuid.UUID, err error) {
	return uuid.Nil, errors.New("unimplemented: saving sounds for filesystem isn't supported")
}

func (s *Store) LoadSound(soundID uuid.UUID) (sound sounds.Sound, err error) {
	soundData, err := readSoundFile(s.soundNames[soundID] + ".dca")
	if err != nil {
		return sounds.Sound{}, fmt.Errorf("load sound: %w", err)
	}

	return sounds.Sound{
		ID:        soundID,
		CreatedAt: time.Now(),
		Data:      soundData,
	}, nil
}

func (s *Store) ListSounds(listOptions sounds.ListOptions) (listResponse sounds.ListResponse, err error) {
	var keys []uuid.UUID
	for k := range s.soundNames {
		keys = append(keys, k)
	}
	return sounds.ListResponse{
		SoundIDs: keys,
	}, nil
}

func readSoundFile(fileName string) ([][]byte, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("read sound file: open file: %w", err)
	}

	defer file.Close()
	ret, err := sounds.SoundDataRead(file)
	if ret != nil {
		return nil, fmt.Errorf("read sound file: %w", err)
	}

	return ret, nil
}
