package game

type PlayerSnapshot struct {
	UserID   int     `json:"user_id"`
	Username string  `json:"username"`
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
	HP       int     `json:"hp"`
	Kills    int     `json:"kills"`
}

type MobSnapshot struct {
	ID int     `json:"id"`
	X  float64 `json:"x"`
	Y  float64 `json:"y"`
	HP int     `json:"hp"`
}

type GameState struct {
	Type    string           `json:"type"`
	Players []PlayerSnapshot `json:"players"`
	Mobs    []MobSnapshot    `json:"mobs"`
}
