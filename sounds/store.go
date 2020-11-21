package sounds

import uuid "github.com/satori/go.uuid"

type SoundStore interface {
	SaveSound(sound *Sound) (soundID uuid.UUID, err error)
	LoadSound(soundID uuid.UUID) (sound Sound, err error)
	ListSounds(listOptions ListOptions) (listResponse ListResponse, err error)
}

type ListOptions struct {
	UserID    string    `json:"user_id"`
	SoundName string    `json:"sound_name"`
	AfterID   uuid.UUID `json:"after_id"`
	Limit     int       `json:"limit"`
}

type ListResponse struct {
	SoundIDs []uuid.UUID
}
