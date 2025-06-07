package main

import (
	"embed"
	"image"
	_ "image/png"
	"io"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed  assets/*
var assets embed.FS

//go:embed  sounds/*
var sounds embed.FS

var airdefense = mustLoadImage("assets/airdefense.png")
var irondome = mustLoadImage("assets/irondome.png")

var airplane = mustLoadImage("assets/airplane.png")
var bomb = mustLoadImage("assets/bomb.png")
var missile = mustLoadImage("assets/missile.png")

var fire1 = mustLoadImage("assets/fire1.png")
var fire2 = mustLoadImage("assets/fire2.png")
var fire3 = mustLoadImage("assets/fire3.png")

var explosionFrames = []*ebiten.Image{fire1, fire2, fire3, fire2} // Adiciona uma sequência para a animação

var ground = mustLoadImage("assets/ground.png")

var launch = mustLoadSound("sounds/missile.mp3")
var explosion = mustLoadSound("sounds/explosion.mp3")

func mustLoadImage(name string) *ebiten.Image {
	f, err := assets.Open(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}

	return ebiten.NewImageFromImage(img)
}

func mustLoadSound(name string) []byte {
	f, err := sounds.Open(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	bytes, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}
	return bytes
}
