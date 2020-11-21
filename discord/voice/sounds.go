package voice

import (
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
