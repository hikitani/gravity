package gravity

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Position struct {
	X float64
	Y float64
}

func (p *Position) DistanceTo(other *Position) Distance {
	dx, dy := other.X-p.X, other.Y-p.Y
	return Distance{
		Distance: math.Sqrt(float64(dx*dx) + float64(dy*dy)),
		Dx:       dx,
		Dy:       dy,
	}
}

type Distance struct {
	Distance float64
	Dx       float64
	Dy       float64
}

type Velocity struct {
	X float64
	Y float64
}

type Render struct {
	Image *ebiten.Image
}

type MouseState uint8

const (
	MouseStateNone MouseState = iota
	MouseStatePressed
	MouseStateClicked
)

type MouseEvents struct {
	X              int
	Y              int
	LeftSideState  MouseState
	RightSideState MouseState
}

type GravityAttraction struct {
	Radious      float64
	Acceleration float64
}

type IsStaticMarker struct{}
