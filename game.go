package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/ungerik/go3d/float64/vec2"
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"time"
)

type GameEntities interface {
	Draw(*ebiten.Image)
	Update(*ebiten.Image, int64)
}

type Game struct {
	ShowFPS      bool
	lastUpdate   time.Time
	bgColor      color.Color
	gameSize     image.Point
	entities     []GameEntities
	frameCounter int
	gif          *gif.GIF
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
	game.entities = append(game.entities,
		starfield,
		planets,
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

var cheapPalette color.Palette

func init() {
	cs := []color.Color{}
	for _, r := range []uint8{0x00, 0x80, 0xff} {
		for _, g := range []uint8{0x00, 0x80, 0xff} {
			for _, b := range []uint8{0x00, 0x80, 0xff} {
				cs = append(cs, color.RGBA{r, g, b, 0xff})
			}
		}
	}
	cheapPalette = color.Palette(cs)
}

func (g *Game) palette() color.Palette {
	if 1 < g.frameCounter/25 {
		return cheapPalette
	}
	return palette.Plan9
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(g.bgColor)

	for _, e := range g.entities {
		e.Draw(screen)
	}

	if g.ShowFPS {
		ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f", ebiten.CurrentTPS(), ebiten.CurrentFPS()))
	}

	if g.frameCounter%25 == 0 {
		if g.gif == nil {
			num := g.frameCounter / 25
			g.gif = &gif.GIF{
				Image:     make([]*image.Paletted, num),
				Delay:     make([]int, num),
				LoopCount: -1,
			}
		}
		g.frameCounter = (g.frameCounter + 25) % 25
		capture := image.NewRGBA(screen.Bounds())

		draw.Draw(capture, screen.Bounds(), screen, image.Point{}, draw.Src)
		go func() {
			img := image.NewPaletted(capture.Bounds(), g.palette())
			draw.FloydSteinberg.Draw(img, img.Bounds(), capture, capture.Bounds().Min)
			g.gif.Image = append(g.gif.Image, img)
			g.gif.Delay = append(g.gif.Delay, g.frameCounter/25)
		}()
	}

	g.frameCounter++
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return windowWidth, windowHeight
}
