package game

import (
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Vector struct {
	X, Y float64
}

func (v Vector) PivotRotate(p Vector, d Vector) Vector {
	sh := Vector{X: v.X - p.X, Y: v.Y - p.Y}
	return Vector{
		X: p.X + sh.X*d.Y - sh.Y*d.X,
		Y: p.Y + sh.X*d.X + sh.Y*d.Y,
	}
}

func (v Vector) Minus(a Vector) Vector {
	return Vector{X: v.X - a.X, Y: v.Y - a.Y}
}

func (v Vector) OrtogonalLeft() Vector {
	return Vector{X: -v.Y, Y: v.X}
}

func (v Vector) DotPrduct(a Vector) float64 {
	return v.X*a.X + v.Y*a.Y
}

func (v Vector) Magnitude() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}
func (v Vector) Normalized() Vector {
	magnitude := v.Magnitude()
	return Vector{X: v.X / magnitude, Y: v.Y / magnitude}
}

type Box struct {
	Vertex []Vector
	Center Vector
}

func (b Box) DrawBorder(screen *ebiten.Image) {
	borderColor := &color.RGBA{G: 255}
	v := b.Vertex
	for i := 0; i < len(b.Vertex); i++ {
		vector.StrokeLine(
			screen,
			float32(v[i].X),
			float32(v[i].Y),
			float32(v[(i+1)%len(v)].X),
			float32(v[(i+1)%len(v)].Y),
			1.0,
			borderColor,
			false,
		)
	}
}

func (b Box) Rotate(direction Vector) {
	for i := 0; i < len(b.Vertex); i++ {
		b.Vertex[i] = b.Vertex[i].PivotRotate(b.Center, direction)
	}
}

func (b Box) IntersectsCircle(c Vector, r float64) bool {
	p := c.Minus(b.Center)
	pNorm := p.Normalized()
	var maxProjection float64
	for i := 0; i < len(b.Vertex); i++ {
		v := b.Center.Minus(b.Vertex[i])
		proj := v.DotPrduct(pNorm)
		if i == 0 || maxProjection < proj {
			maxProjection = proj
		}
	}
	axisMagnitude := p.Magnitude()
	if axisMagnitude <= 0 {
		log.Printf("axis magnitude %v", axisMagnitude)
	}
	if axisMagnitude > 0 && axisMagnitude-r-maxProjection > 0 {
		return false
	}
	return true
}
