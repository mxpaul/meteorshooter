package game

import (
	"fmt"
	"math"

	"github.com/mxpaul/meteorshooter/assets"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/colorm"
)

type Player struct {
	Position  Vector
	Sprite    *ebiten.Image
	Speed     float64
	Canon     *CanonSimple
	InHit     bool
	translate float64
	blinkRate float64
	blinkUp   bool
}

func NewPlayer(
	initialPos Vector,
	sprite *ebiten.Image,
	canon *CanonSimple,
) Player {
	p := Player{
		Position: initialPos,
		Sprite:   sprite,
		Speed:    float64(WindowHeightPixels/ebiten.TPS()) / 2,
		Canon:    canon,
	}
	return p
}

func (p *Player) UpdatePosition(g *Game) error {
	var delta Vector

	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		delta.Y = p.Speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		delta.Y = -p.Speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		delta.X = -p.Speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		delta.X = p.Speed
	}

	// Check for diagonal movement
	if delta.X != 0 && delta.Y != 0 {
		factor := p.Speed / math.Sqrt(delta.X*delta.X+delta.Y*delta.Y)
		delta.X *= factor
		delta.Y *= factor
	}

	p.Position.X += delta.X
	p.Position.Y += delta.Y

	p.LimitPositionToWindow(g.Window)
	return nil
}

func (p *Player) LimitPositionToWindow(window Window) {
	halfW, halfH := Halves(p.Sprite)
	leftXLimit := halfW
	rightXLimit := float64(window.Width) - halfW
	topYLimit := halfH
	bottomYLimit := float64(window.Height) - halfH

	if p.Position.X < leftXLimit {
		p.Position.X = leftXLimit
	}
	if p.Position.X > rightXLimit {
		p.Position.X = rightXLimit
	}
	if p.Position.Y < topYLimit {
		p.Position.Y = topYLimit
	}
	if p.Position.Y > bottomYLimit {
		p.Position.Y = bottomYLimit
	}
}

func (p *Player) Update(g *Game) error {
	if p.InHit {
		if p.blinkUp {
			p.translate += p.blinkRate
		} else {
			p.translate -= p.blinkRate
		}
		if p.translate > 1.0 {
			p.translate = 1.0
			p.blinkUp = false
		}
		if p.translate < 0 {
			p.translate = 0
			p.InHit = false
		}
		return nil
	}
	if err := p.UpdatePosition(g); err != nil {
		return fmt.Errorf("player update position failed: %w", err)
	}
	if err := p.Canon.Update(g, p.Position); err != nil {
		return fmt.Errorf("player canon update failed: %w", err)
	}
	return nil
}

func (p Player) Draw(screen *ebiten.Image) {
	halfW, halfH := Halves(p.Sprite)

	op := &colorm.DrawImageOptions{}
	op.GeoM.Translate(p.Position.X-halfW, p.Position.Y-halfH)

	cm := colorm.ColorM{}
	cm.Translate(p.translate, p.translate, p.translate, 0.0)

	colorm.DrawImage(screen, p.Sprite, cm, op)

	p.Canon.Draw(screen, cm)
	//p.Box().DrawBorder(screen)
}

func (p Player) Box() Box {
	halfW, halfH := Halves(p.Sprite)
	return Box{
		Center: p.Position,
		Vertex: []Vector{
			{p.Position.X - halfW, p.Position.Y - halfH},
			{p.Position.X + halfW, p.Position.Y - halfH},
			{p.Position.X + halfW, p.Position.Y + halfH},
			{p.Position.X - halfW, p.Position.Y + halfH},
		}}
}

func (p Player) IntersectsCircle(c Vector, r float64) bool {
	return p.Box().IntersectsCircle(c, r)
}

func (p *Player) Hit(audioContext *audio.Context) {
	p.InHit = true
	p.translate = 0.0
	p.blinkRate = 2.5 / float64(ebiten.TPS())
	p.blinkUp = true
	audioContext.NewPlayerFromBytes(assets.PlayerHitBytes).Play()
}
