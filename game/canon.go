package game

import (
	"fmt"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"

	"github.com/mxpaul/meteorshooter/assets"
)

type CanonSimple struct {
	Position      Vector
	Rotation      float64
	ShootCooldown *Timer
	Sprite        *ebiten.Image
}

func NewSimpleCanon(
	sprite *ebiten.Image,
) *CanonSimple {
	c := &CanonSimple{
		ShootCooldown: NewReadyTimer(time.Second / 2),
		Sprite:        sprite,
	}
	return c
}

func (c *CanonSimple) Update(g *Game, newPosition Vector) error {
	c.Position = newPosition

	if err := c.HandleRotation(); err != nil {
		return fmt.Errorf("canon rotation error", err)
	}

	c.ShootCooldown.Update()
	if c.ShootCooldown.IsReady() && ebiten.IsKeyPressed(ebiten.KeySpace) {
		c.ShootCooldown.Reset()
		g.AddMissle(NewMissle(c.Position, c.Rotation, c.PivotY()))
		g.AudioContext.NewPlayerFromBytes(assets.CanonShootBytes).Play()
	}

	return nil
}

func (c *CanonSimple) HandleRotation() error {
	speed := 1.2 * math.Pi / float64(ebiten.TPS())

	switch {
	case ebiten.IsKeyPressed(ebiten.KeyW) && ebiten.IsKeyPressed(ebiten.KeyD):
		c.Rotation = math.Pi / 4.0
	case ebiten.IsKeyPressed(ebiten.KeyW) && ebiten.IsKeyPressed(ebiten.KeyA):
		c.Rotation = -math.Pi / 4.0
	case ebiten.IsKeyPressed(ebiten.KeyS) && ebiten.IsKeyPressed(ebiten.KeyD):
		c.Rotation = 3.0 * math.Pi / 4.0
	case ebiten.IsKeyPressed(ebiten.KeyS) && ebiten.IsKeyPressed(ebiten.KeyA):
		c.Rotation = -3.0 * math.Pi / 4.0
	case ebiten.IsKeyPressed(ebiten.KeyW):
		c.Rotation = 0
	case ebiten.IsKeyPressed(ebiten.KeyS):
		c.Rotation = math.Pi
	case ebiten.IsKeyPressed(ebiten.KeyD):
		c.Rotation = math.Pi / 2.0
	case ebiten.IsKeyPressed(ebiten.KeyA):
		c.Rotation = -math.Pi / 2.0
	}

	if ebiten.IsKeyPressed(ebiten.KeyDelete) {
		c.Rotation -= speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyPageDown) {
		c.Rotation += speed
	}
	return nil
}

func (c CanonSimple) PivotX() float64 { return float64(c.Sprite.Bounds().Dx()) / 2 }

func (c CanonSimple) PivotY() float64 { return float64(c.Sprite.Bounds().Dy()) * 2 / 3 }

func (c CanonSimple) Draw(screen *ebiten.Image, cm colorm.ColorM) {
	pivotX, pivotY := c.PivotX(), c.PivotY()
	halfW, halfH := Halves(c.Sprite)

	op := &colorm.DrawImageOptions{}
	// Canon rotation
	op.GeoM.Translate(-pivotX, -pivotY)
	op.GeoM.Rotate(c.Rotation)
	op.GeoM.Translate(pivotX, pivotY)
	// Canon position
	op.GeoM.Translate(c.Position.X-halfW, c.Position.Y-halfH)

	colorm.DrawImage(screen, c.Sprite, cm, op)
}
