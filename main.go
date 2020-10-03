package main

import "C"
import (
	"github.com/hajimehoshi/ebiten"
	_ "github.com/silbinarywolf/preferdiscretegpu"
	log "github.com/sirupsen/logrus"
	"image"
	"image/color"
	"math/rand"
	"time"
)

var game *Game

func init() {
	// init random seed
	rand.Seed(time.Now().UnixNano())
	game = newGame(color.Black, image.Point{windowWidth, windowHeight})
	game.ShowFPS = true
	log.SetLevel(log.DebugLevel)
}

func main() {
	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle(gameTitle)

	defer func() {
		GetCaptureInstance().WriteAndClose()
	}()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
