package voice

import (
	"sort"
)

func (m *Manager) ListSounds() []string {
	keys := make([]string, 0, len(m.Sounds))
	for _, v := range m.Sounds {
		keys = append(keys, v.Name)
	}
	sort.Strings(keys)

	return keys
}
