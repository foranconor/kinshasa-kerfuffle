package scape

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/lucasb-eyer/go-colorful"
	"image/color"
	"math/rand"
)

var types = []string{
	"Apartments",
	"Office",
	"Retail",
	"Government",
}

type Building struct {
	Type    string
	Name    string
	Levels  int
	Cost    int
	Outline rl.Rectangle
	Color   Paint
	Windows []rl.Rectangle
}

type Paint struct {
	Color  color.RGBA
	Shadow color.RGBA
	Window color.RGBA
}

type Scape struct {
	Name      string
	Buildings []Building
}

var MaxBuildings = 30
var MinBuildings = 10

func Colors(buildings int) []Paint {
	colors := make([]Paint, buildings)
	chroma := 0.5
	luminance := 0.7
	//hue := rand.Float64() * 360
	hue := 0.0
	step := 360.0 / float64(buildings)
	c := colorful.Hcl(hue, chroma, luminance)
	s := colorful.Hcl(hue, chroma, luminance-0.2)
	w := colorful.Hcl(hue, chroma+0.2, luminance-0.4)
	for i := 0; i < buildings; i += 1 {

		c = colorful.Hcl(hue, chroma, luminance)
		s = colorful.Hcl(hue, chroma, luminance-0.2)
		w = colorful.Hcl(hue, chroma+0.2, luminance-0.4)
		c = c.Clamped()
		s = s.Clamped()
		w = w.Clamped()
		r, g, b := c.RGB255()
		colors[i].Color = color.RGBA{r, g, b, 255}
		r, g, b = s.RGB255()
		colors[i].Shadow = color.RGBA{r, g, b, 255}
		r, g, b = w.RGB255()
		colors[i].Window = color.RGBA{r, g, b, 255}
		hue += step
		if hue > 360 {
			hue -= 360
		}
	}
	// shuffle colors
	rand.Shuffle(len(colors), func(i, j int) {
		colors[i], colors[j] = colors[j], colors[i]
	})
	return colors
}

func InitScape(x, y int) Scape {
	scape := Scape{
		Name: "Kinshasa",
	}
	buildings := rand.Intn(MaxBuildings-MinBuildings) + MinBuildings
	width := x / buildings
	i := 0
	w := 0.0
	for i < buildings {
		height := (rand.Intn(y/2) + y/10)
		bWidth := width + rand.Intn(width)
		if w+float64(bWidth) > float64(x-20) {
			bWidth = int(float64(x) - w)
		}
		building := Building{
			Cost:    rand.Intn(1000) + 100,
			Outline: rl.Rectangle{X: float32(w), Y: float32(y - height), Width: float32(bWidth), Height: float32(height)},
		}
		scape.Buildings = append(scape.Buildings, building)
		w += float64(bWidth)
		if w >= float64(x) {
			break
		}
		i++
	}
	colors := Colors(len(scape.Buildings))
	for i, building := range scape.Buildings {
		scape.Buildings[i] = fleshOutBuilding(building, colors[i])
	}

	return scape
}

func fleshOutBuilding(building Building, paint Paint) Building {
	floorHeight := 30
	windowHeight := rand.Intn(8) + 5
	windowWidth := rand.Intn(5) + 5
	building.Color = paint
	building.Levels = int(building.Outline.Height) / floorHeight // 20 pixels per floor
	floorHeight = int(building.Outline.Height) / building.Levels
	building.Windows = make([]rl.Rectangle, 0)
	windowSpacing := 20
	windowsPerFloor := int(building.Outline.Width / float32(windowSpacing))
	windowGap := (building.Outline.Width - float32(windowsPerFloor*windowWidth)) / float32(windowsPerFloor+1)
	for i := 0; i < building.Levels; i++ {
		for j := 0; j < int(windowsPerFloor); j++ {
			window := rl.Rectangle{
				X:      building.Outline.X + windowGap + float32(j*(windowWidth+int(windowGap))),
				Y:      building.Outline.Y + float32(i*floorHeight) + 5,
				Width:  float32(windowWidth),
				Height: float32(windowHeight),
			}
			building.Windows = append(building.Windows, window)
		}
	}
	building.Type = types[rand.Intn(len(types))]
	building.Name = randomName()

	return building
}

func randomName() string {
	return "Building"
}
