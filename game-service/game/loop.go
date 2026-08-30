package game

import (
	"context"
	"encoding/json"
	"log"
	"math"
	"math/rand"
	"time"

	"emoji-survivors/game-service/repository"
	"emoji-survivors/game-service/ws"
)

const (
	BaseMobCount = 5
	TickRate     = 30
	MapWidth     = 800
	MapHeight    = 600
	SwordRange   = 40.0
)

type GameLoop struct {
	nextMobID int
	hub       *ws.Hub
	mobs      map[int]*Mob
	scoreRepo *repository.ScoreRepository
}

func NewGameLoop(hub *ws.Hub, scoreRepo *repository.ScoreRepository) *GameLoop {
	return &GameLoop{
		hub:       hub,
		scoreRepo: scoreRepo,
		mobs:      make(map[int]*Mob),
	}
}

func (g *GameLoop) onPlayerDeath(p *ws.Player) {
	log.Printf("Игрок %d умер с %d очков", p.UserID, p.Kills)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := g.scoreRepo.UpdateIfHigher(ctx, p.UserID, p.Kills)
	if err != nil {
		log.Printf("Не удалось сохранить новые данные пользователя %d - %d очков", p.UserID, p.Kills)
	}
}

func (g *GameLoop) resolveAttack(p *ws.Player) {
	for id, m := range g.mobs {
		if CirclesOverlap(m.X, m.Y, MobRadius, p.X, p.Y, SwordRange) {
			delete(g.mobs, id)
			p.Kills++
		}
	}
}

func (g *GameLoop) resolveMobDamage(p *ws.Player) {
	for _, m := range g.mobs {
		if CirclesOverlap(m.X, m.Y, MobRadius, p.X, p.Y, ws.Radius) {
			p.HP -= MobDamagePerTick
		}
	}
}

func movePlayer(p *ws.Player, dx, dy, dt float64) {
	length := math.Hypot(dx, dy)
	if length < 0.001 {
		return
	}

	dx, dy = dx/length, dy/length

	p.X = clamp(p.X+dx*ws.Speed*dt, 0, MapWidth)
	p.Y = clamp(p.Y+dy*ws.Speed*dt, 0, MapHeight)
}

func clamp(v, min, max float64) float64 {
	if v < min {
		return min
	}

	if v > max {
		return max
	}

	return v
}

func nearestPlayer(x, y float64, players []*ws.Player) *ws.Player {
	var nearestP *ws.Player
	minDist := math.MaxFloat64

	for _, p := range players {
		if p.HP <= 0 {
			continue
		}

		dist := math.Hypot(p.X-x, p.Y-y)

		if dist < minDist {
			nearestP = p
			minDist = dist
		}
	}

	return nearestP
}

func (g *GameLoop) broadcastState() {
	players := g.hub.Players()

	state := GameState{
		Type:    "state",
		Players: make([]PlayerSnapshot, 0, len(players)),
		Mobs:    make([]MobSnapshot, 0, len(g.mobs)),
	}

	for _, p := range players {
		state.Players = append(state.Players, PlayerSnapshot{
			UserID:   p.UserID,
			Username: p.Username,
			X:        p.X,
			Y:        p.Y,
			HP:       p.HP,
			Kills:    p.Kills,
		})
	}
	for _, m := range g.mobs {
		state.Mobs = append(state.Mobs, MobSnapshot{ID: m.ID, X: m.X, Y: m.Y, HP: m.HP})
	}

	data, err := json.Marshal(state)
	if err != nil {
		log.Printf("не удалось сериализовать состояние игры: %v", err)
		return
	}

	for _, p := range players {
		p.Send(data)
	}
}

func (g *GameLoop) Run(ctx context.Context) {
	ticker := time.NewTicker(time.Second / TickRate)
	defer ticker.Stop()

	dt := 1.0 / float64(TickRate)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			g.tick(dt)
		}
	}
}

func (g *GameLoop) tick(dt float64) {
	playerCount := g.hub.Count()

	g.spawnMobs(playerCount)
	g.moveMobs(dt)
	g.handlePlayers()
	g.broadcastState()
}

func (g *GameLoop) spawnMobs(playerCount int) {
	if playerCount == 0 {
		return
	}

	desired := BaseMobCount * playerCount
	for len(g.mobs) < desired {
		g.nextMobID++
		g.mobs[g.nextMobID] = &Mob{
			ID: g.nextMobID,
			X:  rand.Float64() * MapWidth,
			Y:  rand.Float64() * MapHeight,
			HP: MobStartHP,
		}
	}
}

func (g *GameLoop) moveMobs(dt float64) {
	players := g.hub.Players()
	if len(players) == 0 {
		return
	}

	for _, m := range g.mobs {
		target := nearestPlayer(m.X, m.Y, players)
		if target != nil {
			m.MoveToward(target.X, target.Y, dt)
		}
	}
}

func (g *GameLoop) handlePlayers() {
	dt := 1.0 / float64(TickRate)
	players := g.hub.Players()

	for _, p := range players {
		if p.HP <= 0 {
			continue
		}

		dx, dy, attack := p.ConsumeInput()
		movePlayer(p, dx, dy, dt)

		if attack {
			g.resolveAttack(p)
		}

		g.resolveMobDamage(p)

		if p.HP <= 0 {
			g.onPlayerDeath(p)
		}
	}
}
