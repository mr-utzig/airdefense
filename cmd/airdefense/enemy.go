package main

import (
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mr-utzig/airdefense/pkg/utils"
)

type Airplane struct {
	Image    *ebiten.Image
	Position Vector
	Speed    float64
}

func (a *Airplane) Draw(screen *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}

	if a.Speed < 0 {
		utils.Rotate(45.0*math.Pi, a.Image, opts)
	}

	opts.GeoM.Translate(a.Position.X, a.Position.Y)
	screen.DrawImage(a.Image, opts)
}

func (g *Game) spawnAirplane() {
	img := airplane
	imgWidth := float64(img.Bounds().Dx())
	imgHeight := float64(img.Bounds().Dy())

	minY := 10.0
	maxY := (WindowHeight * 0.4) - imgHeight
	if maxY < minY {
		maxY = minY
	}
	startY := minY + rand.Float64()*(maxY-minY)

	var startX float64
	var speed float64

	if rand.Intn(2) == 0 { // From left
		startX = -imgWidth
		speed = AirplaneSpeed
	} else { // From right
		startX = WindowWidth
		speed = -AirplaneSpeed
	}

	newAirplane := &Airplane{
		Image:    img,
		Position: Vector{X: startX, Y: startY},
		Speed:    speed,
	}
	g.Airplanes = append(g.Airplanes, newAirplane)
}
