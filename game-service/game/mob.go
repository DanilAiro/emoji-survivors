package game

type Mob struct {
	Name   string
	Health int
	Damage int
	Target *Player
}
