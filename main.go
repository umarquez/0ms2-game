package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	_ "github.com/silbinarywolf/preferdiscretegpu"
	log "github.com/sirupsen/logrus"
	"image"
	"image/color"
	"image/gif"
	"math/rand"
	"os"
	"path/filepath"
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
		fname := filepath.Join("capture", fmt.Sprintf("%v.gif", time.Now().Unix()))
		fOut, _ := os.Create(fname)
		_ = gif.EncodeAll(fOut, game.gif)
		_ = fOut.Close()
	}()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
