package ws

import (
	"sync"
)

type Hub struct {
  mu      sync.RWMutex
  players map[int]*Player
}

func NewHub() *Hub {
  return &Hub{players: make(map[int]*Player)}
}


func (h *Hub) Add(p *Player) {
  h.mu.Lock()
  defer h.mu.Unlock()
  h.players[p.UserID] = p
}

func (h *Hub) Remove(userID int) {
  h.mu.Lock()
  defer h.mu.Unlock()
  if p, ok := h.players[userID]; ok {
    close(p.send)
    delete(h.players, userID)
  }
}

func (h *Hub) Count() int {
  h.mu.RLock()
  defer h.mu.RUnlock()
  return len(h.players)
}


func (h *Hub) Players() []*Player {
  h.mu.RLock()
  defer h.mu.RUnlock()
  result := make([]*Player, 0, len(h.players))
  for _, p := range h.players {
    result = append(result, p)
  }
  return result
}
