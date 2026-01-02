package game

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Meteor struct {
	Position  Vector        // Where it is
	Direction Vector        // Where go next
	Velocity  float64       // Speed
	Rotation  float64       // Current angle
	Spin      float64       // Angular velocity
	Sprite    *ebiten.Image // Personal look
	Scale     float64
}

func NewMeteor(
	pos Vector,
	angle float64,
	velocity float64,
	spin float64,
	sprite *ebiten.Image,
) *Meteor {
	m := &Meteor{
		Position:  pos,
		Velocity:  velocity,
		Direction: Vector{X: math.Sin(angle), Y: math.Cos(angle)},
		Spin:      spin,
		Sprite:    sprite,
		Scale:     0.5,
	}
	//log.Printf("New meteor data: velocity: %v; angle: %v; dir: %+v; spin: %v", velocity, angle, m.Direction, spin)
	return m
}

func (m *Meteor) Update() {
	//speed := float64(WindowHeightPixels/ebiten.TPS()) / 3
	m.Rotation += m.Spin
	if m.Rotation > 2*math.Pi {
		m.Rotation = 2*math.Pi - m.Rotation
	}

	m.Position.X += m.Velocity * m.Direction.X
	m.Position.Y -= m.Velocity * m.Direction.Y
}

func (m *Meteor) Radius() float64 { return m.Scale * float64(m.Sprite.Bounds().Dx()) / 2 }

func (m Meteor) Draw(screen *ebiten.Image) {
	pivotX, pivotY := Halves(m.Sprite)

	op := &ebiten.DrawImageOptions{}
	// Rotation
	op.GeoM.Translate(-pivotX, -pivotY)
	op.GeoM.Rotate(m.Rotation)
	op.GeoM.Translate(pivotX, pivotY)
	op.GeoM.Scale(m.Scale, m.Scale)
	// Position
	op.GeoM.Translate(m.Position.X-pivotX/2, m.Position.Y-pivotY/2)

	screen.DrawImage(m.Sprite, op)
}

func (m *Meteor) IsMeteorFarAway(window Window) bool {
	r := m.Radius()

	leftXLimit := -float64(window.Width) - r
	rightXLimit := 2*float64(window.Width) + r
	topYLimit := -float64(window.Height) - r
	bottomYLimit := 2*float64(window.Height) + r

	if m.Position.X < leftXLimit {
		return true
	}
	if m.Position.X > rightXLimit {
		return true
	}
	if m.Position.Y < topYLimit {
		return true
	}
	if m.Position.Y > bottomYLimit {
		return true
	}
	return false
}
