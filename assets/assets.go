package assets

import (
	"embed"
	"image"
	_ "image/png"
	"io"
	"io/fs"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
)

//go:embed *
var assets embed.FS

var (
	PlayerSprite       = mustLoadImage("player.png")
	CanonSprite        = mustLoadImage("canon_simple.png")
	MissleSprite       = mustLoadImage("missle1.png")
	MeteorSprites      = mustLoadImages("meteors/*.png")
	CanonShootBytes    = mustLoadOgg("sfx/canon_shoot.ogg")
	PlayerHitBytes     = mustLoadOgg("sfx/player_hit.ogg")
	MeteorExplodeBytes = mustLoadOgg("sfx/meteor_explode.ogg")
	SpaceAmbientWav    = mustLoadFile("music/spaceambient.wav")
)

const SampleRate = 44100

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

func mustLoadImages(path string) []*ebiten.Image {
	matches, err := fs.Glob(assets, path)
	if err != nil {
		panic(err)
	}

	images := make([]*ebiten.Image, len(matches))
	for i, match := range matches {
		images[i] = mustLoadImage(match)
	}

	return images
}

func mustLoadOgg(name string) (b []byte) {
	f, err := assets.Open(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	s, err := vorbis.DecodeWithSampleRate(SampleRate, f)
	if err != nil {
		panic(err)
	}
	b, err = io.ReadAll(s)
	if err != nil {
		panic(err)
	}

	return b
}

func mustLoadFile(name string) (b []byte) {
	f, err := assets.Open(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	b, err = io.ReadAll(f)
	if err != nil {
		panic(err)
	}

	return b
}
