package voice

import (
	"errors"
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
	return errors.New("WIP")
}
