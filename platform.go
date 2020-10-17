package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/ungerik/go3d/float64/vec2"
	"image"
	"path/filepath"
)

var initY float64

const platformFrameInterval = 75

var imgPlatform *ebiten.Image

func init() {
	imgPlatform = loadSprite(filepath.Join(spritesPath, platformSpriteFile))
}

type Platform struct {
	position     vec2.T
	currentFrame int
	accumulator  int64
	player       *J0hn
	initPos      float64
}

func NewPlatform(player *J0hn) *Platform {
	return &Platform{
		player:  player,
		initPos: player.position[1],
	}
}

func (p *Platform) SetPosition(position *vec2.T) {
	p.position = *position.Scale(1 / j0hnScale)
}

func (p *Platform) Draw(screen *ebiten.Image) {
	if p.position[1] > windowHeight/j0hnScale {
		return
	}
	op := ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(p.position[0]), float64(p.position[1]))
	op.GeoM.Scale(j0hnScale, j0hnScale)

	x1, y1 := p.currentFrame*platformSize, 0
	x2, y2 := x1+playerSize, y1+playerSize
	_ = screen.DrawImage(imgPlatform.SubImage(image.Rect(x1, y1, x2, y2)).(*ebiten.Image), &op)
}

func (p *Platform) Update(_ *ebiten.Image, delta int64) {
	if p.position[1] == p.player.upPosition[1] {
		return
	}

	p.accumulator += delta
	w, _ := imgPlatform.Size()
	platformTotalFrames := w / platformSize

	if p.accumulator >= platformFrameInterval {
		p.accumulator = 0

		if (p.player.isLifting || p.player.flying) && p.player.position[1] < p.player.upPosition[1] {
			p.player.position[1]++
		}

		if p.currentFrame != platformTotalFrames-1 {
			p.currentFrame++
			p.player.position[1]--
		}

		v := copyVector(*p.player.velocity)
		v.Scale(platformFrameInterval / 1000.0)

		p.position.Add(&v)
	}
}
