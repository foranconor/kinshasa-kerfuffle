package scape

import (
	"github.com/kr/pretty"
	"testing"
)

func TestColors(t *testing.T) {
	n := 10
	colors := Colors(n)
	if len(colors) != n {
		t.Errorf("Expected %d colors, got %d", n, len(colors))
	}
	for i, c := range colors {
		if c.R == 0 && c.G == 0 && c.B == 0 {
			t.Errorf("Color %d is black", i)
		}
		pretty.Println(c)
	}

}
