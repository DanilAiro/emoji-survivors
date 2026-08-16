package game

import "math"

const (
	MobRadius = 12.0
	MobStartHP = 1
	MobSpeed = 60.0
	MobDamagePerTick = 1
)

type Mob struct {
	ID int
	X float64
	Y float64
	HP int
}

func (m *Mob) MoveToward(targetX, targetY, dt float64) {
	dx := targetX - m.X
	dy := targetY - m.Y

	dist := math.Hypot(dx, dy)
	if dist < 0.001 {
		return
	}

	step := MobSpeed * dt
	if step > dist {
		step = dist
	}

	m.X += dx / dist * step
    m.Y += dy / dist * step
} 