package main

import (
	"fmt"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
	log "github.com/sirupsen/logrus"
	"github.com/ungerik/go3d/float64/vec2"
	"golang.org/x/image/colornames"
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"math"
	"path/filepath"
)

const uiSprite = "ui-level-bar.png"
const uiRedBarSprite = "red-level.png"
const uiBlueBarSprite = "blue-level.png"
const fontfile = "ka1.ttf"

const barW = 32

//const barH = 120
const barMargin = 7
const uiMarginLeft = 10

var imgBar *ebiten.Image
var imgO2Level *ebiten.Image
var imgFuelLevel *ebiten.Image

func init() {
	imgBar = loadSprite(filepath.Join(spritesPath, uiSprite))
	imgO2Level = loadSprite(filepath.Join(spritesPath, uiBlueBarSprite))
	imgFuelLevel = loadSprite(filepath.Join(spritesPath, uiRedBarSprite))
}

type UserInterface struct {
	player           *J0hn
	o2Position       vec2.T
	fuelPosition     vec2.T
	distancePosition vec2.T
	o2Level          int
	o2Offset         float64
	fuelLevel        int
	fuelOffset       float64
	distanceLevel    int
	uiScale          float64
	src              image.Image
	font             *truetype.Font
	ctxFont          *freetype.Context
}

func NewUi(player *J0hn) *UserInterface {
	ui := new(UserInterface)
	ui.uiScale = 3
	ui.o2Position = vec2.T{
		uiMarginLeft,
		windowHeight / ui.uiScale,
	}
	ui.fuelPosition = vec2.T{
		uiMarginLeft + (barW * ui.uiScale),
		windowHeight / ui.uiScale,
	}
	ui.player = player
	ui.src = image.NewUniform(colornames.Green)
	fontBytes, err := ioutil.ReadFile(fontfile)
	if err != nil {
		log.Println(err)
	}

	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Println(err)
	}
	ui.font = f

	imgClip := image.NewRGBA(image.Rect(0, 0, 500, 200))
	draw.Draw(imgClip, imgClip.Bounds(), image.NewUniform(colornames.Green), image.Point{}, draw.Src)

	ctxFont := freetype.NewContext()
	ctxFont.SetDPI(72)
	ctxFont.SetFontSize(20)
	ctxFont.SetClip(imgClip.Bounds())
	ctxFont.SetSrc(imgClip)
	ctxFont.SetFont(ui.font)
	ui.ctxFont = ctxFont

	return ui
}

func (ui *UserInterface) Draw(screen *ebiten.Image) {
	//if !ui.player.isLifting && ui.player.flying {}
	ui.ctxFont.SetDst(screen)
	optO2 := &ebiten.DrawImageOptions{}
	optO2.GeoM.Translate(ui.o2Position[0]/ui.uiScale, ui.o2Position[1]/ui.uiScale)
	optO2.GeoM.Scale(ui.uiScale, ui.uiScale)

	optO2Fill := &ebiten.DrawImageOptions{}
	optO2Fill.GeoM.Translate(ui.o2Position[0]/ui.uiScale, (ui.o2Position[1]+ui.o2Offset)/ui.uiScale)
	optO2Fill.GeoM.Scale(ui.uiScale, ui.uiScale)
	_ = screen.DrawImage(imgO2Level.SubImage(image.Rect(0, 0, barW, ui.o2Level)).(*ebiten.Image), optO2Fill)
	_ = screen.DrawImage(imgBar, optO2)

	optFuel := &ebiten.DrawImageOptions{}
	optFuel.GeoM.Translate(ui.fuelPosition[0]/ui.uiScale, ui.fuelPosition[1]/ui.uiScale)
	optFuel.GeoM.Scale(ui.uiScale, ui.uiScale)

	optFuelFill := &ebiten.DrawImageOptions{}
	optFuelFill.GeoM.Translate(ui.fuelPosition[0]/ui.uiScale, (ui.fuelPosition[1]+ui.fuelOffset)/ui.uiScale)
	optFuelFill.GeoM.Scale(ui.uiScale, ui.uiScale)
	_ = screen.DrawImage(imgFuelLevel.SubImage(image.Rect(0, 0, barW, ui.fuelLevel)).(*ebiten.Image), optFuelFill)
	_ = screen.DrawImage(imgBar, optFuel)

	text.Draw(screen, fmt.Sprintf("%vkm", math.Round(ui.player.relativePosition[1])), truetype.NewFace(ui.font, &truetype.Options{
		Size:              30,
		DPI:               72,
		Hinting:           0,
		GlyphCacheEntries: 0,
		SubPixelsX:        0,
		SubPixelsY:        0,
	}), uiMarginLeft*2, windowHeight-10, colornames.Green)

	text.Draw(screen, "O2", truetype.NewFace(ui.font, &truetype.Options{
		Size:              20,
		DPI:               72,
		Hinting:           0,
		GlyphCacheEntries: 0,
		SubPixelsX:        0,
		SubPixelsY:        0,
	}), int(ui.o2Position[0])+15,
		int(ui.o2Position[1])-4,
		color.RGBA{0x5b, 0x6E, 0xE1, 0xFF},
	)

	text.Draw(screen, "Fuel", truetype.NewFace(ui.font, &truetype.Options{
		Size:              20,
		DPI:               72,
		Hinting:           0,
		GlyphCacheEntries: 0,
		SubPixelsX:        0,
		SubPixelsY:        0,
	}), int(ui.fuelPosition[0]),
		int(ui.fuelPosition[1])-4,
		color.RGBA{0xac, 0x32, 0x32, 0xFF},
	)
}

func (ui *UserInterface) Update(screen *ebiten.Image, delta int64) {
	_, h := imgO2Level.Size()
	ui.o2Level = int(float64(h-14) * ui.player.o2 / 100)
	ui.o2Level += barMargin
	ui.o2Offset = float64(h-ui.o2Level) * ui.uiScale
	ui.o2Offset -= barMargin * ui.uiScale

	ui.fuelLevel = int(float64(h-14) * ui.player.fuel / 100)
	ui.fuelLevel += barMargin
	ui.fuelOffset = float64(h-ui.fuelLevel) * ui.uiScale
	ui.fuelOffset -= barMargin * ui.uiScale
}
