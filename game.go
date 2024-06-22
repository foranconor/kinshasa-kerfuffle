package main

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"

	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/foranconor/kinshasa-kerfuffle/banana"
	"github.com/foranconor/kinshasa-kerfuffle/gorilla"
	"github.com/foranconor/kinshasa-kerfuffle/scape"
	"github.com/foranconor/kinshasa-kerfuffle/tools"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/kr/pretty"
)

const (
	sW         = 1700
	sH         = 800
	title      = "Kinshasa Kerfuffle"
	numPlayers = 2
)

var (
	gameOver   bool
	paused     bool
	bananaSent bool
	turn       int
	players    []gorilla.Gorilla
	city       scape.Scape
	banane     banana.Banana
	explosions []banana.Explosion
	spin       float32    // spin of the banana in rads/s
	sky        color.RGBA = color.RGBA{121, 212, 253, 255}
	gTop       color.RGBA = color.RGBA{0, 0, 255, 128}
	gBottom    color.RGBA = color.RGBA{0, 255, 0, 128}
	lights     float32
)

func gameState() string {
	if gameOver {
		return "Game Over"
	}
	if paused {
		return "Paused"
	}
	out := fmt.Sprintf("Player %s's turn\n", players[turn].Name)
	out += fmt.Sprintf("Banana sent: %t\n", bananaSent)
	out += fmt.Sprintf("Banana: x: %0.2f, y: %0.2f\n", banane.Pos.X, banane.Pos.Y)
	return out

}

func init() {
	rl.InitWindow(sW, sH, title)
	setup()
}

func setup() {
	city = scape.InitScape(sW, sH)
	gTop, gBottom = tools.CityToGradient(city.City)
	players = gorilla.InitGorillas(city, numPlayers)
	turn = 0
	gameOver = false
	paused = false
	bananaSent = false
	banane = banana.Banana{}
	explosions = make([]banana.Explosion, 0)
	spin = 0
	sun := tools.CityLight(city.City)
	switch sun {
	case "day":
		lights = 1
	case "civil":
		lights = 0.8
	case "nautical":
		lights = 0.5
	case "astronomical":
		lights = 0.2
	default:
		lights = 0.7
	}
}

func update() {
	//pretty.Println(gameState())
	if !gameOver {
		if rl.IsKeyPressed(rl.KeyP) {
			paused = !paused
		}
		if !paused {
			if !bananaSent {
				// aiming
				bananaSent = updateGorilla(turn)

			} else {
				// firing
				if updateBanana() {
					alive := 0
					for _, gorilla := range players {
						if gorilla.IsAlive {
							alive++
						}
					}
					pretty.Println(alive)
					if alive == 1 {
						gameOver = true
					} else {
						bananaSent = false
						banane.Active = false
						turn++
						if turn >= len(players) {
							turn = 0
						}
					}
				}
			}
		}
	} else {
		if rl.IsKeyPressed(rl.KeyEnter) {
			setup()
		}
	}

}

func updateGorilla(i int) bool {
	if !players[i].IsAlive {
		turn++
		if turn >= len(players) {
			turn = 0
		}
		return false
	}

	if rl.IsMouseButtonDown(rl.MouseLeftButton) {
		players[i].Aim = rl.GetMousePosition()
	}

	wheel := rl.GetMouseWheelMove()
	if wheel != 0 {
		spin += float32(wheel)
	}

	if rl.IsKeyPressed(rl.KeySpace) {
		// fire the banana with a scaled vector opposite to the aim
		speed := rl.Vector2Subtract(players[i].Pos, players[i].Aim)
		speed = rl.Vector2Scale(speed, 0.1)
		opposite := rl.Vector2{
			X: -speed.X,
			Y: -speed.Y,
		}
		banane = banana.Banana{
			Pos:    players[i].Pos,
			Speed:  opposite,
			Active: true,
		}
		return true
	}

	return false
}

func updateBanana() bool {
	if banane.Active {
		banane.Pos = rl.Vector2Add(banane.Pos, banane.Speed)
		// apply gravity
		banane.Speed.Y += 0.2
		// apply wind resistance
		// scale speed vector proportional to the speed and elevation
		elevation := int(city.City.Elevation)
		kpa := 101.325 * math.Pow(1-0.0065*float64(elevation)/288.15, 5.2559)
		inv := 1/kpa - 1/101.325
		pretty.Println(kpa)
		pretty.Println(inv)
		banane.Speed = rl.Vector2Scale(banane.Speed, float32(0.99-inv))

	}
	// check if off the screen
	if banane.Pos.X < 0 || banane.Pos.X > float32(sW) || banane.Pos.Y > float32(sH) {
		banane.Active = false
		bananaSent = false
		return true
	}
	// check for player collision
	for i, gorilla := range players {
		if gorilla.IsAlive && i != turn {
			if rl.CheckCollisionPointCircle(banane.Pos, gorilla.Pos, 10) {
				// gorilla hit
				players[i].IsAlive = false
				bananaSent = false
				explosions = append(explosions, banana.Explosion{
					Pos:    gorilla.Pos,
					Radius: 50,
					Active: true,
				})
				return true
			}
		}
	}
	// check for inside an explosion
	for _, explosion := range explosions {
		if explosion.Active {
			if rl.CheckCollisionPointCircle(banane.Pos, explosion.Pos, explosion.Radius) {
				// hit an explosion
				return false
			}
		}
	}

	// check for building collision
	for _, building := range city.Buildings {
		if rl.CheckCollisionPointRec(banane.Pos, building.Outline) {
			// hit a building
			rad := rand.Intn(50) + 25

			explosions = append(explosions, banana.Explosion{
				Pos:    banane.Pos,
				Radius: float32(rad),
				Active: true,
			})
			// switch the banana off
			banane.Active = false
			bananaSent = false
			// check if the explosion hit a gorilla
			for i, gorilla := range players {
				if gorilla.IsAlive {
					if rl.CheckCollisionPointCircle(gorilla.Pos, banane.Pos, float32(rad)) {
						// gorilla hit
						pretty.Println("Gorilla hit")
						players[i].IsAlive = false
						explosions = append(explosions, banana.Explosion{
							Pos:    gorilla.Pos,
							Radius: 200,
							Active: true,
						})
					}
				}
			}
			return true
		}
	}
	return false
}

func draw() {
	rl.BeginDrawing()
	rl.ClearBackground(sky)
	drawBuildings()
	drawExplosions()
	drawGorillas()
	drawBanana()
	rl.DrawRectangleGradientV(0, 0, sW, sH, gTop, gBottom)
	drawWindows()
	drawAim()
	drawHud()
	rl.EndDrawing()

}

func drawHud() {
	// location details
	c := city.City
	rl.DrawText(fmt.Sprintf("%s, %s", c.Name, c.Country), 10, 10, 20, rl.White)
	pr := message.NewPrinter(language.English)
	pop := pr.Sprintf("Population: %d", c.Population)
	rl.DrawText(pop, 10, 30, 20, rl.White)
	elev := pr.Sprintf("Elevation: %dm", c.Elevation)
	rl.DrawText(elev, 10, 50, 20, rl.White)
	latlon := fmt.Sprintf("Lat: %0.2f°, Lon: %0.2f°", c.Latitude, c.Longitude)
	rl.DrawText(latlon, 10, 70, 20, rl.White)
	t := pr.Sprintf("Time: %s", c.Time.Format("15:04"))
	rl.DrawText(t, 10, 90, 20, rl.White)

	if gameOver {
		txt := "Press Enter to restart"
		txtSize := 100
		rl.DrawText(txt, sW/2-rl.MeasureText(txt, int32(txtSize))/2, sH/2, int32(txtSize), rl.White)
	}
	if paused {
		txt := "Press P to unpause"
		txtSize := 100
		rl.DrawText(txt, sW/2-rl.MeasureText(txt, int32(txtSize))/2, sH/2, int32(txtSize), rl.White)
	}
	// draw aim data in the top center
	angle := rl.Vector2Subtract(players[turn].Pos, players[turn].Aim)
	degrees := math.Atan2(float64(angle.Y), float64(angle.X)) * 180 / math.Pi
	power := rl.Vector2Length(angle)
	angTxt := fmt.Sprintf("Angle: %0.2f", degrees)
	powTxt := fmt.Sprintf("Power: %0.2f", power)
	spinTxt := fmt.Sprintf("Spin: %0.2f", spin)
	rl.DrawText(angTxt, sW/2-rl.MeasureText(angTxt, 20)/2, 10, 20, rl.White)
	rl.DrawText(powTxt, sW/2-rl.MeasureText(powTxt, 20)/2, 30, 20, rl.White)
	rl.DrawText(spinTxt, sW/2-rl.MeasureText(spinTxt, 20)/2, 50, 20, rl.White)
	// draw player data
	y := int32(10)
	for i, gorilla := range players {
		txt := gorilla.Name
		if i == turn {
			txt = fmt.Sprintf("# %s", gorilla.Name)
		}
		if i == turn {
			rl.DrawText(txt, sW-rl.MeasureText(txt, 20)-10, y, 20, rl.White)
		} else if gorilla.IsAlive {
			rl.DrawText(txt, sW-rl.MeasureText(txt, 20)-10, y, 20, rl.LightGray)

		} else {
			rl.DrawText(txt, sW-rl.MeasureText(txt, 20)-10, y, 20, rl.Gray)
		}
		y += 20
	}
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

func drawWindows() {
	thing := rand.New(rand.NewSource(99))

	for _, building := range city.Buildings {
		for _, window := range building.Windows {
			// 50% chance of a light being on
			mid := rl.Vector2{
				X: window.X + window.Width/2,
				Y: window.Y + window.Height/2,
			}
			if thing.Float32() > lights {
				// check if it's exploded
				exploded := false
				for _, explosion := range explosions {
					if explosion.Active {
						if rl.CheckCollisionPointCircle(mid, explosion.Pos, explosion.Radius+100) {
							exploded = true
						}
					}
				}
				if !exploded {
					rl.DrawRectangleRec(window, rl.Yellow)
				}
			}
		}
	}
}

func drawGorillas() {
	for _, gorilla := range players {
		if gorilla.IsAlive {
			rl.DrawRectangleV(gorilla.Pos, gorilla.Size, rl.Red)
			rl.DrawText(gorilla.Name, int32(gorilla.Pos.X), int32(gorilla.Pos.Y), 10, rl.White)
		} else {
			rl.DrawRectangleV(gorilla.Pos, gorilla.Size, rl.Gray)
			rl.DrawText(gorilla.Name, int32(gorilla.Pos.X), int32(gorilla.Pos.Y), 10, rl.Gray)
		}
	}
}

func drawExplosions() {
	for _, explosion := range explosions {
		if explosion.Active {
			rl.DrawCircleV(explosion.Pos, explosion.Radius, sky)
		}
	}
}

func drawBanana() {
	if banane.Active {
		rl.DrawCircleV(banane.Pos, 5, rl.Yellow)
	}
}

func drawAim() {
	if !players[turn].IsAlive {
		return
	}
	rl.DrawLineV(players[turn].Pos, players[turn].Aim, rl.Red)
	rl.DrawCircleV(players[turn].Pos, 10, rl.Red)
	rl.DrawCircleV(players[turn].Aim, 5, rl.Red)
}

func main() {
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		update()
		draw()
	}
}
