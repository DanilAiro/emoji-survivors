package game

import "math"

func CirclesOverlap(x1, y1, r1, x2, y2, r2 float64) bool {
	dist := math.Hypot(x2 - x1, y2 - y1)

	return dist < r1 + r2
}