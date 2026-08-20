package ws

import (
	"encoding/json"
	"sync"

	"github.com/gorilla/websocket"
)

const (
	StartHP = 100
	Radius  = 15.0
	Speed   = 200.0
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

type clientCommand struct {
	Action string  `json:"action"` // move или attack
	DX     float64 `json:"dx"`
	DY     float64 `json:"dy"`
}

func NewPlayer(userID int, username string, conn *websocket.Conn) *Player {
	return &Player{
		UserID:   userID,
		Username: username,

		X:  400,
		Y:  300,
		HP: StartHP,

		conn: conn,
		send: make(chan []byte, 16),
	}
}

func (p *Player) WritePump() {
	defer p.conn.Close()

	for msg := range p.send {
		if err := p.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			return
		}
	}
}

func (p *Player) ReadPump(onClose func()) {
	defer onClose()

	for {
		_, msg, err := p.conn.ReadMessage()
		if err != nil {
			return
		}

		var cmd clientCommand
		if err := json.Unmarshal(msg, &cmd); err != nil {
			continue
		}

		switch cmd.Action {
		case "move":
			p.SetMove(cmd.DX, cmd.DY)
		case "attack":
			p.TriggerAttack()
		}
	}
}

func (p *Player) Send(msg []byte) {
	defer func() {
		_ = recover()
	}()

	select {
	case p.send <- msg:
	default:
	}
}

func (p *Player) SetMove(dx, dy float64) {
  p.inputMu.Lock()
  defer p.inputMu.Unlock()
  p.dx, p.dy = dx, dy
}

func (p *Player) TriggerAttack() {
  p.inputMu.Lock()
  defer p.inputMu.Unlock()
  p.attack = true
}

func (p *Player) ConsumeInput() (dx, dy float64, attack bool) {
	p.inputMu.Lock()
	defer p.inputMu.Unlock()
	dx, dy, attack = p.dx, p.dy, p.attack
	p.attack = false
	return
}
