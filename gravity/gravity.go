package gravity

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hikitani/blueprint/ecs"
)

const (
	ScreenWidth  = 600
	ScreenHeight = 600
)

var _ ebiten.Game = &Game{}

type Game struct {
	*ecs.World
}

func NewGame() *Game {
	w := ecs.New(ScreenWidth, ScreenHeight)

	w.AddEntity(&Input{})
	w.AddDrawer(&RenderSystem{})
	w.
		AddLogic(&InputChecker{}).
		AddLogic(&BlockSpawnerByClick{}).
		AddLogic(&GravitySystem{}).
		AddLogic(&MovementSystem{})

	return &Game{
		World: w,
	}
}
