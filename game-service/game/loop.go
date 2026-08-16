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

	err := g.scoreRepo.Update(ctx, p.UserID, p.Kills)
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
		if CirclesOverlap(m.X, m.Y, MobRadius, p.X, p.Y, ws.PlayerRadius) {
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

	p.X = clamp(p.X+dx*ws.PlayerSpeed*dt, 0, MapWidth)
	p.Y = clamp(p.Y+dy*ws.PlayerSpeed*dt, 0, MapHeight)
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

func NearestPlayer(x, y float64, players []*ws.Player) *ws.Player {
	var nearestP *ws.Player
	minDist := math.MaxFloat64

	for _, p := range players {
		if p.HP < 0 {
			continue
		}

		dist := math.Hypot(p.X - x, p.Y - y)

		if dist < minDist {
			nearestP = p
			minDist = dist
		}
	}

	return nearestP
}