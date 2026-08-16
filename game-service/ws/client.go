package ws

import (
	"sync"

	"github.com/gorilla/websocket"
)

const (
	PlayerSpeed   = 200.0
	PlayerRadius  = 15.0
	PlayerStartHP = 100
)

type Player struct {
	UserID   int
	Username string

	X, Y  float64
	HP    int
	Kills int

	conn *websocket.Conn
	send chan []byte

	inputMu sync.Mutex
	dx, dy  float64
	attack  bool
}

func NewPlayer(userID int, username string, conn *websocket.Conn) *Player {
	return &Player{
		UserID:   userID,
		Username: username,

		X:  400,
		Y:  300,
		HP: PlayerStartHP,

		conn: conn,
		send: make(chan []byte, 16),
	}
}
