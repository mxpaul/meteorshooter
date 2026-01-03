package game

import (
	"math"

	"github.com/mxpaul/meteorshooter/assets"

	"github.com/hajimehoshi/ebiten/v2"
)

type Missle struct {
	Position  Vector
	Direction Vector
	Rotation  float64
	Sprite    *ebiten.Image
}

func NewMissle(pos Vector, angle float64, distance float64) *Missle {
	m := &Missle{
		Position: Vector{
			pos.X + math.Sin(angle)*distance,
			pos.Y - math.Cos(angle)*distance,
		},
		Direction: Vector{
			math.Sin(angle),
			math.Cos(angle),
		},
		Rotation: angle,

		Sprite: assets.MissleSprite,
	}
	return m
}

func (m *Missle) Update(g *Game) (keep bool) {
	speed := float64(WindowHeightPixels/ebiten.TPS()) / 5 // 1.5

	m.Position.X += speed * m.Direction.X
	m.Position.Y -= speed * m.Direction.Y

	return m.IsMissleInWindow(g.Window)
}

func (m *Missle) IsMissleInWindow(window Window) bool {
	x, y := float64(m.Sprite.Bounds().Dx()), float64(m.Sprite.Bounds().Dy())
	h := math.Sqrt(x*x + y*y)

	leftXLimit := -h
	rightXLimit := float64(window.Width) + h
	topYLimit := -h
	bottomYLimit := float64(window.Height) - h

	if m.Position.X < leftXLimit {
		return false
	}
	if m.Position.X > rightXLimit {
		return false
	}
	if m.Position.Y < topYLimit {
		return false
	}
	if m.Position.Y > bottomYLimit {
		return false
	}
	return true
}

func (m Missle) PivotX() float64 { return float64(m.Sprite.Bounds().Dx()) / 2 }

func (m Missle) PivotY() float64 { return float64(m.Sprite.Bounds().Dy()) }

func (m Missle) Draw(screen *ebiten.Image) {
	pivotX, pivotY := m.PivotX(), m.PivotY()

	op := &ebiten.DrawImageOptions{}
	// Canon rotation
	op.GeoM.Translate(-pivotX, -pivotY)
	op.GeoM.Rotate(m.Rotation)
	op.GeoM.Translate(pivotX, pivotY)
	// Canon position
	op.GeoM.Translate(m.Position.X-pivotX, m.Position.Y-pivotY)

	screen.DrawImage(m.Sprite, op)
	// m.Box().DrawBorder(screen)
}

func (m Missle) Box() Box {
	pivotX, pivotY := m.PivotX(), m.PivotY()
	r := Box{
		Center: m.Position,
		Vertex: []Vector{
			{m.Position.X - pivotX, m.Position.Y - pivotY},
			{m.Position.X + pivotX, m.Position.Y - pivotY},
			{m.Position.X + float64(m.Sprite.Bounds().Dx()) - pivotX, m.Position.Y + float64(m.Sprite.Bounds().Dy()) - pivotY},
			{m.Position.X - float64(m.Sprite.Bounds().Dx()) + pivotX, m.Position.Y + float64(m.Sprite.Bounds().Dy()) - pivotY},
		}}
	r.Rotate(m.Direction)
	return r
}

func (m Missle) IntersectsCircle(c Vector, r float64) bool {
	return m.Box().IntersectsCircle(c, r)
}
