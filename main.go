package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hikitani/blueprint/gravity"
)

func main() {
	ebiten.SetWindowSize(gravity.ScreenWidth, gravity.ScreenHeight)
	ebiten.SetWindowTitle("Gravity")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(gravity.NewGame()); err != nil {
		log.Print(err)
	}
}
