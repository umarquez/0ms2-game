package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/ungerik/go3d/float64/vec2"
	"image"
	"image/color"
	"time"
)

type GameEntities interface {
	Draw(*ebiten.Image)
	Update(*ebiten.Image, int64)
}

type Game struct {
	ShowFPS    bool
	lastUpdate time.Time
	bgColor    color.Color
	gameSize   image.Point
	entities   []GameEntities
	lastDraw   time.Time
}

func newGame(bg color.Color, windowSize image.Point) *Game {
	game := &Game{
		bgColor:  bg,
		gameSize: windowSize,
	}

	playerPosition := vec2.T{
		(windowWidth - (playerSize * j0hnScale)) / 2,
		(windowHeight - (playerSize * j0hnScale)) - 20,
	}

	playerPosition.Scale(1 / j0hnScale)

	player := NewJ0hn().SetPosition(playerPosition)
	starfield := NewStarfield(player)
	planets := NewPlanetSpawner(player)
	powerups := NewPowerupSpawner(player)
	game.entities = append(game.entities,
		starfield,
		planets,
		powerups,
		player,
	)

	return game
}

func (g *Game) Update(screen *ebiten.Image) error {
	d := time.Since(g.lastUpdate).Milliseconds()

	for _, e := range g.entities {
		e.Update(screen, d)
	}

	g.lastUpdate = time.Now()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(g.bgColor)

	for _, e := range g.entities {
		e.Draw(screen)
	}

	if g.ShowFPS {
		ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f", ebiten.CurrentTPS(), ebiten.CurrentFPS()))
	}

	GetCaptureInstance().Capture(screen, time.Since(g.lastDraw).Milliseconds())
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return windowWidth, windowHeight
}
