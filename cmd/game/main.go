package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/mxpaul/meteorshooter/game"
)

func main() {
	g := game.NewGame()

	ebiten.SetWindowTitle("Meteor shooter")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	err := ebiten.RunGame(g)
	if err != nil {
		log.Fatalf("RunGame error: %v", err)
	}
	log.Printf("missle count: %d; meteor count: %d", len(g.Missle), len(g.Meteor))
}
