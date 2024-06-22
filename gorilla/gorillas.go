package gorilla

import (
	"github.com/foranconor/kinshasa-kerfuffle/scape"
	rl "github.com/gen2brain/raylib-go/raylib"
	"math/rand"
)

var names = []string{
	"Mek",
	"T'hu",
	"Fang",
	"Mik",
	"T'ses",
	"T'minvith",
	"V'lunval",
	"Sinas",
	"T'hukur",
	"Vinek",
}

type Gorilla struct {
	Name    string
	Pos     rl.Vector2
	Size    rl.Vector2
	Aim     rl.Vector2
	PrevAim rl.Vector2

	Team    int
	IsAlive bool
	IsHuman bool
}

func InitGorillas(city scape.Scape, players int) []Gorilla {
	// shuffle the names
	rand.Shuffle(len(names), func(i, j int) {
		names[i], names[j] = names[j], names[i]
	})

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
		gorillas[i].Name = names[i]
		j += 2
	}
	return gorillas
}
