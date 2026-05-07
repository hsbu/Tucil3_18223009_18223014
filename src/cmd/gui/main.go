package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenW = 800
	screenH = 600
)

type Screen interface {
	Update() (Screen, error)
	Draw(screen *ebiten.Image)
}

type Game struct {
	current Screen
}

func (g *Game) Update() error {
	next, err := g.current.Update()
	if err != nil {
		return err
	}
	g.current = next
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.current.Draw(screen)
}

func (g *Game) Layout(_, _ int) (int, int) {
	return screenW, screenH
}

func main() {
	ebiten.SetWindowSize(screenW, screenH)
	ebiten.SetWindowTitle("Ice Sliding Puzzle Solver")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	g := &Game{current: newFileSelectScreen()}
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
