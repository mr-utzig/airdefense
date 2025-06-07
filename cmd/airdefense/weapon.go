package main

import (
	"bytes"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/mr-utzig/airdefense/pkg/utils"
)

type Weapon struct {
	Image        *ebiten.Image
	Position     Vector
	Rotation     float64
	Projectile   Projectile
	Projectiles  []Projectile
	LastShotTick int
}

type Projectile struct {
	Image    *ebiten.Image
	Position *Vector
	Speed    float64
	Fire     *ebiten.Image
	Target   *Airplane
	Rotation float64
}

func (s Weapon) draw(screen *ebiten.Image) {
	if s.hasProjectiles() {
		for _, projectile := range s.Projectiles {
			missileOpt := &ebiten.DrawImageOptions{}
			utils.Rotate(projectile.Rotation, projectile.Image, missileOpt)
			missileOpt.GeoM.Translate(projectile.Position.X, projectile.Position.Y)

			if projectile.Fire != nil {
				fireWidth := float64(projectile.Fire.Bounds().Dx())
				fireHeight := float64(projectile.Fire.Bounds().Dy())

				fireOpt := &ebiten.DrawImageOptions{}
				utils.Rotate(projectile.Rotation+(45.0*math.Pi), projectile.Fire, fireOpt)

				flightAngle := projectile.Rotation - math.Pi/2.0
				// Distance from missile center to fire center
				offsetDistance := float64(projectile.Image.Bounds().Dy())*0.5 + fireHeight*0.1
				fireX := (projectile.Position.X + fireWidth/2) - offsetDistance*math.Cos(flightAngle) - fireWidth/2
				fireY := (projectile.Position.Y + fireHeight/2) - offsetDistance*math.Sin(flightAngle) - fireHeight/2
				fireOpt.GeoM.Translate(fireX, fireY)

				screen.DrawImage(projectile.Fire, fireOpt)
			}

			screen.DrawImage(projectile.Image, missileOpt)
		}
	}

	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(s.Position.X, s.Position.Y)
	screen.DrawImage(s.Image, opt)
}

func (s Weapon) hasProjectiles() bool {
	return len(s.Projectiles) > 0
}

func (s *Weapon) shoot(currentTick int, audioContext *audio.Context) {
	if currentTick-s.LastShotTick < ShotCooldownTicks {
		return
	}
	pos := s.Position
	newProjectile := s.Projectile
	projStartPos := Vector{X: pos.X, Y: pos.Y}
	newProjectile.Position = &projStartPos
	newProjectile.Target = nil

	s.Projectiles = append(s.Projectiles, newProjectile)
	s.LastShotTick = currentTick

	if audioContext != nil && len(launch) > 0 {
		// Decodifica os dados MP3 para PCM antes de reproduzir
		launchStream, err := mp3.DecodeWithSampleRate(audioContext.SampleRate(), bytes.NewReader(launch))
		if err == nil { // Verificação básica de erro
			soundPlayer, _ := audioContext.NewPlayer(launchStream)
			soundPlayer.Play()
		}
		// Considere registrar o erro se err != nil
	}
}

func NewAirdesense() Weapon {
	pos := Vector{
		X: ScreenWidth - float64((airdefense.Bounds().Dx() / 2)),
		Y: 8.5 + (float64(WindowHeight-ground.Bounds().Dy()) - float64(airdefense.Bounds().Dy())),
	}

	return Weapon{
		Image:    airdefense,
		Position: pos,
		Projectile: Projectile{
			Image: missile,
			Speed: 3.0,
			Fire:  fire1,
		},
	}
}

func NewIrondome() Weapon {
	pos := Vector{
		X: (ScreenWidth - float64(irondome.Bounds().Dx())) - 10,
		Y: 6.5 + (float64(WindowHeight-ground.Bounds().Dy()) - float64(irondome.Bounds().Dy())),
	}

	return Weapon{
		Image:    irondome,
		Position: pos,
		Projectile: Projectile{
			Image: bomb,
			Speed: 1,
			Fire:  fire1,
		},
	}
}
