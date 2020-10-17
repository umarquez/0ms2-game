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
		bgColor:    bg,
		gameSize:   windowSize,
		lastUpdate: time.Now(),
	}

	playerPosition := vec2.T{
		(windowWidth - (playerSize * j0hnScale)) / 2,
		(windowHeight - (playerSize * j0hnScale)) - 77,
	}

	player := NewJ0hn().SetPosition(playerPosition)
	starfield := NewBackgroundSystem(player)
	planets := NewPlanetSpawner(player)
	powerups := NewPowerupSpawner(player)
	ui := NewUi(player)

	platform := NewPlatform(player)
	platform.SetPosition(&vec2.T{(windowWidth - (platformSize * j0hnScale)) / 2, windowHeight - platformSize*3})

	game.entities = append(game.entities,
		starfield,
		planets,
		powerups,
		platform,
		player,
		ui,
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
	_ = screen.Fill(g.bgColor)

	for _, e := range g.entities {
		e.Draw(screen)
	}

	if g.ShowFPS {
		_ = ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f", ebiten.CurrentTPS(), ebiten.CurrentFPS()))
	}
}

func (g *Game) Layout(int, int) (int, int) {
	return windowWidth, windowHeight
}
