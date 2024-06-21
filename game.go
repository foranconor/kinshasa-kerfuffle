package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func init() {
	rl.InitWindow(800, 450, "Kinshasa Kerfuffle")
}

func update() {

}

func draw() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.RayWhite)
	rl.DrawText("Kinshasa Kerfuffle", 190, 200, 20, rl.LightGray)
	rl.EndDrawing()
}

func main() {
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		update()
		draw()
	}
}
