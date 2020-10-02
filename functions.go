package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	log "github.com/sirupsen/logrus"
	"github.com/ungerik/go3d/float64/vec2"
)

func Panic(source string, params map[string]interface{}, err error) {
	log.WithFields(log.Fields{
		"source": source,
		"params": params,
	}).Panic(err)
}

func loadSprite(spritePath string) *ebiten.Image {
	ebImg, _, err := ebitenutil.NewImageFromFile(spritePath, ebiten.FilterNearest)
	if err != nil {
		Panic("loadSprite", map[string]interface{}{"spritePath": spritePath}, err)
	}

	return ebImg
}

func copyVector(src vec2.T) vec2.T {
	var dest [2]float64
	dest[0] = src[0]
	dest[1] = src[1]

	return vec2.T(dest)
}
