package banana

import rl "github.com/gen2brain/raylib-go/raylib"

type Banana struct {
	Pos    rl.Vector2
	Speed  rl.Vector2
	Status int
	Active bool
}
