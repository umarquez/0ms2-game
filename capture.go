package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	log "github.com/sirupsen/logrus"
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const captureInterval = 10000

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

type Capture struct {
	gif          *gif.GIF
	accumulator  int64
	frameCounter int
	wg           *sync.WaitGroup
	record       bool
}

var captureInstance *Capture

func GetCaptureInstance() *Capture {
	if captureInstance == nil {
		captureInstance = new(Capture)
		captureInstance.wg = new(sync.WaitGroup)
		captureInstance.record = true
	}

	return captureInstance
}

func (capture *Capture) palette() color.Palette {
	/*if 1 < capture.frameCounter/25 {
		return cheapPalette
	}*/
	return palette.Plan9
}

func (capture *Capture) Capture(screen *ebiten.Image, delta int64) {
	if !capture.record {
		return
	}
	capture.accumulator += delta
	if capture.accumulator < captureInterval {
		return
	}
	capture.accumulator = 0

	if capture.gif == nil {
		capture.gif = &gif.GIF{
			Image:     []*image.Paletted{},
			Delay:     []int{},
			LoopCount: 0,
		}
	}

	capture.frameCounter++

	src := image.NewRGBA(screen.Bounds())
	draw.Draw(src, screen.Bounds(), screen, image.Point{}, draw.Src)

	capture.wg.Add(1)
	go func(src image.Image) {
		img := image.NewPaletted(src.Bounds(), capture.palette())
		draw.FloydSteinberg.Draw(img, img.Bounds(), src, src.Bounds().Min)
		capture.gif.Image = append(capture.gif.Image, img)
		capture.gif.Delay = append(capture.gif.Delay, 1)
		capture.wg.Done()
	}(src)
}

func (capture *Capture) WriteAndClose() {
	if !capture.record {
		return
	}

	capture.record = false
	capture.wg.Wait()
	fname := filepath.Join(capturesDir, fmt.Sprintf("%v.gif", time.Now().Unix()))
	fOut, err := os.Create(fname)
	if err != nil {
		log.Error(err)
		return
	}

	err = gif.EncodeAll(fOut, capture.gif)
	if err != nil {
		log.Error(err)
		return
	}

	err = fOut.Close()
}
