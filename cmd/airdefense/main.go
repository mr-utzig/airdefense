package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
)

func main() {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	// Initialize audio context. Sample rate can be adjusted (e.g., 44100 or 48000).
	audioContext := audio.NewContext(48000)

	g := &Game{
		Airdefense:            NewAirdesense(),
		Airplanes:             []*Airplane{},
		NextAirplaneSpawnTick: MinAirplaneSpawnCooldownTicks + rand.Intn(RandomAirplaneSpawnCooldownRange),
		AudioContext:          audioContext,
		Explosions:            []Explosion{},
		AirplanesDestroyed:    0,
	}

	ebiten.SetWindowSize(WindowWidth, WindowHeight)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
