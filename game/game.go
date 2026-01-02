package game

import (
	"image/color"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"gametry/assets"
)

const WindowWidthPixels = 1600
const WindowHeightPixels = 1200

type Window struct {
	Width, Height int
}

func Halves(sprite *ebiten.Image) (w, h float64) {
	width := sprite.Bounds().Dx()
	height := sprite.Bounds().Dy()

	return float64(width / 2), float64(height / 2)
}

// =================================================================================
// ================================== Game =========================================
// =================================================================================
type Game struct {
	Window           Window
	AudioContext     *audio.Context
	Player           Player
	Missle           []*Missle
	MeteorSpawnTimer *Timer
	Meteor           []*Meteor
}

func NewGame() *Game {
	playerCanon := NewSimpleCanon(assets.CanonSprite)

	player := NewPlayer(
		Vector{WindowWidthPixels / 2, WindowHeightPixels / 2},
		assets.PlayerSprite,
		playerCanon,
	)

	g := &Game{
		Window:           Window{Width: WindowWidthPixels, Height: WindowHeightPixels},
		Player:           player,
		MeteorSpawnTimer: NewTimer(900*time.Millisecond + time.Millisecond*time.Duration(rand.Intn(100))),
	}

	return g
}

func (g *Game) AddMissle(m *Missle) {
	g.Missle = append(g.Missle, m)
}

func ExcludeIndexFuckOrder[T any](s []T, i int) ([]T, int) {
	j := len(s) - 1
	if i == j {
		return s[:i], i
	}
	s[i], s[j] = s[j], s[i]
	s = s[:j]
	if i > 0 {
		i--
	}
	return s, i
}

func (g *Game) Update() error {
	if g.AudioContext == nil {
		g.AudioContext = audio.NewContext(assets.SampleRate)
	}
	if err := g.Player.Update(g); err != nil {
		return err
	}

	g.SpawnMeteors()
	g.UpdateMeteors()
	g.UpdateMissles()
	g.UpdateCollisions()
	g.RemoveDistantMeteors()

	return nil
}

func (g *Game) SpawnMeteors() {
	g.MeteorSpawnTimer.Update()
	if g.MeteorSpawnTimer.IsReady() {
		g.MeteorSpawnTimer.Reset()

		g.SpawnMeteor()
	}
}

func (g *Game) SpawnMeteor() {
	sprite := assets.MeteorSprites[rand.Intn(len(assets.MeteorSprites))]
	pos := Vector{
		X: float64(rand.Intn(g.Window.Width)),
		Y: (float64(sprite.Bounds().Dx()) / 2),
	}
	velocity := float64(g.Window.Height/ebiten.TPS()) / 5
	spin := (math.Pi * (rand.Float64() - 0.5) * 1.5) / float64(ebiten.TPS())
	angle := math.Pi + (rand.Float64()-0.5)*math.Pi/7
	m := NewMeteor(pos, angle, velocity, spin, sprite)
	g.Meteor = append(g.Meteor, m)
}

func (g *Game) UpdateMeteors() {
	for i := 0; i < len(g.Meteor); i++ {
		g.Meteor[i].Update()
	}
}

func (g *Game) UpdateMissles() {
	for i := 0; i < len(g.Missle); i++ {
		if keep := g.Missle[i].Update(g); !keep {
			g.Missle, i = ExcludeIndexFuckOrder(g.Missle, i)
		}
	}
}

func (g *Game) RemoveDistantMeteors() {
	for i := 0; i < len(g.Meteor); i++ {
		if g.Meteor[i].IsMeteorFarAway(g.Window) {
			g.Meteor, i = ExcludeIndexFuckOrder(g.Meteor, i)
		}
	}
}

func (g *Game) UpdateCollisions() {
	for i := 0; i < len(g.Missle); i++ {
		for j := 0; i > -1 && i < len(g.Missle) && j < len(g.Meteor); j++ {
			m := g.Meteor[j]
			if g.Missle[i].IntersectsCircle(m.Position, m.Radius()) {
				log.Printf("HIT! Missle: %v Meteor: %v", i, j)
				g.Missle, i = ExcludeIndexFuckOrder(g.Missle, i)
				g.Meteor, j = ExcludeIndexFuckOrder(g.Meteor, j)
				g.AudioContext.NewPlayerFromBytes(assets.MeteorExplodeBytes).Play()
			}
		}
	}
	for i := 0; i < len(g.Meteor); i++ {
		m := g.Meteor[i]
		if g.Player.IntersectsCircle(m.Position, m.Radius()) {
			log.Printf("HIT PLAYER Meteor: %v", i)
			g.Meteor, i = ExcludeIndexFuckOrder(g.Meteor, i)
			g.Player.Hit(g.AudioContext)
		}
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.Player.Draw(screen)
	for _, m := range g.Missle {
		m.Draw(screen)
	}
	for _, m := range g.Meteor {
		m.Draw(screen)
	}
	g.DrawBorder(screen)
}

func (g *Game) DrawBorder(screen *ebiten.Image) {
	borderColor := &color.RGBA{G: 255}
	w, h := float32(g.Window.Width), float32(g.Window.Height)
	vector.StrokeLine(screen, 0, 0, w, 0, 2.0, borderColor, false)
	vector.StrokeLine(screen, w, 0, w, h, 2.0, borderColor, false)
	vector.StrokeLine(screen, 0, h, w, h, 2.0, borderColor, false)
	vector.StrokeLine(screen, 0, 0, 0, h, 2.0, borderColor, false)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.Window.Width, g.Window.Height
}

// ================================ Game done ======================================
