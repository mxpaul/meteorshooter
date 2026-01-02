package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"gametry/game"
)

func main() {
	g := game.NewGame()

	ebiten.SetWindowTitle("Игруля")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	err := ebiten.RunGame(g)
	if err != nil {
		log.Fatalf("RunGame error: %v", err)
	}
	log.Printf("missle count: %d", len(g.Missle))
}
