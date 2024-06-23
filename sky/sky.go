package sky

import (
	"image/color"
	"math/rand"

	"github.com/foranconor/kinshasa-kerfuffle/scape"
	"github.com/foranconor/kinshasa-kerfuffle/tools"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Sun struct {
	Radiance float32
	Pos      rl.Vector2
	Color    rl.Color
}

type Moon struct {
	Radiance float32
	Pos      rl.Vector2
	Color    rl.Color
}

type Star struct {
	Pos        rl.Vector2
	Color      rl.Color
	Brightness float32
}

type Particle struct {
	Pos   rl.Vector2
	Speed rl.Vector2
	Color rl.Color
}

type Cloud struct {
	Pos   rl.Vector2
	Kind  string // "cirrus", "cumulus", "stratus"
	Color rl.Color
}
type Sky struct {
	Color     color.RGBA
	Sun       Sun
	Moon      Moon
	Stars     []Star
	Particles []Particle
	Clouds    []Cloud
}

func InitSky(c scape.City, x, y float32) Sky {
	var s Sky

	switch tools.CityLight(c) {

	case "day":
		// suns y position depends on the time of day
		// suns x position
		s = Sky{
			Color: color.RGBA{255, 255, 255, 255},
			Sun: Sun{
				Radiance: 1.0,
				Pos:      rl.Vector2{X: rand.Float32() * x, Y: 100},
				Color:    rl.NewColor(255, 255, 0, 255),
			},
		}

	case "civil":
		s = Sky{
			Color: color.RGBA{255, 255, 255, 255},
			Sun: Sun{
				Radiance: 0.5,
				Pos:      rl.Vector2{X: rand.Float32() * x, Y: 400},
				Color:    rl.NewColor(255, 255, 0, 255),
			},
		}
	case "nautical":
		s = Sky{
			Color: color.RGBA{0, 0, 60, 255},
		}
	case "astronomical":
		s = Sky{
			Color: color.RGBA{0, 0, 64, 255},
		}
		s.Stars = makeStars(100, x, y)
	case "night":
		s = Sky{
			Color: color.RGBA{0, 0, 0, 255},
			Moon: Moon{
				Radiance: 1.0,
				Pos:      rl.Vector2{X: rand.Float32() * x, Y: 100},
				Color:    rl.NewColor(255, 255, 255, 255),
			},
		}
		s.Stars = makeStars(1000, x, y)
	}
	// particles
	numParticles := 20
	for i := 0; i < numParticles; i++ {
		s.Particles = append(s.Particles, Particle{
			Pos:   rl.Vector2{X: rand.Float32() * x, Y: rand.Float32() * y},
			Speed: rl.Vector2{X: 0, Y: 0},
			Color: color.RGBA{100, 255, 100, 128},
		})
	}
	return s
}

func makeStars(n int, x, y float32) []Star {
	stars := make([]Star, 0)
	for i := 0; i < n; i++ {
		roll := rand.Intn(4)
		color := rl.NewColor(255, 255, 255, 255)
		switch roll {
		case 0:
			color.B = 200
		case 1:
			color.G = 200
		case 2:
			color.R = 200
		}

		stars = append(stars, Star{
			Pos: rl.Vector2{X: rand.Float32() * x, Y: rand.Float32() * y},

			Color:      color,
			Brightness: rand.Float32(),
		})
	}

	return stars
}
