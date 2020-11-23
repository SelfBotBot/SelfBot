package ws_player

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

var keys = make(map[uuid.UUID]*VoicePlayer)

func StartAuth(vp *VoicePlayer) uuid.UUID {
	auth := uuid.Must(uuid.NewV4())
	keys[auth] = vp
	go func() {
		// TODO gross.
		time.Sleep(30 * time.Second)
		delete(keys, auth)
	}()
	return auth
}

func IsAuthed(wsConn *websocket.Conn) (*VoicePlayer, bool) {
	wsConn.SetReadLimit(maxMessageSize)
	wsConn.SetReadDeadline(time.Now().Add(pongWait))

	var authRequest PlayRequest
	err := wsConn.ReadJSON(&authRequest)
	if err != nil {
		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			log.Printf("error: %v", err)
		}
		wsConn.Close()
		return &VoicePlayer{}, false
	}

	vp, ok := keys[authRequest.ID]
	if ok {
		delete(keys, authRequest.ID)
	}
	return vp, ok
}
