package scape

import (
	"image/color"
	"math/rand"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/lucasb-eyer/go-colorful"
)

var types = []string{
	"Apartments",
	"Office",
	"Retail",
	"Government",
}

type City struct {
	Name       string
	Country    string
	Population int
	Elevation  int
	Latitude   float64
	Longitude  float64
	Time       time.Time
}

var cities = []City{
	{
		Name:       "Kinshasa",
		Country:    "DRC",
		Population: 17032322,
		Elevation:  240,
		Latitude:   -4.321944,
		Longitude:  15.311944,
	},
	{
		Name:       "Brazzaville",
		Country:    "Congo",
		Population: 2145783,
		Elevation:  320,
		Latitude:   -4.269444,
		Longitude:  15.271111,
	},
	{
		Name:       "Pointe-Noire",
		Country:    "Congo",
		Population: 1420612,
		Elevation:  0,
		Latitude:   -4.7975,
		Longitude:  11.850278,
	},
	{
		Name:       "Libreville",
		Country:    "Gabon",
		Population: 703904,
		Elevation:  0,
		Latitude:   0.390278,
		Longitude:  9.454167,
	},
	{
		Name:       "Yaoundé",
		Country:    "Cameroon",
		Population: 2765600,
		Elevation:  726,
		Latitude:   3.866667,
		Longitude:  11.516667,
	},
	{
		Name:       "Douala",
		Country:    "Cameroon",
		Population: 5066000,
		Elevation:  13,
		Latitude:   4.05,
		Longitude:  9.683333,
	},
	{
		Name:       "Bujumbura",
		Country:    "Burundi",
		Elevation:  774,
		Population: 1143202,
		Latitude:   -3.383333,
		Longitude:  29.366667,
	},
	{
		Name:       "Bata",
		Country:    "Equatorial Guinea",
		Population: 250770,
		Elevation:  5,
		Latitude:   1.863611,
		Longitude:  9.765833,
	},
	{
		Name:       "Bafoussam",
		Country:    "Cameroon",
		Population: 1146000,
		Elevation:  1521,
		Latitude:   5.466667,
		Longitude:  10.416667,
	},
	{
		Name:       "Uyo",
		Country:    "Nigeria",
		Population: 554906,
		Elevation:  70,
		Latitude:   5.033333,
		Longitude:  7.9275,
	},
	{
		Name:       "Calabar",
		Country:    "Nigeria",
		Population: 571500,
		Elevation:  32,
		Latitude:   4.976667,
		Longitude:  8.338333,
	},
	{
		Name:       "Bukavu",
		Country:    "DRC",
		Population: 1133000,
		Elevation:  1498,
		Latitude:   -2.5,
		Longitude:  28.866667,
	},
	{
		Name:       "Goma",
		Country:    "DRC",
		Population: 670000,
		Elevation:  1530,
		Latitude:   -1.679444,
		Longitude:  29.233611,
	},
	{
		Name:       "Bangî",
		Country:    "CAR",
		Population: 812407,
		Elevation:  369,
		Latitude:   4.366667,
		Longitude:  18.583333,
	},
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
	City      City
	Buildings []Building
	Backdrop  [][]Building
}

var MaxBuildings = 40
var MinBuildings = 28

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
		City: cities[rand.Intn(len(cities))],
	}
	hour := rand.Intn(24)
	minute := rand.Intn(60)
	t := time.Date(2020, 1, 1, hour, minute, 0, 0, time.UTC)
	scape.City.Time = t
	z := 1
	if scape.City.Population > 2000000 {
		z = 2
	} else if scape.City.Population < 600000 {
		z = 0
	}
	scape.Buildings = generateBuildings(x, y, z)
	for i := 0; i < z+1; i++ {
		scape.Backdrop = append(scape.Backdrop, generateBuildings(x, y, z))
	}
	return scape
}

func generateBuildings(x, y, z int) []Building {
	row := make([]Building, 0)
	buildings := rand.Intn(MaxBuildings-MinBuildings) + MinBuildings
	width := x / buildings
	i := 0
	w := 0.0
	for i < buildings {
		height := (rand.Intn(y/2) + y/10)
		switch z {
		case 0:
			height -= rand.Intn(y / 8)
			if height < y/10 {
				height = y / 10
			}
		case 2:
			height += rand.Intn(y / 3)
		}

		bWidth := width + rand.Intn(width)
		if w+float64(bWidth) > float64(x-20) {
			bWidth = int(float64(x) - w)
		}
		building := Building{
			Cost:    rand.Intn(1000) + 100,
			Outline: rl.Rectangle{X: float32(w), Y: float32(y - height), Width: float32(bWidth), Height: float32(height)},
		}
		row = append(row, building)
		w += float64(bWidth)
		if w >= float64(x) {
			break
		}
		i++
	}
	colors := Colors(len(row))
	for i, building := range row {
		row[i] = fleshOutBuilding(building, colors[i])
	}
	return row
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
