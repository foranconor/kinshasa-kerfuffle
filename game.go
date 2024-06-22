package main

import (
	"image/color"

	"github.com/foranconor/kinshasa-kerfuffle/banana"
	"github.com/foranconor/kinshasa-kerfuffle/gorilla"
	"github.com/foranconor/kinshasa-kerfuffle/scape"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	sW    = 800
	sH    = 450
	title = "Kinshasa Kerfuffle"
)

var (
	gameOver bool
	paused   bool
	players  []gorilla.Gorilla
	city     scape.Scape
	bananas  []banana.Banana
)

func init() {

	rl.InitWindow(sW, sH, title)

	city = scape.InitScape(sW, sH)
	players = gorilla.InitGorillas(city, 2)

}

func update() {

}

func draw() {
	rl.BeginDrawing()
	rl.ClearBackground(color.RGBA{121, 212, 253, 255})
	drawHud()
	drawBuildings()
	drawGorillas()
	rl.EndDrawing()

}

func drawHud() {
	// city name
	rl.DrawText(city.Name, 10, 10, 20, rl.White)
}

func drawBuildings() {
	for _, building := range city.Buildings {
		rl.DrawRectangleRec(building.Outline, building.Color.Color)
		shade := rl.Rectangle{
			X:      building.Outline.X + building.Outline.Width - 10,
			Y:      building.Outline.Y,
			Width:  10,
			Height: building.Outline.Height,
		}
		rl.DrawRectangleRec(shade, building.Color.Shadow)
		for _, window := range building.Windows {
			rl.DrawRectangleRec(window, building.Color.Window)
		}
	}
}

func drawGorillas() {
	for _, gorilla := range players {
		if gorilla.IsAlive {
			rl.DrawRectangleV(gorilla.Pos, gorilla.Size, rl.Red)
		}
	}
}

func main() {
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		update()
		draw()
	}
}
