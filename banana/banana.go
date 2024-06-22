package banana

import rl "github.com/gen2brain/raylib-go/raylib"

type Banana struct {
	Pos      rl.Vector2
	Speed    rl.Vector2
	Rotation float32
	radius   float32
	Status   int
	Active   bool
}

type Explosion struct {
	Pos    rl.Vector2
	Radius float32
	Active bool
}
