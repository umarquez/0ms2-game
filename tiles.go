package main

import (
	"github.com/hajimehoshi/ebiten"
	log "github.com/sirupsen/logrus"
	"github.com/ungerik/go3d/float64/vec2"
	"github.com/ungerik/go3d/float64/vec3"
	"image"
	"image/color"
	"math/rand"
)

type Tile interface {
	GenericInstance
	Draw(*ebiten.Image)
	Update(vec2.T)
	IsOffscreen() bool
	GetPosition() *vec2.Rect
}

type StarsTile struct {
	*GameInstance
	img      *ebiten.Image
	position *vec2.T
	op       *ebiten.DrawImageOptions
	bounds   *vec2.Rect
	size     *vec2.T
}

func NewStarsTile(position vec2.T) Tile {
	op := ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(position[0]), float64(position[1]))

	max := copyVector(position)

	size := &vec2.T{
		windowWidth,
		windowHeight,
	}
	max.Add(size)

	tile := StarsTile{
		position: &position,
		op:       &op,
		bounds: &vec2.Rect{
			Min: position,
			Max: max,
		},
		size: size,
	}
	tile.GameInstance = NewGenericInstance()

	log.WithFields(map[string]interface{}{
		"position": tile.position,
	}).Debugf("new tile")

	tile.img, _ = ebiten.NewImage(windowWidth, windowHeight, ebiten.FilterNearest)
	for v := float64(windowWidth*windowHeight) * starsProportion; v > 0; v-- {
		tile.img.Set(rand.Intn(windowWidth), rand.Intn(windowHeight), color.White)
	}

	return tile
}

func (t StarsTile) Draw(image *ebiten.Image) {
	_ = image.DrawImage(t.img, t.op)
}

func (t StarsTile) Update(vel vec2.T) {
	t.position.Add(&vel)
	t.op.GeoM.Reset()
	t.op.GeoM.Translate(float64(t.position[0]), float64(t.position[1]))
	t.op.GeoM.Scale(1, 1)

	max := copyVector(*t.position)
	t.bounds = &vec2.Rect{
		Min: copyVector(*t.position),
		Max: *max.Add(t.size),
	}
}

func (t StarsTile) IsOffscreen() bool {
	return t.position[1] > windowHeight
}

func (t StarsTile) GetPosition() *vec2.Rect {
	t.bounds.Min = copyVector(*t.position)
	t.bounds.Max = copyVector(*t.position)
	t.bounds.Max.Add(t.size)
	return t.bounds
}

//var midPoint = 0.46

type InitTile struct {
	*GameInstance
	img      *ebiten.Image
	bounds   *vec2.Rect
	position *vec2.T
	op       *ebiten.DrawImageOptions
	scale    float64
	size     *vec2.T
	//relativeHeight float64
}

func NewInitTile(w, h, scale float64) Tile {
	tile := InitTile{}

	bgImg := image.NewRGBA(image.Rect(0, 0, int(w), int(h)))
	startColor := &vec3.T{203, 219, 255}
	midColor := &vec3.T{99, 155, 255}
	endColor := &vec3.T{0, 0, 0}

	grad := Gradient{}
	grad.AddColor(0, &vec3.T{255, 255, 255})
	grad.AddColor(.3, startColor)
	grad.AddColor(.5, midColor)
	grad.AddColor(.8, midColor)
	grad.AddColor(1, endColor)

	for y := h - 1; y >= 0; y-- {
		p := (h - y) / h

		fillColor := grad.GetColor(p)

		for x := w - 1; x >= 0; x-- {
			bgImg.Set(int(x), int(y), fillColor)
		}
	}

	py := windowHeight - h

	tile.position = &vec2.T{0, py}
	tile.GameInstance = NewGenericInstance()
	tile.op = new(ebiten.DrawImageOptions)
	tile.scale = scale
	tile.img, _ = ebiten.NewImageFromImage(bgImg, ebiten.FilterNearest)
	tile.size = &vec2.T{w, h}
	max := copyVector(*tile.position)
	max.Add(tile.size)
	tile.bounds = &vec2.Rect{
		Min: copyVector(*tile.position),
		Max: max,
	}

	/*
		max := copyVector(origin)

		max.Add(&vec2.T{
			windowWidth,
			windowHeight,
		})

		tile := InitTile{
			position: &position,
			op:       &op,
		}

		log.WithFields(map[string]interface{}{
			"position": tile.position,
		}).Debugf("new tile")

		tile.img, _ = ebiten.NewImage(windowWidth, windowHeight, ebiten.FilterNearest)
		for v := float64(windowWidth*windowHeight) * starsProportion; v > 0; v-- {
			tile.img.Set(rand.Intn(windowWidth), rand.Intn(windowHeight), color.White)
		}*/

	return tile
}

func (t InitTile) Draw(image *ebiten.Image) {
	_ = image.DrawImage(t.img, t.op)
}

func (t InitTile) Update(vel vec2.T) {
	t.position.Add(&vel)
	t.op.GeoM.Reset()
	t.op.GeoM.Translate(float64(t.position[0]), float64(t.position[1]))
	t.op.GeoM.Scale(t.scale, t.scale)
	t.bounds.Min.Add(&vel)
	t.bounds.Max.Add(&vel)
}

func (t InitTile) IsOffscreen() bool {
	return t.position[1] > windowHeight
}

func (t InitTile) GetPosition() *vec2.Rect {
	return t.bounds
}
