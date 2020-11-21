package owo

import (
	"net/http"
	"selfbot/discord/voice"
	"selfbot/sound"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
)

func TestBlah(t *testing.T) {
	owoStore := Store{
		client: Client{
			UploadURL: "",
			Client:    &http.Client{Timeout: time.Second * 10},
		},
	}

	owoStore.SaveSound(&sound.Sound{
		ID:        uuid.Must(uuid.NewV4()),
		Name:      "blah",
		Data:      voice.Goodbye,
		UserID:    "blah",
		CreatedAt: time.Now(),
		Archived:  false,
	})
}
