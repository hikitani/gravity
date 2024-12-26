package gravity

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	blockSize = 16
)

var (
	blockImage        = ebiten.NewImage(blockSize, blockSize)
	gravityBlockImage = ebiten.NewImage(blockSize, blockSize)
)

func init() {
	blockImage.Fill(color.White)
	gravityBlockImage.Fill(color.RGBA{255, 100, 100, 0})
}
