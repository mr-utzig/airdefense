package utils

import (
	"github.com/hajimehoshi/ebiten/v2"
)

func Rotate(theta float64, img *ebiten.Image, opt *ebiten.DrawImageOptions) {
	hw := float64(img.Bounds().Dx() / 2)
	hh := float64(img.Bounds().Dy() / 2)

	opt.GeoM.Translate(-hw, -hh)
	opt.GeoM.Rotate(theta)
	opt.GeoM.Translate(hw, hh)
}
