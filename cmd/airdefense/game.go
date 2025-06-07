package main

import (
	"bytes"
	"fmt"
	"image"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	ScreenWidth  = 240
	ScreenHeight = 240
	WindowWidth  = ScreenWidth * 2
	WindowHeight = ScreenHeight * 2

	FrameCount        = 10
	ShotCooldownTicks = 60

	AirplaneSpeed                    = 2.0
	ExplosionFrameDurationTicks      = 5 // Quantos ticks cada quadro da explosão dura
	MinAirplaneSpawnCooldownTicks    = 120
	RandomAirplaneSpawnCooldownRange = 180
)

type Vector struct {
	X float64
	Y float64
}

type Explosion struct {
	Position           Vector
	Frames             []*ebiten.Image
	CurrentFrameIndex  int
	FrameDurationTicks int
	CurrentTickInFrame int
}

type Game struct {
	Airdefense Weapon
	Irondome   Weapon
	Missile    Weapon
	Bomb       Weapon
	Airplanes  []*Airplane
	Explosions []Explosion

	AudioContext          *audio.Context
	Counter               int
	NextAirplaneSpawnTick int
	AirplanesDestroyed    int
}

func (g *Game) Update() error {
	g.Counter++

	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		g.Airdefense.shoot(g.Counter, g.AudioContext)
	}

	if g.Airdefense.hasProjectiles() {
		currentProjectiles := g.Airdefense.Projectiles
		g.Airdefense.Projectiles = []Projectile{}
		for i := range currentProjectiles {
			p := &currentProjectiles[i]

			// Target acquisition for tracking missiles
			if p.Target == nil && len(g.Airplanes) > 0 {
				// Simple: target the first available airplane.
				p.Target = g.Airplanes[0]
			}

			// Movement
			if p.Target != nil {
				// Check if the target is still valid
				targetStillValid := false
				for idx := range g.Airplanes {
					if p.Target == g.Airplanes[idx] {
						targetStillValid = true
						break
					}
				}

				if p.Target != nil && targetStillValid { // If target is still valid and set
					dx := p.Target.Position.X - p.Position.X
					dy := p.Target.Position.Y - p.Position.Y
					dist := math.Sqrt(dx*dx + dy*dy)

					if dist < p.Speed { // Close enough to hit/overlap in this step
						// Move directly to target's last known position if very close
						p.Position.X = p.Target.Position.X
						if p.Position.Y > p.Target.Position.Y {
							p.Position.Y = p.Target.Position.Y
						}
					} else if dist > 0 { // Move towards target
						p.Position.X += (dx / dist) * p.Speed
						if p.Position.Y > p.Target.Position.Y {
							p.Position.Y += (dy / dist) * p.Speed
						}
					}
					// Update rotation to point towards target
					p.Rotation = math.Atan2(dy, dx) + math.Pi/2.0
				} else {
					p.Position.Y -= p.Speed
					if p.Rotation < 0 {
						p.Position.X -= p.Speed
					} else {
						p.Position.X += p.Speed
					}
				}
			}

			if p.Position.Y > 0 || (p.Position.X > 0 && p.Position.X < WindowWidth) {
				g.Airdefense.Projectiles = append(g.Airdefense.Projectiles, *p)
			}
		}
	}

	// Spawn new airplanes
	if g.Counter >= g.NextAirplaneSpawnTick {
		g.spawnAirplane()
		g.NextAirplaneSpawnTick = g.Counter + MinAirplaneSpawnCooldownTicks + rand.Intn(RandomAirplaneSpawnCooldownRange)
	}

	// Update airplanes (movement and removal if off-screen)
	keptAirplanes := []*Airplane{}
	for i := range g.Airplanes {
		plane := g.Airplanes[i]
		plane.Position.X += plane.Speed

		// Check if airplane is off-screen
		onScreen := true
		if plane.Speed > 0 { // Moving right
			if plane.Position.X > WindowWidth {
				onScreen = false
			}
		} else { // Moving left
			if plane.Position.X < -float64(plane.Image.Bounds().Dx()) {
				onScreen = false
			}
		}
		if onScreen {
			keptAirplanes = append(keptAirplanes, plane)
		}
	}
	g.Airplanes = keptAirplanes

	// Collision detection: Missiles vs Airplanes
	survivingProjectiles := []Projectile{}
	hitAirplaneIndices := make(map[int]struct{}) // Store indices of airplanes hit this frame
	for i := range g.Airdefense.Projectiles {
		p := &g.Airdefense.Projectiles[i]
		collidedThisFrame := false

		if p.Position != nil {
			for j := range g.Airplanes {
				if _, alreadyHit := hitAirplaneIndices[j]; alreadyHit {
					continue
				}
				plane := g.Airplanes[j]
				if checkCollision(*p, *plane) {
					hitAirplaneIndices[j] = struct{}{} // Mark airplane as hit
					// Create explosion at airplane's position
					newExplosion := Explosion{
						Position:           plane.Position, // Ou p.Position, dependendo do efeito desejado
						Frames:             explosionFrames,
						CurrentFrameIndex:  0,
						FrameDurationTicks: ExplosionFrameDurationTicks,
						CurrentTickInFrame: 0,
					}
					g.Explosions = append(g.Explosions, newExplosion)

					g.AirplanesDestroyed++
					if g.AudioContext != nil && len(explosion) > 0 {
						// Decodifica os dados MP3 para PCM antes de reproduzir
						explosionStream, err := mp3.DecodeWithSampleRate(g.AudioContext.SampleRate(), bytes.NewReader(explosion))
						if err == nil { // Verificação básica de erro
							explosionPlayer, _ := g.AudioContext.NewPlayer(explosionStream)
							explosionPlayer.Play()
						}
						// Considere registrar o erro se err != nil
					}
					collidedThisFrame = true
					break
				}
			}
		}
		if !collidedThisFrame {
			survivingProjectiles = append(survivingProjectiles, *p)
		}
	}
	g.Airdefense.Projectiles = survivingProjectiles

	newAirplanesList := []*Airplane{}
	for i := range g.Airplanes {
		if _, wasHit := hitAirplaneIndices[i]; !wasHit {
			newAirplanesList = append(newAirplanesList, g.Airplanes[i])
		}
	}
	g.Airplanes = newAirplanesList

	// Update explosions
	activeExplosions := []Explosion{}
	for i := range g.Explosions {
		expl := &g.Explosions[i]
		expl.CurrentTickInFrame++
		if expl.CurrentTickInFrame >= expl.FrameDurationTicks {
			expl.CurrentTickInFrame = 0
			expl.CurrentFrameIndex++
		}
		if expl.CurrentFrameIndex < len(expl.Frames) {
			activeExplosions = append(activeExplosions, *expl)
		}
	}
	g.Explosions = activeExplosions

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.Airdefense.draw(screen)

	for _, plane := range g.Airplanes {
		plane.Draw(screen)
	}

	// Draw explosions
	for i := range g.Explosions {
		expl := &g.Explosions[i]
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(expl.Position.X, expl.Position.Y)
		screen.DrawImage(expl.Frames[expl.CurrentFrameIndex], opts)
	}
	drawGround(screen)

	// Display destroyed airplanes count
	ebitenutil.DebugPrint(screen, fmt.Sprintf("Destroyed: %d", g.AirplanesDestroyed))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func drawGround(screen *ebiten.Image) {
	groundWidth := ground.Bounds().Dx()

	y := float64(WindowHeight - ground.Bounds().Dy())
	for i := range (WindowWidth / groundWidth) + 1 {
		x := float64(i * groundWidth)

		opt := &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(x, y)

		screen.DrawImage(ground, opt)
	}
}

func checkCollision(p Projectile, a Airplane) bool {
	if p.Image == nil || p.Position == nil || a.Image == nil {
		return false
	}

	projRect := image.Rect(
		int(p.Position.X),
		int(p.Position.Y),
		int(p.Position.X)+p.Image.Bounds().Dx(),
		int(p.Position.Y)+p.Image.Bounds().Dy(),
	)
	planeRect := image.Rect(
		int(a.Position.X),
		int(a.Position.Y),
		int(a.Position.X)+(a.Image.Bounds().Dx()/2),
		int(a.Position.Y)+(a.Image.Bounds().Dy()/2),
	)
	return projRect.Overlaps(planeRect)
}
