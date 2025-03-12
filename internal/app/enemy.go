package app

import "math/rand"

type Enemy struct {
	X, Y         int
	ID           string
	Name         string
	HP           int
	Strength     int
	Agility      int
	Intelligence int
}

func NewEnemy(x, y int, level int) Enemy {
	return Enemy{
		X:            x,
		Y:            y,
		ID:           "enemy_" + string(rand.Intn(1000)),
		Name:         "Goblin",
		HP:           30 + level*5,
		Strength:     5 + level*2,
		Agility:      5 + level*2,
		Intelligence: 3 + level,
	}
}
