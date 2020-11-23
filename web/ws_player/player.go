package ws_player

import (
	"log"
	"selfbot/discord/voice"
	"time"

	uuid "github.com/satori/go.uuid"

	"github.com/gorilla/websocket"

	"github.com/markbates/goth"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type PlayRequest struct {
	ID uuid.UUID `json:"id"`
}

type VoicePlayer struct {
	GuildID string
	User    goth.User
	vm      *voice.Manager
	conn    *websocket.Conn
}

func NewVoicePlayer(guildID string, user goth.User, vm *voice.Manager) *VoicePlayer {
	ret := &VoicePlayer{
		GuildID: guildID,
		User:    user,
		vm:      vm,
	}
	return ret
}

func (p *VoicePlayer) Start(wsConn *websocket.Conn) {
	p.conn = wsConn
	go p.reader()
	go p.writer()
}

func (p *VoicePlayer) reader() {
	defer func() {
		p.conn.Close()
	}()
	p.conn.SetReadLimit(maxMessageSize)
	p.conn.SetReadDeadline(time.Now().Add(pongWait))
	p.conn.SetPongHandler(func(string) error { p.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	var playReq PlayRequest
	for {
		err := p.conn.ReadJSON(&playReq)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		go p.vm.Play(p.GuildID, playReq.ID)
	}
}

func (p *VoicePlayer) writer() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		p.conn.Close()
	}()
	for {
		select {
		case <-ticker.C:
			p.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := p.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}

	}
}
