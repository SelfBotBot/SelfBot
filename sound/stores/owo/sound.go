package owo

import (
	"selfbot/sound"
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type StoredSound struct {
	gorm.Model
	ID     uuid.UUID `gorm:"primarykey"`
	Name   string    `gorm:"index:ux_storedsound_name_userid,unique"`
	UserID string    `gorm:"index:ux_storedsound_name_userid,unique"`
	OwoURL string
}

func (s *StoredSound) FromSound(sound sound.Sound) {
	s.ID = sound.ID
	s.UserID = sound.UserID
	s.Name = sound.Name
	s.CreatedAt = sound.CreatedAt
	if sound.Archived {
		if !s.DeletedAt.Valid || s.DeletedAt.Time.IsZero() {
			s.DeletedAt = gorm.DeletedAt{
				Time:  time.Now(),
				Valid: true,
			}
		}
	}

	return
}

func (s *StoredSound) ToSound() sound.Sound {
	var ret = sound.Sound{
		ID:        s.ID,
		Name:      s.Name,
		UserID:    s.UserID,
		CreatedAt: s.CreatedAt,
		Archived:  s.DeletedAt.Valid && !s.DeletedAt.Time.IsZero(),
	}

	return ret
}
