package tools

import (
	"image/color"
	"time"

	"github.com/foranconor/kinshasa-kerfuffle/scape"
	"github.com/kr/pretty"
)

func CityToGradient(c scape.City) (color.RGBA, color.RGBA) {
	light := CityLight(c)
	// alpha is determined by the population of the city pop 3000000 = 255 alpha, pop 10000000 = 200 alpha
	alphaFactor := 255 / 10000000
	alpha := uint8(c.Population * alphaFactor)

	pretty.Println(c, light, alpha)
	switch light {
	case "day":
		alpha = 100
		top := color.RGBA{240, 240, 240, alpha}
		bottom := color.RGBA{170, 170, 170, alpha}
		return top, bottom
	case "civil":
		alpha = 128
		return color.RGBA{134, 48, 43, alpha}, color.RGBA{94, 113, 152, alpha}
	case "nautical":
		alpha = 220
		return color.RGBA{60, 10, 40, alpha}, color.RGBA{10, 10, 0, alpha}
	case "astronomical":
		alpha = 230
		return color.RGBA{0, 88, 181, alpha}, color.RGBA{10, 10, 0, alpha}
	default:
		alpha = 240
		return color.RGBA{0, 0, 0, alpha}, color.RGBA{10, 10, 0, alpha}
	}
}

func CityLight(c scape.City) string {
	// location details
	t := c.Time
	light := "night"
	civilAm := time.Date(t.Year(), t.Month(), t.Day(), 5, 42, 0, 0, time.UTC)
	civilPm := time.Date(t.Year(), t.Month(), t.Day(), 18, 18, 0, 0, time.UTC)
	nauticalAm := time.Date(t.Year(), t.Month(), t.Day(), 5, 15, 0, 0, time.UTC)
	nauticalPm := time.Date(t.Year(), t.Month(), t.Day(), 18, 19, 0, 0, time.UTC)
	astronomicalAm := time.Date(t.Year(), t.Month(), t.Day(), 4, 49, 0, 0, time.UTC)
	astronomicalPm := time.Date(t.Year(), t.Month(), t.Day(), 18, 45, 0, 0, time.UTC)
	nightAm := time.Date(t.Year(), t.Month(), t.Day(), 4, 49, 0, 0, time.UTC)
	nightPm := time.Date(t.Year(), t.Month(), t.Day(), 19, 11, 0, 0, time.UTC)
	if t.After(civilAm) && t.Before(civilPm) {
		light = "day"
	} else if t.After(civilPm) && t.Before(nauticalPm) || t.After(nauticalAm) && t.Before(civilAm) {
		light = "civil"
	} else if t.After(nauticalPm) && t.Before(astronomicalPm) || t.After(astronomicalAm) && t.Before(nauticalAm) {
		light = "nautical"
	} else if t.After(astronomicalPm) && t.Before(nightPm) || t.After(nightAm) && t.Before(astronomicalAm) {
		light = "astronomical"
	} else if t.After(nightPm) || t.Before(nightAm) {
		light = "night"
	}
	_ = light
	return light
}
