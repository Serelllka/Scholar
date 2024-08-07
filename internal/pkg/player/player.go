package player

import "github.com/gorilla/websocket"

type Player struct {
	conn *websocket.Conn
}
