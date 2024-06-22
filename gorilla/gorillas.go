package gorilla

import (
	"github.com/foranconor/kinshasa-kerfuffle/scape"
	rl "github.com/gen2brain/raylib-go/raylib"
	"math/rand"
)

type Gorilla struct {
	Pos     rl.Vector2
	Size    rl.Vector2
	Aim     rl.Vector2
	PrevAim rl.Vector2

	Team    int
	IsAlive bool
	IsHuman bool
}

func InitGorillas(city scape.Scape, players int) []Gorilla {
	gorillas := make([]Gorilla, players)
	// divide the city into a as many sections as there are players with a DMZ in the middle between them
	// place a gorilla on a random building in each section
	sectionLength := len(city.Buildings) / (players*2 - 1)
	j := 0
	for i := 0; i < players; i++ {
		// pick a random building
		b := city.Buildings[j*sectionLength+rand.Intn(sectionLength)]
		// place the gorilla on the building
		gorillas[i].Pos = rl.Vector2{
			X: b.Outline.X + b.Outline.Width/2,
			Y: b.Outline.Y - 10,
		}
		gorillas[i].Size = rl.Vector2{
			X: 10,
			Y: 10,
		}
		gorillas[i].Team = i
		gorillas[i].IsAlive = true
		gorillas[i].IsHuman = true
		j += 2
	}
	return gorillas
}
