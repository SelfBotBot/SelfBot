package voice

import (
	"errors"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Session struct {
	connection    *discordgo.VoiceConnection
	buffer        [][]byte
	bufferUpdated chan struct{}
	quit          chan struct{}
	speaking      bool
	vm            *Manager
}

func NewSession(s *discordgo.Session, vm *Manager, GuildID string, ChannelID string) (*Session, error) {
	if vm.stopping {
		return nil, errors.New("bot shutting down")
	}

	var err error
	ret := &Session{
		buffer:        make([][]byte, 0),
		bufferUpdated: make(chan struct{}),
		quit:          make(chan struct{}),
		vm:            vm,
	}

	ret.connection, err = s.ChannelVoiceJoin(GuildID, ChannelID, false, true)
	return ret, err
}

func (v *Session) StartLoop() {
	var data []byte

	tryReady := 0
	for !v.connection.Ready {
		time.Sleep(1 * time.Second)
		tryReady++
		if tryReady > 30 {
			v.Stop()
			return
		}
	}

	for {
		select {
		case <-v.quit:
			return
		default:
			if len(v.buffer) > 0 {
				v.setSpeaking(true)
				data, v.buffer = v.buffer[0], v.buffer[1:]
				v.connection.OpusSend <- data
			} else {
				v.setSpeaking(false)
				<-v.bufferUpdated
			}
		}
	}
}

func (v *Session) setSpeaking(s bool) {
	if s != v.speaking {
		if err := v.connection.Speaking(s); err != nil {
			v.Stop()
		}
		v.speaking = s
	}
}

func (v *Session) SetBuffer(data [][]byte) {
	if v.vm.stopping {
		return
	}
	isZero := len(v.buffer) == 0
	v.buffer = data
	if isZero && len(v.buffer) != 0 {
		v.bufferUpdated <- struct{}{}
	}
}

func (v *Session) Stop() error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Error recovered in Session.Stop() for "+v.connection.GuildID, r) // TODO logging.
		}
		if err := v.connection.Disconnect(); err != nil {
			fmt.Println("Unable to disconnect?!", err)
		}
		time.Sleep(50 * time.Millisecond)
		v.connection.Close()
	}()

	// Remove voice session from bot and stop the loop.
	delete(v.vm.SessionByGuild, v.connection.GuildID)
	delete(v.vm.SessionByChannel, v.connection.ChannelID)
	close(v.quit)

	if v.connection.Ready {
		//Broadcast "Goodbye".
		v.setSpeaking(true)
		for _, data := range Goodbye {
			v.connection.OpusSend <- data
		}
		v.setSpeaking(false)
	}

	return nil
}
